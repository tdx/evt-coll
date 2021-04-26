package evtcoll_test

import (
	"testing"
	"time"

	evtcoll "github.com/tdx/evt-coll"
	"github.com/tdx/evt-coll/api"

	"github.com/stretchr/testify/require"
)

type evt struct {
	id   uint64
	data string
}

func (e *evt) ID() uint64 {
	return e.id
}

func (e *evt) Data() interface{} {
	return "not used in test"
}

func (e *evt) String() string {
	return "not used in test"
}

func (e *evt) StageFormat(evtStage api.EventStageType) string {
	return "not used in test"
}

func Test(t *testing.T) {

	var count uint64

	cbFunc := func(evtStage api.EventStageType, evt api.Event, n uint64) {
		count += n
	}

	s := evtcoll.New(cbFunc)
	s.RegisterRule(1, api.DelayRule{
		Second: time.Duration(2 * time.Second),
		Next:   time.Duration(4 * time.Second),
	})

	// first
	e := &evt{1, "test"}
	s.Event(e)
	time.Sleep(time.Second)
	require.Equal(t, uint64(1), count)
	require.Equal(t, 1, s.State().ActiveEvents)

	// second
	count = 0
	s.Event(e)
	s.Event(e)

	time.Sleep(3 * time.Second)
	require.Equal(t, uint64(2), count)

	// next
	count = 0
	s.Event(e)
	time.Sleep(5 * time.Second)
	require.Equal(t, uint64(1), count)

	time.Sleep(5 * time.Second)
	require.Equal(t, uint64(1), count)
	require.Equal(t, 0, s.State().ActiveEvents)
}
