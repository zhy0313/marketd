package markets

import (
	"github.com/gnuos/marketd/engine"
)

type Metric interface {
	Write(*engine.InfluxDB, string, string)
}

type Market interface {
	Open(name string) (Client, error)
}

type Client interface {
	Close()
	Name() string
	Query(*Rows) chan string
	Write(*engine.InfluxDB, string, *Rows) error
}

func Open(name string) (Client, error) {
	m := markets[name]
	return m.Open(name)
}

type Rows struct {
	Data map[string]Metric
}

func (r *Rows) Query(key string) Metric {
	return r.Data[key]
}

func (r *Rows) Add(key string, metric Metric) {
	r.Data[key] = metric
}

func (r *Rows) Del(key string) {
	delete(r.Data, key)
}
