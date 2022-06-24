package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/josh-tracey/eventual-agent/internal/profile"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Pool   *Pool
	Send   chan interface{}
	RefID  string
	closed bool
}

func NewClient(id string, conn *websocket.Conn, pool *Pool) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Pool: pool,
		Send: make(chan interface{}, 32),
	}
}

func (c *Client) AddRefID(refID string) {
	c.RefID = refID
}

func (c *Client) close() {

	defer func() {
		if r := recover(); r != nil {
			c.Pool.Logging.Error("websocket::Client.close => %s", r)
		}
	}()

	if !c.closed {
		if err := c.Conn.Close(); err != nil {
			c.Pool.Logging.Trace("websocket was already closed: %+v", err)
		}
		close(c.Send)
		c.closed = true
	}
}

var (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	PongWait = 30 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10
	// Maximum message size allowed from peer.
	MaxMessageSize int64 = 64 * 1024
)

func dispatch(c *Client, data map[string]interface{}) {

	defer func() {
		if r := recover(); r != nil {
			c.Pool.Logging.Error("websocket::Client.dispatch => %s", r)
		}
	}()

	defer profile.Duration(*c.Pool.Logging, time.Now(), "Client::ReadListen::dispatch")

	switch data["type"].(string) {
	case "publish":
		c.Pool.Logging.Trace("dispatch => publish")
		c.Pool.Publish <- *c.NewPublishRequest(data)
	case "subscribe":
		c.Pool.Logging.Trace("dispatch => subscribe")
		c.Pool.Subscribe <- *c.NewSubscribeRequest(data)
	case "unsubscribe":
		c.Pool.Logging.Trace("dispatch => unsubscribe")
		c.Pool.Unsubscribe <- *c.NewSubscribeRequest(data)
	default:
		c.Pool.Logging.Trace("dispatch => invalid request")
		c.Conn.WriteJSON("Invalid Request")
	}
}

func (c *Client) WriteListen() {

	write := func(mt int, payload interface{}, json bool) error {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(WriteWait)); err != nil {
			return err
		}
		if json {

			return websocket.WriteJSON(c.Conn, payload)
		} else {
			return c.Conn.WriteMessage(mt, payload.([]uint8))
		}
	}

	ticker := time.NewTicker(PingPeriod)
	defer func() {
		if err := recover(); err != nil {
			c.Pool.Logging.Error("websocket::Client.WriteListen => %s", err)
		}
		ticker.Stop()
		c.close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				if err := write(websocket.CloseMessage, []byte{}, false); err != nil {
					c.Pool.Logging.Trace("socket already closed: %+v", err)
					panic(err)
				}
			}
			if err := write(websocket.TextMessage, message, true); err != nil {
				c.Pool.Logging.Trace("failed to write socket message: %+v", err)
				panic(err)
			}
		case <-ticker.C:
			if err := write(websocket.PingMessage, []byte{}, false); err != nil {
				c.Pool.Logging.Trace("failed to ping socket: %+v", err)
				panic(err)
			}
		}
	}
}

func (c *Client) ReadListen() {

	defer func() {
		if err := recover(); err != nil {
			c.Pool.Logging.Error("websocket::Client.ReadListen => %s", err)
		}
		c.close()
	}()
	c.Conn.SetReadLimit(MaxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(PongWait)); err != nil {
		c.Pool.Logging.Error("failed to set socket read deadline: %+v", err)
	}
	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	})

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Pool.Logging.Error("IsUnexpectedCloseError: %+v", err.Error())
				break
			}
			c.Pool.Logging.Error("ReadMessage: %+v", err.Error())
			break
		}
		var data map[string]interface{}

		mErr := json.Unmarshal(p, &data)

		if mErr != nil {
			c.Pool.Logging.Error("json.Unmarshal: %+v", err.Error())
			break
		}
		dispatch(c, data)
	}
}
