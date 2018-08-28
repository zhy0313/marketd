package config

import (
	"log"

	"github.com/gnuos/marketd/config/ini"
)

type Config struct {
	Listen   string //marketd websocket的监听地址
	Influxdb string //运行的InfluxDB数据库的http访问地址
	LogPath  string //marketd websocket的日志路径
}

func LoadConfig(confPath *string) *Config {
	var config = &Config{}

	conf, err := ini.LoadFromFile(confPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Listen, _ = conf.Get("marketd", "listen")
	config.Influxdb, _ = conf.Get("marketd", "influxdb")
	config.LogPath, _ = conf.Get("marketd", "logPath")

	config.Check()

	return config
}

func (c Config) Check() {
	if c.Listen == "" {
		c.Listen = ":3000"
	}

	if c.Influxdb == "" {
		c.Influxdb = "http://localhost:8086"
	}

	if c.LogPath == "" {
		c.LogPath = "../marketd.log"
	}
}
