package websocket

import (
	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
)

func NewPublishEvent(m map[string]interface{}) *core.PublishEvent {
	channels, ok := m["channels"].([]interface{})
	if !ok {
		panic("channels is not a string slice")
	}
	return &core.PublishEvent{
		Channels: core.ConvertToStringSlice(channels),
		Event:    m["event"].(core.CloudEvent),
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

func NewHistoryMessage(m map[string]interface{}) *core.HistoryMessage {
	return &core.HistoryMessage{
		Channel: m["channel"].(string),
		Consume: m["consume"].(bool),
	}
}
