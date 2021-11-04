package websocket

import (
	"time"
)

type Message interface {
	isMessage()
}

type PublishEvent struct {
	Type     string        `json:"type"`
	Channels []interface{} `json:"channels"`
	Event    CloudEvent    `json:"event"`
}

type SubscribeMessage struct {
	Type     string        `json:"type"`
	Channels []interface{} `json:"channels"`
}

func (p PublishEvent) isMessage() {}

func (p SubscribeMessage) isMessage() {}

type CloudEvent struct {
	Id              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	Data            map[string]interface{} `json:"data"`
	DataContentType string                 `json:"datacontenttype"`
	Time            int                    `json:"time"`
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

func NewPublishEvent(m map[string]interface{}) *PublishEvent {
	return &PublishEvent{
		Channels: m["channels"].([]interface{}),
		Event:    m["event"].(CloudEvent),
	}
}

func NewSubscribeMessage(m map[string]interface{}) *SubscribeMessage {
	return &SubscribeMessage{
		Channels: m["channels"].([]interface{}),
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

	return &PublishRequest{
		PublishEvent: PublishEvent{
			Type:     m["type"].(string),
			Channels: m["channels"].([]interface{}),
			Event: CloudEvent{
				Id:              string(event["id"].(string)),
				Source:          string(event["source"].(string)),
				Type:            string(event["type"].(string)),
				Data:            data,
				SpecVersion:     "1.0",
				DataContentType: "application/json",
				Time:            int(time.Now().Unix()),
				Meta:            meta,
			},
		},
		Client: c,
	}
}

func (c *Client) NewSubscribeRequest(m map[string]interface{}) *SubscribeRequest {
	return &SubscribeRequest{
		SubscribeMessage: SubscribeMessage{
			Type:     m["type"].(string),
			Channels: m["channels"].([]interface{}),
		},
		Client: c,
	}
}
