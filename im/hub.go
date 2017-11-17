package im

// hub maintains the set of active clients and broadcasts messages to the
// clients.
import (
	"fmt"

	//	"github.com/golang/protobuf/proto"
	"github.com/orcaman/concurrent-map"
)

type Hub struct {
	// Registered clients.
	//clients map[*Client]int
	clients cmap.ConcurrentMap

	// Inbound messages from the clients.
	hubCh  chan Msg
	userCh chan Msg

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var hubIns *Hub

func init() {
	InitLonlyClients()
	hubIns = NewHub()
	go hubIns.run()
}

func NewHub() *Hub {
	return &Hub{
		hubCh:      make(chan Msg, 1024),
		userCh:     make(chan Msg, 1024),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		//clients:    make(map[*Client]bool),
		clients: cmap.New(),
	}
}

func Broadcast(data string) error {
	select {
	case hubIns.hubCh <- Msg{Data: data}:
	default:
		return fmt.Errorf("broadcast err: overflow")
	}
	return nil
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// h.clients[client] = ONLINE
			h.clients.Set(client.String(), client)
		case client := <-h.unregister:
			_ = client
			// 同时删除这个对象?
			// h.clients.

		case message := <-h.hubCh:
			h.clients.IterCb(func(key string, v interface{}) {
				if client, ok := v.(*Client); ok && client.isOnline {
					select {
					case client.send <- message:
					default:
						// 妥当待确认
						client.Close()
					}
				}
			})
		case message := <-h.userCh:
			if v, ok := h.clients.Get(message.ToUser); ok {
				if client, ok := v.(*Client); ok && client.isOnline {
					select {
					case client.send <- message:
					default:
						// 妥当待确认
						client.Close()
					}
				}
			}
		}
	}
}
