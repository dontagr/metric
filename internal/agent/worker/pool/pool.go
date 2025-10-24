package pool

import (
	"sync"

	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/agent/config"
)

type WPool struct {
	cfg      *config.Config
	workerWG *sync.WaitGroup
	log      *zap.SugaredLogger
	workers  int
}

func NewWPool(cfg *config.Config, log *zap.SugaredLogger) *WPool {
	workers := 1
	if cfg.RateLimit != 0 {
		workers = cfg.RateLimit
	}

	return &WPool{
		log:      log,
		workers:  workers,
		workerWG: &sync.WaitGroup{},
	}
}

type ProcessJobFunc func(id int, job any)

func (w *WPool) GetWG() *sync.WaitGroup {
	return w.workerWG
}

func (w *WPool) Start(fn ProcessJobFunc) chan any {
	jobs := make(chan any, w.workers)
	for id := 1; id <= w.workers; id++ {
		w.workerWG.Add(1)
		go w.worker(id, jobs, fn)
	}

	return jobs
}

func (w *WPool) worker(id int, jobs chan any, fn ProcessJobFunc) {
	defer w.workerWG.Done()
	w.log.Infof("Worker %d: Начал работу", id)

	for job := range jobs {
		fn(id, job)
	}
	w.log.Infof("Worker %d: Завершил работу", id)
}
