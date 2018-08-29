package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/devfeel/dotweb"
	"github.com/gnuos/marketd/markets"
	_ "github.com/gnuos/marketd/service"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
)

var (
	services = markets.AllMarkets()
)

type HeartBeat struct {
	Pong uint64 `json:"pong"`
}

func market(ctx dotweb.Context) error {
	var err error
	var strMsg string

	for {
		ctx.WebSocket().SendMessage(fmt.Sprintf("{\"ping\":%d", time.Now().UnixNano()/1000000))

		if strMsg, err = ctx.WebSocket().ReadMessage(); err != nil {
			ctx.WebSocket().SendMessage("connect has a error => " + err.Error())
			log.Warn("ReadMessage error => "+err.Error(), "HttpRequest")
			break
		} else {
			hb := new(HeartBeat)

			if err := json.Unmarshal([]byte(strMsg), hb); err != nil {
				ctx.WebSocket().Conn.WriteClose(websocket.CloseNormalClosure)
			}
			ctx.WebSocket().SendMessage(fmt.Sprintf("{\"ping\":%d", time.Now().UnixNano()/1000000))
		}
	}

	return nil
}

func index(ctx dotweb.Context) error {
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

	return ctx.WriteJson(maps)
}
