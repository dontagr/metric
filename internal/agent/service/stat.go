package service

import (
	"math/rand"
	"runtime"
)

type Stats struct {
	MemStats    runtime.MemStats
	PollCount   int
	RandomValue float64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Update() {
	s.PollCount++
	s.RandomValue = rand.Float64()
	runtime.ReadMemStats(&s.MemStats)
}
