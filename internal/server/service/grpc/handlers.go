package grpc

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	metrics "github.com/dontagr/metric/grpc/gen"
	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service/interfaces"
	serviceModels "github.com/dontagr/metric/internal/server/service/models"
	"github.com/dontagr/metric/models"
)

type Handler struct {
	metrics.UnimplementedMetricServiceServer
	service interfaces.Service
	hashKey string
	log     *zap.SugaredLogger
}

func NewHandler(cnf *config.Config, service interfaces.Service, log *zap.SugaredLogger) *Handler {
	return &Handler{service: service, hashKey: cnf.Security.Key, log: log}
}

func (h *Handler) Value(ctx context.Context, in *metrics.ValueRequest) (*metrics.ValueResponse, error) {
	oldMetric, errEcho := h.service.GetMetric(serviceModels.RequestMetric{
		MType: in.Type,
		MName: in.Id,
	})
	if errEcho != nil {
		return nil, status.Errorf(errorCodeTransfer(errEcho.Code), errEcho.Message)
	}

	md := metadata.New(map[string]string{models.HashAlgKey: oldMetric.Hash})
	if err := grpc.SendHeader(ctx, md); err != nil {
		h.log.Errorf("failed to send header: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to send metadata")
	}

	return &metrics.ValueResponse{
		Id:    oldMetric.ID,
		Type:  oldMetric.MType,
		Delta: oldMetric.Delta,
		Value: oldMetric.Value,
		Hash:  oldMetric.Hash,
	}, nil
}

func (h *Handler) Values(_ context.Context, _ *emptypb.Empty) (*metrics.ValuesResponse, error) {
	collection, errEcho := h.service.GetAllMetrics()
	if errEcho != nil {
		return nil, status.Errorf(errorCodeTransfer(errEcho.Code), errEcho.Message)
	}

	result := metrics.ValuesResponse{}
	for _, row := range collection {
		result.Metric = append(result.Metric, &metrics.ValueResponse{
			Id:    row.ID,
			Type:  row.MType,
			Delta: row.Delta,
			Value: row.Value,
			Hash:  row.Hash,
		})
	}

	return &result, nil
}

func (h *Handler) Update(_ context.Context, in *metrics.UpdateRequest) (*metrics.ValueResponse, error) {
	metric, errEcho := h.service.UpdateMetric(serviceModels.RequestMetric{
		MType: in.Type,
		MName: in.Id,
		Delta: in.Delta,
		Value: in.Value,
		Hash:  in.Hash,
	})
	if errEcho != nil {
		return nil, status.Errorf(errorCodeTransfer(errEcho.Code), errEcho.Message)
	}

	return &metrics.ValueResponse{
		Id:    metric.ID,
		Type:  metric.MType,
		Delta: metric.Delta,
		Value: metric.Value,
		Hash:  metric.Hash,
	}, nil
}

func (h *Handler) Updates(_ context.Context, msg *metrics.UpdatesRequest) (*metrics.ValuesResponse, error) {
	requestArrayMetric := serviceModels.RequestArrayMetric{}
	for _, row := range msg.Metric {
		metric := serviceModels.RequestMetric{
			MType: row.Type,
			MName: row.Id,
			Delta: row.Delta,
			Value: row.Value,
		}

		if row.Hash != nil && *row.Hash != "" {
			metric.Hash = row.Hash
		}

		requestArrayMetric = append(requestArrayMetric, metric)
	}

	upMetrics, errEcho := h.service.UpdateMetrics(requestArrayMetric)
	if errEcho != nil {
		return nil, status.Errorf(errorCodeTransfer(errEcho.Code), errEcho.Message)
	}

	result := metrics.ValuesResponse{}
	for _, row := range upMetrics {
		result.Metric = append(result.Metric, &metrics.ValueResponse{
			Id:    row.ID,
			Type:  row.MType,
			Delta: row.Delta,
			Value: row.Value,
			Hash:  row.Hash,
		})
	}

	return &result, nil
}

func (h *Handler) Ping(_ context.Context, msg *metrics.PingRequest) (*metrics.PingResponse, error) {
	h.log.Infof("msg: %s", msg)

	return &metrics.PingResponse{Message: "pong"}, nil
}

func errorCodeTransfer(errorCode int) codes.Code {
	switch errorCode {
	case http.StatusOK:
		return codes.OK
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusInternalServerError:
		return codes.Internal
	default:
		return codes.Internal
	}
}
