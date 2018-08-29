package engine

import (
	"log"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

type Metric struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
}

type InfluxDB struct {
	Client influx.Client
}

func Open(url string) (*InfluxDB, error) {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{Addr: url})
	if err != nil {
		return nil, err
	}

	var db = &InfluxDB{
		Client: client,
	}

	return db, nil
}

func (db *InfluxDB) Close() error {
	return db.Client.Close()
}

func (db *InfluxDB) Query(dbname string, cmd string) (res []influx.Result, err error) {
	defer db.Client.Close()

	q := influx.Query{
		Command:  cmd,
		Database: dbname,
	}

	if response, err := db.Client.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (db *InfluxDB) Write(dbname string, metrics ...Metric) {
	defer db.Client.Close()

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  dbname,
		Precision: "us",
	})
	if err != nil {
		log.Println(err)
	}

	for _, m := range metrics {
		pt, err := influx.NewPoint(
			m.Name,
			m.Tags,
			m.Fields,
			time.Now(),
		)
		if err != nil {
			log.Println(err)
			continue
		}
		bp.AddPoint(pt)
	}

	if err := db.Client.Write(bp); err != nil {
		log.Println(err)
	}
}
