package markets

import (
	influx "github.com/influxdata/influxdb/client/v2"
)

type Metric interface {
	Write(influx.Client, string)
}

type Market interface {
	Open(name string) (Client, error)
}

type Client interface {
	Name() string
	Query(client influx.Client, output chan string)
}

func Open(name string) (Client, error) {
	m := markets[name]
	return m.Open(name)
}
