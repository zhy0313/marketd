package huobi

import (
	"sync"

	"github.com/gnuos/marketd/markets"
	"github.com/labstack/gommon/log"
)

type Huobi struct{}

type HandleFunc func() (string, error)

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

func (hc *HuobiClient) Query() string {
	var wg sync.WaitGroup

	var out = "{"

	for name, handle := range hc.Handlers {
		wg.Add(1)

		//开goroutine去获取火币的API信息
		//name: 相当于表的名字
		//handle: 用于获取具体接口数据的函数
		go func(name string, handle HandleFunc) {
			defer wg.Done()
			jsonStr, err := handle()
			if err != nil {
				log.Error(err)
				return
			}

			out += (`"` + name + `":` + jsonStr + ",")
		}(name, handle)
	}

	wg.Wait()

	return out[:len(out)-1] + "}"
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
