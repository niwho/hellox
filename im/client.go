package im

import (
	//"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/niwho/hellox/logs"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID  string //唯一标志
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	isOnline bool
	// Buffered channel of outbound messages.
	send chan Msg
}

func (c *Client) String() string {
	return c.ID
}

func (c *Client) Close() {
	c.isOnline = false
	c.conn.Close()
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msgJson, err := simplejson.NewJson(message)
		if err != nil {
			fmt.Printf("parse json err:%v", err)
			continue
		}
		to, _ := msgJson.Get("to_user").String()
		text, _ := msgJson.Get("text").String()
		// 是不是直接向to_user用户写消息
		c.hub.hubCh <- Msg{
			FromUser: c.String(),
			ToUser:   to,
			Data:     string(text),
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// 简单发送
			w.Write([]byte(message.Data))

			//		// Add queued chat messages to the current websocket message.
			//		n := len(c.send)
			//		for i := 0; i < n; i++ {
			//			w.Write(newline)
			//			w.Write(<-c.send)
			//		}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	hub := hubIns
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 用户+设备唯一识别码
	uid := "123456" + conn.RemoteAddr().String()
	var client *Client
	if v, ok := hub.clients.Get(uid); ok {
		client, _ = v.(*Client)
	} else {
		client = &Client{ID: uid, hub: hub, conn: conn, send: make(chan Msg, 256), isOnline: true}
		client.hub.register <- client
	}

	logs.Log(logs.F{"uid": uid}).Info("wsssssssssss")
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
