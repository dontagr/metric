package worker

import (
	"github.com/dontagr/metric/models"
)

type batchModel struct{}

func (bm *batchModel) GetJobs(s *Sender, jobs chan any) {
	metrics := make([]*models.Metrics, 0, len(EnableStats))
	for index, mType := range EnableStats {
		metric, err := s.getMetric(mType, index)
		if err != nil {
			s.log.Errorf("get metrics for index %s: %v", index, err)
			continue
		}

		metrics = append(metrics, metric)
	}

	jobs <- metrics
}
