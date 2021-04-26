package evtcoll

import (
	"sync/atomic"
	"time"

	"github.com/tdx/evt-coll/api"
)

type evtState struct {
	count uint64
	stage api.EventStageType
	rule  api.DelayRule
	codes map[int]int
}

//
func (s *svc) addEvent(evt api.Event) {
	s.evtStateStoreLock.RLock()
	st, ok := s.evtStateStore[evt.ID()]
	s.evtStateStoreLock.RUnlock()

	if ok {
		atomic.AddUint64(&st.count, 1)
		if evt.Code() != 0 {
			s.evtStateStoreLock.Lock()
			st.codes[evt.Code()]++
			s.evtStateStoreLock.Unlock()
		}
		return
	}

	s.evtStateStoreLock.Lock()
	defer s.evtStateStoreLock.Unlock()

	st = &evtState{
		rule:  s.evtRule(evt.ID()),
		codes: make(map[int]int),
	}

	if evt.Code() != 0 {
		st.codes[evt.Code()]++
	}

	s.evtStateStore[evt.ID()] = st

	go s.collectEvents(st, evt)
}

func (s *svc) delEvent(id uint64) {
	s.evtStateStoreLock.Lock()
	defer s.evtStateStoreLock.Unlock()

	delete(s.evtStateStore, id)
}

func (s *svc) events() int {
	s.evtStateStoreLock.RLock()
	defer s.evtStateStoreLock.RUnlock()

	return len(s.evtStateStore)
}

func (s *svc) evtRule(id uint64) api.DelayRule {
	s.evtRulesLock.RLock()
	defer s.evtRulesLock.RUnlock()

	r, ok := s.evtRules[id]
	if !ok {
		return api.DefaultDelayRule
	}

	return r
}

func (s *svc) collectEvents(
	st *evtState, firstEvent api.Event) {

	st.stage = api.EventStageFirst

	s.callback(st.stage, time.Duration(time.Minute), firstEvent, 1, nil)

	st.stage = api.EventStageSecond

	d := st.rule.Second
	for {
		if st.stage == api.EventStageNext {
			time.Sleep(st.rule.Next)
		} else {
			time.Sleep(st.rule.Second)
		}

		count := atomic.SwapUint64(&st.count, 0)
		if count == 0 {
			s.delEvent(firstEvent.ID())
			return
		}

		s.callback(st.stage, d, firstEvent, count, st.codesCopy(s))

		if st.stage == api.EventStageSecond {
			d = st.rule.Next
			st.stage = api.EventStageNext
		}
	}
}

func (st *evtState) codesCopy(s *svc) map[int]int {

	s.evtStateStoreLock.RLock()
	defer s.evtStateStoreLock.RUnlock()

	if len(st.codes) == 0 {
		return nil
	}

	m := make(map[int]int, len(st.codes))
	for k, v := range st.codes {
		m[k] = v
	}

	return m
}
