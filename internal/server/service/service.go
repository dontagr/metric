package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dontagr/metric/models"
)

type UpdateHandler struct {
	metricFactory *MetricFactory
	store         Store
}

func NewUpdateHandler(mf *MetricFactory, st Store) *UpdateHandler {
	return &UpdateHandler{
		metricFactory: mf,
		store:         st,
	}
}

func (h *UpdateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	mType, mName, mValue, mLen := getPathVars(req.URL.Path)
	if mLen == 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if mLen != 5 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isValidMetricType(mType) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.metricFactory.GetMetric(mType)
	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	newMetric, err := metric.ConvertToMetrics(mName, mValue)
	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	oldMetric := h.store.LoadMetric(mName, mType)
	err = metric.Process(oldMetric, newMetric)
	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.store.SaveMetric(newMetric)

	fmt.Println(newMetric)
	if newMetric.Value != nil {
		fmt.Println(*newMetric.Value)
	}
	if newMetric.Delta != nil {
		fmt.Println(*newMetric.Delta)
	}
}

func getPathVars(path string) (string, string, string, int) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 5 {
		return "", "", "", len(pathParts)
	}

	return pathParts[2], pathParts[3], pathParts[4], len(pathParts)
}

func isValidMetricType(mType string) bool {
	if mType == "" {
		return false
	}
	if models.Counter != mType && models.Gauge != mType {
		return false
	}
	return true
}
