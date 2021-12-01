package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/josh-tracey/eventual-agent/src/profile"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

var (
	writeWait = 10 * time.Second
)

func (c *Client) Send(data interface{}) chan error {

	errChan := make(chan error, 4)

	go func() {
		defer profile.Duration(*c.Pool.Logging, time.Now(), "client.Send")
		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		errChan <- c.Conn.WriteJSON(data)
	}()

	return errChan
}

func getData(c Client) (map[string]interface{}, error) {
	defer profile.Duration(*c.Pool.Logging, time.Now(), "getData")

	_, p, err := c.Conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			c.Pool.Logging.Warn("%v", err)
		}
		return nil, err
	}
	var data map[string]interface{}

	mErr := json.Unmarshal(p, &data)

	if mErr != nil {
		return nil, mErr
	}
	return data, nil
}

func dispatch(c Client, data map[string]interface{}) {
	defer profile.Duration(*c.Pool.Logging, time.Now(), "dispatch")

	switch data["type"].(string) {
	case "publish":
		c.Pool.Publish <- *c.NewPublishRequest(data)
	case "subscribe":
		c.Pool.Subcribe <- *c.NewSubscribeRequest(data)
	case "unsubscribe":
		c.Pool.Unsubscribe <- *c.NewSubscribeRequest(data)
	default:
		c.Conn.WriteJSON("Invalid Request")
	}
}

func (c *Client) Read() {

	defer func() {
		c.Pool.UnsubscribeAll <- SubscribeRequest{Client: c}
		c.Conn.Close()
	}()

	for {
		data, err := getData(*c)
		if err != nil {
			c.Pool.Logging.Error(err.Error())
		}
		dispatch(*c, data)
	}
}
