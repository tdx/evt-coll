package evtcoll

import (
	"sync/atomic"
	"time"

	"github.com/tdx/evt-call/evtcoll/api"
)

type evtState struct {
	count  uint64
	isNext bool
	rule   api.DelayRule
}

//
func (s *svc) addEvent(evt api.Event) {
	s.evtStateStoreLock.RLock()
	st, ok := s.evtStateStore[evt.ID()]
	s.evtStateStoreLock.RUnlock()

	if ok {
		atomic.AddUint64(&st.count, 1)
		return
	}

	s.evtStateStoreLock.Lock()
	defer s.evtStateStoreLock.Unlock()

	st = &evtState{
		rule: s.evtRule(evt.ID()),
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

	s.senderFunc(firstEvent, 1)

	for {
		if st.isNext {
			time.Sleep(st.rule.Next)
		} else {
			time.Sleep(st.rule.Second)
		}

		count := atomic.SwapUint64(&st.count, 0)
		if count == 0 {
			s.delEvent(firstEvent.ID())
			return
		}

		s.senderFunc(firstEvent, count)

		st.isNext = true
	}
}
