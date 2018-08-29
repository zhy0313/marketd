package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devfeel/dotweb"
	"github.com/gnuos/marketd/engine"
	_ "github.com/gnuos/marketd/service"
	"github.com/labstack/gommon/log"
)

var (
	configPath = flag.String("config", "marketd.ini", "Configuration file path.")

	influxdb *engine.InfluxDB
)

func main() {
	var err error

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) < 3 {
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

	app.HttpServer.GET("/", index)
	app.HttpServer.WebSocket("/ws", market)
	log.Fatal(app.StartServer(cfg.Listen))
}
