package service

import "sync"

type StatUpdateWg struct {
	sync.WaitGroup
}

type StatSendWg struct {
	sync.WaitGroup
}

func newStatUpdateWg() *StatUpdateWg {
	statWg := StatUpdateWg{
		sync.WaitGroup{},
	}

	return &statWg
}

func newStatSendWg() *StatSendWg {
	statWg := StatSendWg{
		sync.WaitGroup{},
	}

	return &statWg
}
