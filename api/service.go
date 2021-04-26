package api

import "time"

type Event interface {
	ID() uint64
	Code() int
	Data() interface{}
	String() string
	StageFormat(
		stage EventStageType,
		duration time.Duration) string
}

type EventStageType int

const (
	EventStageFirst EventStageType = iota
	EventStageSecond
	EventStageNext
)

type CallbackFunc func(
	evtStage EventStageType,
	duration time.Duration,
	evt Event,
	count uint64,
	codes map[int]int,
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

func (e EventStageType) String() string {
	switch e {
	case EventStageFirst:
		return "first"
	case EventStageSecond:
		return "second"
	default:
		return "next"
	}
}
