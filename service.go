package evtcoll

import (
	"sync"

	"github.com/tdx/evt-coll/api"
)

type svc struct {
	callback          api.CallbackFunc
	evtRulesLock      sync.RWMutex
	evtRules          map[uint64]api.DelayRule
	evtStateStoreLock sync.RWMutex
	evtStateStore     map[uint64]*evtState
}

func New(cbFunc api.CallbackFunc) api.Collector {
	if cbFunc == nil {
		return nil
	}
	return &svc{
		callback:      cbFunc,
		evtRules:      make(map[uint64]api.DelayRule),
		evtStateStore: make(map[uint64]*evtState),
	}
}

//
func (s *svc) RegisterRule(evtID uint64, rule api.DelayRule) {
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
