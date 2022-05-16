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
		s[i] = fmt.Sprint(v)
	}
	return s
}

func NewPublishEvent(m map[string]interface{}) *PublishEvent {
	channels, _ := m["channels"].([]interface{})
	return &PublishEvent{
		Channels: convertToStringSlice(channels),
		Event:    m["event"].(CloudEvent),
	}
}

func NewSubscribeMessage(m map[string]interface{}) *SubscribeMessage {
	channels, _ := m["channels"].([]interface{})
	return &SubscribeMessage{
		Channels: convertToStringSlice(channels),
	}
}

func (c *Client) NewPublishRequest(m map[string]interface{}) *PublishRequest {
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
	channels, _ := m["channels"].([]interface{})

	return &PublishRequest{
		PublishEvent: PublishEvent{
			Type:     m["type"].(string),
			Channels: convertToStringSlice(channels),
			Event: CloudEvent{
				ID:              string(event["id"].(string)),
				Source:          string(event["source"].(string)),
				Type:            string(event["type"].(string)),
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
	c.Pool.Logging.Trace("NewSubscribeRequest: %+v", m)
	channels, _ := m["channels"].([]interface{})
	return &SubscribeRequest{
		SubscribeMessage: SubscribeMessage{
			Type:     m["type"].(string),
			Channels: convertToStringSlice(channels),
		},
		Client: c,
	}
}
