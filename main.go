package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devfeel/dotweb"
	"github.com/gnuos/marketd/engine"
	"github.com/gnuos/marketd/markets"
	_ "github.com/gnuos/marketd/service"
	"github.com/labstack/gommon/log"
)

var (
	influxdb *engine.InfluxDB

	services = markets.AllMarkets()

	configPath = flag.String("config", "marketd.ini", "Configuration file path.")
)

func main() {
	var err error

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if _, err := os.Open(*configPath); os.IsNotExist(err) {
		flag.Usage()
		os.Exit(1)
	}

	cfg := LoadConfig(configPath)

	influxdb, err = engine.Open(cfg.Influxdb)
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err == nil {
		log.SetOutput(logFile)
	}

	app := dotweb.New()
	app.SetProductionMode()

	app.HttpServer.GET("/", indexHandler)
	app.HttpServer.WebSocket("/ws", serveWS)

	log.Fatal(app.ListenAndServe(cfg.Listen))
}
