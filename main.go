package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gnuos/marketd/config"
	"github.com/gnuos/marketd/engine"
	"github.com/gnuos/marketd/markets"
	_ "github.com/gnuos/marketd/service"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	influxdb *engine.InfluxDB
	logFile  *os.File

	configPath = flag.String("config", "", "Configuration file path.")

	services = markets.AllMarkets()
	infoChan = make(chan string)
	done     = make(chan bool)
)

func market(c echo.Context) error {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Market Daemon Complete!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

func index(c echo.Context) error {
	defer close(infoChan)
	defer close(done)

	for _, srv := range services {
		m, err := markets.Open(srv)
		if err != nil {
			log.Fatal(err)
		}
		m.Query(influxdb.Client, infoChan)

		done <- true
		read(m)
	}

	return c.String(200, "Everything is over")
}

func read(cl markets.Client) {
	for {
		if <-done {
			cl.Close()
			return
		}

		js, more := <-infoChan
		if more {
			log.Println(js)
		} else {
			close(infoChan)
			return
		}
	}
}

func main() {
	var err error

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg := config.LoadConfig(configPath)

	influxdb, err = engine.Open(cfg.Influxdb)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.OpenFile(cfg.LogPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(1)

	e.GET("/", index)
	e.GET("/ws", market)

	e.Logger.Fatal(e.Start(cfg.Listen))
}
