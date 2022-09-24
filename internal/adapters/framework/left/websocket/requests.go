package websocket

import (
	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
)

func NewPublishEvent(m map[string]interface{}) *core.PublishEvent {
	return &core.PublishEvent{
		Channel: m["channel"].(string),
		Event:   m["event"].(core.CloudEvent),
	}
}

func NewSubscribeMessage(m map[string]interface{}) *core.SubscribeMessage {
	channels, ok := m["channels"].([]interface{})
	if !ok {
		panic("channels is not a string slice")
	}
	return &core.SubscribeMessage{
		Channels: core.ConvertToStringSlice(channels),
	}
}

func (c *Client) NewPublishRequest(m map[string]interface{}) *core.PublishRequest[*Client] {
	defer func() {
		if r := recover(); r != nil {
			c.Pool.Logging.Error("websocket::Client.NewPublishRequest => %s", r)
		}
	}()

	event := m["event"].(map[string]interface{})

	meta, ok := event["meta"].(string)

	if !ok {
		meta = ""
	}

	subject, ok := event["subject"].(string)

	if !ok {
		subject = "*"
	}

	return &core.PublishRequest[*Client]{
		PublishEvent: core.PublishEvent{
			Type:    m["type"].(string),
			Channel: m["channel"].(string),
			Event: core.CloudEvent{
				ID:              string(event["id"].(string)),
				Source:          string(event["source"].(string)),
				Type:            string(event["type"].(string)),
				Subject:         string(subject),
				Data:            string(event["data"].(string)),
				SpecVersion:     string(event["specversion"].(string)),
				DataContentType: "application/json",
				Time:            string(event["time"].(string)),
				Meta:            string(meta),
			},
		},
		Client: c,
	}
}

func (c *Client) NewSubscribeRequest(m map[string]interface{}) *core.SubscribeRequest[*Client] {
	defer func() {
		if r := recover(); r != nil {
			c.Pool.Logging.Error("websocket::Client.NewSubscribeRequest => %s", r)
		}
	}()

	c.Pool.Logging.Trace("NewSubscribeRequest: %+v", m)
	channels, ok := m["channels"].([]interface{})
	if !ok {
		c.Pool.Logging.Debug("channels is not a string slice: %+v", m["channels"])
		panic("channels is not a string slice")
	}
	return &core.SubscribeRequest[*Client]{
		SubscribeMessage: core.SubscribeMessage{
			Type:     m["type"].(string),
			Channels: core.ConvertToStringSlice(channels),
		},
		Client: c,
	}
}
