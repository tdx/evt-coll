package api

import "time"

type Event interface {
	ID() uint64
	Data() interface{}
	String() string
}

type SenderFunc func(
	evt Event,
	count uint64,
)

// Collector counts events by delay rule.
// The first event is always dispatched, then the number of events is counted
// during the SecondDelay. All of the following events counted during
// the NextDelay.
type DelayRule struct {
	Second time.Duration
	Next   time.Duration
}

var DefaultDelayRule = DelayRule{
	Second: time.Duration(time.Minute),
	Next:   time.Duration(3 * time.Minute),
}

type State struct {
	ActiveEvents int
}

type Collector interface {
	RegisterRule(evtID uint64, rule DelayRule)
	Event(evt Event)
	State() State
}
