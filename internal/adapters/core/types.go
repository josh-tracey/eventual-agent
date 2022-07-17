package core

import (
	"fmt"
)

// Message - Message duck type
type Message interface {
	isMessage()
}
type BaseClient interface {
	isClient()
}

type CoreClient struct {
	RefID string
}

func (c *CoreClient) AddRefID(refID string) {
	c.RefID = refID
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

type SubscribeRequest[T any] struct {
	SubscribeMessage
	Client T
}

type PublishRequest[T any] struct {
	PublishEvent
	Client T
}

func ConvertToStringSlice(input []interface{}) []string {
	s := make([]string, len(input))
	for i, v := range input {
		s[i] = fmt.Sprintf("%s", v)
	}
	return s
}
