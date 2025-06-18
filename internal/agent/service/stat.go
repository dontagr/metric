package service

import (
	"math/rand"
	"runtime"
)

type Stats struct {
	PollCount     int
	RandomValue   float64
	Alloc         uint64
	BuckHashSys   uint64
	Frees         uint64
	GCCPUFraction float64
	GCSys         uint64
	HeapAlloc     uint64
	HeapIdle      uint64
	HeapInuse     uint64
	HeapObjects   uint64
	HeapReleased  uint64
	HeapSys       uint64
	LastGC        uint64
	Lookups       uint64
	MSpanInuse    uint64
	MSpanSys      uint64
	MCacheInuse   uint64
	MCacheSys     uint64
	Mallocs       uint64
	NextGC        uint64
	NumGC         uint32
	NumForcedGC   uint32
	TotalAlloc    uint64
	Sys           uint64
	StackInuse    uint64
	StackSys      uint64
	OtherSys      uint64
	PauseTotalNs  uint64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Update() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	s.PollCount++
	s.RandomValue = rand.Float64()
	s.Alloc = memStats.Alloc
	s.BuckHashSys = memStats.BuckHashSys
	s.Frees = memStats.Frees
	s.GCCPUFraction = memStats.GCCPUFraction
	s.GCSys = memStats.GCSys
	s.HeapAlloc = memStats.HeapAlloc
	s.HeapIdle = memStats.HeapIdle
	s.HeapInuse = memStats.HeapInuse
	s.HeapObjects = memStats.HeapObjects
	s.HeapReleased = memStats.HeapReleased
	s.HeapSys = memStats.HeapSys
	s.LastGC = memStats.LastGC
	s.Lookups = memStats.Lookups
	s.MSpanInuse = memStats.MSpanInuse
	s.MSpanSys = memStats.MSpanSys
	s.MCacheInuse = memStats.MCacheInuse
	s.MCacheSys = memStats.MCacheSys
	s.Mallocs = memStats.Mallocs
	s.NextGC = memStats.NextGC
	s.NumGC = memStats.NumGC
	s.NumForcedGC = memStats.NumForcedGC
	s.TotalAlloc = memStats.TotalAlloc
	s.Sys = memStats.Sys
	s.StackInuse = memStats.StackInuse
	s.StackSys = memStats.StackSys
	s.OtherSys = memStats.OtherSys
	s.PauseTotalNs = memStats.PauseTotalNs
}
