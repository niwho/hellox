package im

import (
	"fmt"
	//	"sync"
	"time"

	"github.com/niwho/hellox/logs"
	"github.com/orcaman/concurrent-map"
)

var LonlyClientsIns *LonlyClients

type LonlyClients struct {
	FC cmap.ConcurrentMap
	// sync.Lock
}

func PushFC(client *Client) {
	if LonlyClientsIns == nil {
		return
	}
	LonlyClientsIns.Push(client)
}

func PopFC(client *Client) {
	if LonlyClientsIns == nil {
		return
	}
	LonlyClientsIns.Pop(client)
	// 等待5s重连，继续前一次匹配
	matched := client.extra.Matched
	if matched != nil && client.extra.Matched.extra.Matched == client {
		go func() {
			select {
			case <-time.After(time.Second * 3):
				logs.Log(logs.F{"client": client, "extra": client.extra.Matched}).Info("dddddddd")
				if matched.isOnline && client.isOnline {
					logs.Log(nil).Info("is online reconnect")
					return
				} else {
					logs.Log(logs.F{"mo": matched.isOnline, "co": client.isOnline}).Info("is offline free")
					matched.extra.Matched = nil
					client.extra.Matched = nil
				}
			}
		}()
	} else {
		client.extra.Matched = nil
	}
}

func InitLonlyClients() {
	LonlyClientsIns = &LonlyClients{
		FC: cmap.New(),
	}
	go func() {
		for {
			select {
			case <-time.After(time.Second * 1):
				LonlyClientsIns.AutoMatching()
			}
		}
	}()
}

func (ll *LonlyClients) Push(client *Client) {
	if client != nil {
		ll.FC.Set(client.String(), client)
	}
}

func (ll *LonlyClients) Pop(client *Client) {
	if client != nil {
		ll.FC.Remove(client.String())
	}
}

func (ll *LonlyClients) AutoMatching() {
	var priveClient *Client
	ll.FC.IterCb(func(key string, v interface{}) {
		if client, ok := v.(*Client); ok && client.isOnline && client.extra.Matched == nil {
			if priveClient == nil {
				priveClient = client
			} else {
				priveClient.extra.Matched = client
				msg := Msg{
					Stamp: time.Now().Unix(),
				}
				client.extra.Matched = priveClient
				msg.Data = fmt.Sprintf("已经匹配【%s】", priveClient.String())
				client.PushMessage(msg)
				msg.Data = fmt.Sprintf("已经匹配【%s】", client.String())
				priveClient.PushMessage(msg)
				priveClient = nil
			}
		}
	})
}
