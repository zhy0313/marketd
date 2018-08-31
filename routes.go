package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/devfeel/dotweb"
	"github.com/gnuos/marketd/markets"
	"github.com/labstack/gommon/log"
)

func indexHandler(ctx dotweb.Context) error {
	var maps = make(map[string]map[string]markets.Metric)
	for _, srv := range services {
		m, err := markets.Open(srv)
		if err != nil {
			log.Fatal(err)
		}
		rows := &markets.Rows{
			Data: make(map[string]markets.Metric),
		}

		for info := range m.Query(rows) {
			fmt.Fprintln(ioutil.Discard, info)
		}

		maps[m.Name()] = rows.Data
	}

	res, err := json.Marshal(maps)
	if err != nil {
		return err
	}

	return ctx.WriteBlob("application/json", res)
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
