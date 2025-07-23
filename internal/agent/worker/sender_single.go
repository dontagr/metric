package worker

type singleModel struct{}

func (sb *singleModel) GetJobs(s *Sender, jobs chan any) {
	for index, mType := range EnableStats {
		metric, err := s.getMetric(mType, index)
		if err != nil {
			s.log.Errorf("get metrics for index %s: %v", index, err)
			continue
		}

		jobs <- metric
	}
}
