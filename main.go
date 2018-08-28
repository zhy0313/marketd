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

func list(c echo.Context) error {
	done := make(chan struct{})
	defer close(infoChan)
	for _, srv := range services {
		m, err := markets.Open(srv)
		if err != nil {
			log.Fatal(err)
		}
		m.Query(influxdb.Client, infoChan)
	}

	done <- struct{}{}

	select {
	case <-done:
		c.String(200, "Everything is over")
	case info := <-infoChan:
		log.Println(info)
	}
	return nil
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

	logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(1)

	e.GET("/", list)
	e.GET("/ws", market)
	e.Logger.Fatal(e.Start(cfg.Listen))
}
