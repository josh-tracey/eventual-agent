package websocket

import (
	"fmt"
)

// Message - Message duck type
type Message interface {
	isMessage()
}

// PublishEvent - Publish incoming message type
type PublishEvent struct {
	Type     string     `json:"type"`
	Channels []string   `json:"channels"`
	Event    CloudEvent `json:"event"`
}

// SubscribeMessage - Subscribe incoming message type
type SubscribeMessage struct {
	Type     string   `json:"type"`
	Channels []string `json:"channels"`
}

func (p PublishEvent) isMessage() {}

func (p SubscribeMessage) isMessage() {}

// CloudEvent - https://github.com/cloudevents/spec/blob/v1.0.1/spec.md
type CloudEvent struct {
	ID              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	Subject         string                 `json:"subject"`
	Data            map[string]interface{} `json:"data"`
	DataContentType string                 `json:"datacontenttype"`
	Time            string                 `json:"time"`
	SpecVersion     string                 `json:"specversion"`
	Meta            map[string]interface{} `json:"meta"`
}

type SubscribeRequest struct {
	SubscribeMessage
	Client *Client
}

type PublishRequest struct {
	PublishEvent
	Client *Client
}

func convertToStringSlice(input []interface{}) []string {
	s := make([]string, len(input))
	for i, v := range input {
		s[i] = fmt.Sprintf("%s", v)
	}
	return s
}

func NewPublishEvent(m map[string]interface{}) *PublishEvent {
	channels, ok := m["channels"].([]interface{})
	if !ok {
		panic("channels is not a string slice")
	}
	return &PublishEvent{
		Channels: convertToStringSlice(channels),
		Event:    m["event"].(CloudEvent),
	}
}

func NewSubscribeMessage(m map[string]interface{}) *SubscribeMessage {
	channels, ok := m["channels"].([]interface{})
	if !ok {
		panic("channels is not a string slice")
	}
	return &SubscribeMessage{
		Channels: convertToStringSlice(channels),
	}
}

func (c *Client) NewPublishRequest(m map[string]interface{}) *PublishRequest {
	defer func() {
		if r := recover(); r != nil {
			c.Pool.Logging.Error("websocket::Client.NewPublishRequest => %s", r)
		}
	}()

	event := m["event"].(map[string]interface{})
	data := make(map[string]interface{})
	for k, v := range event["data"].(map[string]interface{}) {
		data[k] = v
	}

	meta := make(map[string]interface{})

	if event["meta"] != nil {
		for k, v := range event["meta"].(map[string]interface{}) {
			meta[k] = v
		}
	}
	channels, ok := m["channels"].([]interface{})

	if !ok {
		panic("channels is not a string slice")
	}

	subject, ok := event["subject"].(string)

	if !ok {
		subject = "*"
	}

	return &PublishRequest{
		PublishEvent: PublishEvent{
			Type:     m["type"].(string),
			Channels: convertToStringSlice(channels),
			Event: CloudEvent{
				ID:              string(event["id"].(string)),
				Source:          string(event["source"].(string)),
				Type:            string(event["type"].(string)),
				Subject:         string(subject),
				Data:            data,
				SpecVersion:     string(event["specversion"].(string)),
				DataContentType: "application/json",
				Time:            string(event["time"].(string)),
				Meta:            meta,
			},
		},
		Client: c,
	}
}

func (c *Client) NewSubscribeRequest(m map[string]interface{}) *SubscribeRequest {
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
	return &SubscribeRequest{
		SubscribeMessage: SubscribeMessage{
			Type:     m["type"].(string),
			Channels: convertToStringSlice(channels),
		},
		Client: c,
	}
}
