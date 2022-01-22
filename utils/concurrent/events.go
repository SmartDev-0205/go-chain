package concurrent

import (
	"sync"

	"github.com/galaxy-foundation/icicb-base/hash"
)

type EventsSet struct {
	sync.RWMutex
	Val hash.EventsSet
}

func WrapEventsSet(v hash.EventsSet) *EventsSet {
	return &EventsSet{
		RWMutex: sync.RWMutex{},
		Val:     v,
	}
}
