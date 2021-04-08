package evtcoll_test

import (
	"testing"
	"time"

	"github.com/tdx/evt-call/evtcoll"
	"github.com/tdx/evt-call/evtcoll/api"

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
	return e.data
}

func (e *evt) String() string {
	return e.data
}

func Test(t *testing.T) {

	var count uint64

	sFunc := func(evt api.Event, n uint64) {
		count += n
	}

	s := evtcoll.New(sFunc)
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
