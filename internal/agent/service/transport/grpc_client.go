package transport

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	metrics "github.com/dontagr/metric/grpc/gen"
	"github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/models"
)

type GRPCManager struct {
	client metrics.MetricServiceClient
	log    *zap.SugaredLogger
}

func NewGRPCManager(cfg *config.Config, log *zap.SugaredLogger, lc fx.Lifecycle) (*GRPCManager, error) {
	ip, err := getIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get IP: %v", err)
	}

	conn, err := grpc.NewClient(
		cfg.GRPCBindAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			ipInterceptor(ip),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("new GRPC Client: %v", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return conn.Close()
		},
	})

	return &GRPCManager{client: metrics.NewMetricServiceClient(conn), log: log}, nil
}

func (h *GRPCManager) NewRequest(income any, _ []string, w int) error {
	switch v := income.(type) {
	case []*models.Metrics:
		collection := make([]*metrics.UpdateRequest, 0, len(v))
		for _, i := range v {
			collection = append(collection, &metrics.UpdateRequest{
				Id:    i.ID,
				Type:  i.MType,
				Delta: i.Delta,
				Value: i.Value,
				Hash:  &i.Hash,
			})
		}

		_, err := h.client.Updates(context.Background(), &metrics.UpdatesRequest{Metric: collection})
		if err != nil {
			return err
		}

		return nil
	case *models.Metrics:
		_, err := h.client.Update(context.Background(), &metrics.UpdateRequest{
			Id:    v.ID,
			Type:  v.MType,
			Delta: v.Delta,
			Value: v.Value,
			Hash:  &v.Hash,
		})
		if err != nil {
			return err
		}

		return nil
	default:
		h.log.Errorf("wrong type income variable")
	}

	return nil
}
