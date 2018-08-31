package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devfeel/dotweb"
	"github.com/gnuos/marketd/markets"
	"github.com/labstack/gommon/log"
)

func main() {
	var (
		err     error
		host    string
		port    string
		logPath string
	)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.StringVar(&host, "host", "", "Default host.")
	flag.StringVar(&port, "port", "3000", "Default port to use for HTTP")
	flag.StringVar(&logPath, "log", "marketd.log", "Process log file")

	flag.Parse()

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err == nil {
		log.SetOutput(logFile)
	}

	app := dotweb.New()
	app.SetProductionMode()

	app.HttpServer.GET("/markets", allMarket)
	app.HttpServer.GET("/markets/:name", oneMarket)
	app.HttpServer.WebSocket("/ws", serveWS)

	log.Fatal(app.ListenAndServe(host + ":" + port))
}

func GetMarket(name string) string {
	m, err := markets.Open(name)
	if err != nil {
		log.Error(err)
		return ""
	}

	return m.Query()
}
