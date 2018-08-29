package main

import (
	"strconv"

	"github.com/gnuos/marketd/utils/ini"
	"github.com/labstack/gommon/log"
)

type config struct {
	Listen   int    //marketd websocket的监听地址
	Influxdb string //运行的InfluxDB数据库的http访问地址
	LogPath  string //日志文件的存放路径
}

func LoadConfig(confPath *string) *config {
	var c = &config{}

	conf, err := ini.LoadFromFile(confPath)
	if err != nil {
		log.Fatal(err)
	}

	if listen, ok := conf.Get("marketd", "listen"); ok {
		port, err := strconv.Atoi(listen)
		if err != nil {
			panic(err)
		}
		if port > 65535 || port < 1 {
			panic("端口超出界限，请填写一个在 1~65535 之间的端口")
		}
	} else {
		c.Listen = 0
	}
	c.Influxdb, _ = conf.Get("marketd", "influxdb")
	c.LogPath, _ = conf.Get("marketd", "logpath")

	c.Check()

	return c
}

func (c *config) Check() {
	if c.Listen == 0 {
		c.Listen = 3000
	}

	if c.Influxdb == "" {
		c.Influxdb = "http://localhost:8086"
	}

	if c.LogPath == "" {
		c.LogPath = "./marketd.log"
	}
}
