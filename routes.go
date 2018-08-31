package main

import (
	"github.com/devfeel/dotweb"
	"github.com/gnuos/marketd/markets"
)

var (
	services = markets.AllMarkets()
)

func allMarket(ctx dotweb.Context) error {
	res := "{"
	for _, srv := range services {
		data := GetMarket(srv)
		res += (`"` + srv + `":{` + data + "},")
	}

	return ctx.WriteBlob("application/json;charset=UTF-8", []byte(res[:len(res)-1]+"}"))
}

func oneMarket(ctx dotweb.Context) error {
	name := ctx.GetRouterName("name")
	data := GetMarket(name)

	return ctx.WriteBlob("application/json;charset=UTF-8", []byte("{"+data+"}"))
}

func serveWS(ctx dotweb.Context) error {
	client := &Client{
		closed: make(chan struct{}),
		send:   make(chan string, 1000),
		ws:     ctx.WebSocket(),
	}

	client.wg.Add(2)

	go client.readLoop()
	go client.pushLoop()

	client.wg.Wait()

	return nil
}
