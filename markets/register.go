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
		panic("Register market is nil")
	}
	if _, dup := markets[name]; dup {
		panic("Register called twice for market " + name)
	}
	markets[name] = market
}

func unregisterAllMarkets() {
	marketsMu.Lock()
	defer marketsMu.Unlock()
	markets = make(map[string]Market)
}

// Markets returns a sorted list of the names of the registered markets.
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
