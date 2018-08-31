package markets

import "errors"

type Market interface {
	Open(name string) (Client, error)
}

type Client interface {
	Close()
	Name() string
	Query() string
}

func Open(name string) (Client, error) {
	m, ok := markets[name]
	if !ok {
		return nil, errors.New(name + " serive is unregistered.")
	}

	return m.Open(name)
}
