package main

import (
	"github.com/gnuos/marketd/utils/ini"
	"github.com/labstack/gommon/log"
)

type config struct {
	Listen   string //marketd websocket的监听地址
	Influxdb string //运行的InfluxDB数据库的http访问地址
	LogPath  string //日志文件的存放路径
}

func LoadConfig(confPath *string) *config {
	var c = &config{}

	conf, err := ini.LoadFromFile(confPath)
	if err != nil {
		log.Fatal(err)
	}

	c.Listen, _ = conf.Get("marketd", "listen")
	c.Influxdb, _ = conf.Get("marketd", "influxdb")
	c.LogPath, _ = conf.Get("marketd", "logpath")

	c.Check()

	return c
}

func (c *config) Check() {
	if c.Listen == "" {
		c.Listen = "localhost:3000"
	}

	if c.Influxdb == "" {
		c.Influxdb = "http://localhost:8086"
	}

	if c.LogPath == "" {
		c.LogPath = "./marketd.log"
	}
}
