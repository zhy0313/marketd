package main

import (
	"io"
	stdLog "log"
	"sync"
	"time"

	"github.com/devfeel/dotweb"
	_ "github.com/gnuos/marketd/service"
	"github.com/labstack/gommon/log"
)

type Client struct {
	wg     sync.WaitGroup
	closed chan struct{}
	send   chan string
	ws     *dotweb.WebSocket
}

func (c *Client) readLoop() {
	defer func() {
		c.ws.Conn.Close()
		c.wg.Done()
	}()

	for {
		//判断是否为心跳信息
		message, err := c.ws.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Warn("连接被客户端关闭了！")
				c.closed <- struct{}{}
				break
			}

			log.Errorf("error: %v", err)
			break
		}

		switch message {
		case "__PING__":
			if err := c.ws.SendMessage("__PONG__"); err != nil {
				log.Error(err)
				return
			}
		case "__PONG__":
			stdLog.Println(c.ws.Request().RemoteAddr, " is connected!")
		default:
			c.send <- message
		}
	}
}

func (c *Client) pushLoop() {
	var ticker = time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		c.ws.Conn.Close()
		c.wg.Done()
	}()

	for {
		select {
		case <-c.closed:
			return
		case message, ok := <-c.send:
			if !ok {
				c.ws.Conn.WriteClose(1000)
				return
			}
			//将消息正式发送给客户端
			err := c.ws.SendMessage(message)
			if err != nil {
				log.Error(err)
			}
		case <-ticker.C: //发送心跳信息
			c.send <- "__PING__"
		default:
			time.Sleep(200 * time.Millisecond)
			//持续向客户端发送消息
			for _, srv := range services {
				data := GetMarket(srv)
				if err := c.ws.SendMessage(`{"` + srv + `":` + data + "}"); err != nil {
					log.Error(err)
				}
			}
		}
	}
}
