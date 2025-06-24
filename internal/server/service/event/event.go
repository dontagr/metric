package event

type Event struct {
	Metrics chan interface{}
}

func NewEvent() *Event {
	return &Event{
		Metrics: make(chan interface{}),
	}
}
