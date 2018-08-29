package huobi

import (
	"log"
	"sync"

	"github.com/gnuos/marketd/markets"
	influx "github.com/influxdata/influxdb/client/v2"
)

type HandleFunc func() (string, markets.Metric, error)

type Huobi struct{}

type HuobiClient struct {
	db influx.Client

	MarketAPI map[string]string
	TradeAPI  map[string]string
	Handlers  map[string]HandleFunc
}

func (hc *HuobiClient) Close() {
}

func (hc *HuobiClient) Name() string {
	return "huobi"
}

func (hc *HuobiClient) Query() chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	for name, handle := range hc.Handlers {
		wg.Add(1)
		go func(name string, handle HandleFunc) {
			jsonString, metric, err := handle()

			if err != nil {
				log.Println(name+": ", err)
				return
			}
			out <- name + " ===> " + jsonString

			metric.Write(hc.db, name)

			wg.Done()
		}(name, handle)
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (hc *HuobiClient) Metrics() []string {
	var metrics = make([]string, 0)

	for k, _ := range hc.Handlers {
		metrics = append(metrics, k)
	}

	return metrics
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
