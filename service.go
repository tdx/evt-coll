package evtcoll

import (
	"sync"

	"github.com/tdx/evt-call/evtcoll/api"
)

type svc struct {
	senderFunc        api.SenderFunc
	evtRulesLock      sync.RWMutex
	evtRules          map[int]api.DelayRule
	evtStateStoreLock sync.RWMutex
	evtStateStore     map[int]*evtState
}

func New(sender api.SenderFunc) api.Collector {
	if sender == nil {
		return nil
	}
	return &svc{
		senderFunc:    sender,
		evtRules:      make(map[int]api.DelayRule),
		evtStateStore: make(map[int]*evtState),
	}
}

//
func (s *svc) RegisterRule(evtID int, rule api.DelayRule) {
	s.evtRulesLock.Lock()
	defer s.evtRulesLock.Unlock()

	s.evtRules[evtID] = rule
}

//
func (s *svc) Event(evt api.Event) {
	s.addEvent(evt)
}

// Active events in collector
func (s *svc) State() api.State {
	return api.State{
		ActiveEvents: s.events(),
	}
}
