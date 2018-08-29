package huobi

import (
	"sync"

	"github.com/gnuos/marketd/engine"
	"github.com/gnuos/marketd/markets"
	"github.com/labstack/gommon/log"
)

type Huobi struct{}

type HandleFunc func() (string, markets.Metric, error)

type HuobiClient struct {
	MarketAPI map[string]string
	TradeAPI  map[string]string
	Handlers  map[string]HandleFunc
}

func (hc *HuobiClient) Close() {
}

func (hc *HuobiClient) Name() string {
	return "huobi"
}

func (hc *HuobiClient) Query(rows *markets.Rows) (out chan string) {
	var wg sync.WaitGroup
	out = make(chan string)

	for name, handle := range hc.Handlers {
		wg.Add(1)

		//开goroutine去获取火币的API信息，把每次的消息都发到管道中
		//name: 相当于表的名字
		//handle: 用于获取具体接口数据的函数
		go func(name string, result *markets.Rows, handle HandleFunc) {
			defer wg.Done()
			jsonString, metric, err := handle()
			if err != nil {
				log.Error(err)
				return
			}

			result.Add(name, metric)
			out <- jsonString
		}(name, rows, handle)
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	return
}

func (hc *HuobiClient) Write(db *engine.InfluxDB, dbname string, rows *markets.Rows) error {
	for table, metric := range rows.Data {
		metric.Write(db, dbname, table)
	}

	return db.Close()
}

func (h *Huobi) Open(name string) (markets.Client, error) {
	hc := &HuobiClient{
		MarketAPI: marketAPI,
		TradeAPI:  tradeAPI,
		Handlers:  handlers,
	}

	return hc, nil
}

func init() {
	markets.Register("huobi", &Huobi{})
}
