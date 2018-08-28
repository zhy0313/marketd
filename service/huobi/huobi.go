package huobi

import (
	"log"
	"sync"

	"github.com/gnuos/marketd/markets"
	influx "github.com/influxdata/influxdb/client/v2"
)

type Huobi struct{}

type HuobiClient struct {
	MarketAPI map[string]string
	TradeAPI  map[string]string
	Handlers  map[string]func() (string, markets.Metric)
}

func (hc *HuobiClient) Name() string {
	return "huobi"
}

func (hc *HuobiClient) Query(client influx.Client, output chan string) {
	var wg sync.WaitGroup
	for name, handle := range hc.Handlers {
		wg.Add(1)
		go func(name string, handle func() (string, markets.Metric)) {
			defer wg.Done()
			jsonString, metric := handle()

			if jsonString == "" {
				return
			}

			log.Println(jsonString)

			output <- jsonString
			metric.Write(client, name)
		}(name, handle)
	}
	wg.Wait()
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
