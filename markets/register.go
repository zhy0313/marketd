package markets

import (
	"sort"
	"sync"
)

var (
	marketsMu sync.RWMutex
	markets   = make(map[string]Market)
)

func Register(name string, market Market) {
	marketsMu.Lock()
	defer marketsMu.Unlock()
	if market == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := markets[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	markets[name] = market
}

func unregisterAllMarkets() {
	marketsMu.Lock()
	defer marketsMu.Unlock()
	// For tests.
	markets = make(map[string]Market)
}

// Drivers returns a sorted list of the names of the registered drivers.
func AllMarkets() []string {
	marketsMu.RLock()
	defer marketsMu.RUnlock()
	var list []string
	for name := range markets {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
