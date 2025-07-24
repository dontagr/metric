package service

import "sync"

type StatUpdateWg struct {
	sync.WaitGroup
}

type StatSendWg struct {
	sync.WaitGroup
}

func NewStatUpdateWg() *StatUpdateWg {
	statWg := StatUpdateWg{
		sync.WaitGroup{},
	}

	return &statWg
}

func NewStatSendWg() *StatSendWg {
	statWg := StatSendWg{
		sync.WaitGroup{},
	}

	return &statWg
}
