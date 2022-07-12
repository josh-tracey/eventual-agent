package websocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPublishEvent(t *testing.T) {
	inputEvent := map[string]interface{}{
		"channels": []string{"TestChannel"},
		"event": CloudEvent{
			Type:            "TempUpdate",
			Source:          "com.adriftdev.server001",
			Data:            map[string]interface{}{"cpu01": 56.0, "cpu02": 78.0},
			DataContentType: "application/json",
			ID:              "1234567",
			Time:            "2020-01-01T00:00:00Z",
			SpecVersion:     "1.0",
			Subject:         "CpuTemps",
		}}

	event := NewPublishEvent(inputEvent)

	require.Equal(t, []string{"TestChannel"}, event.Channels)
	require.Equal(t, "TempUpdate", event.Event.Type)
	require.Equal(t, "com.adriftdev.server001", event.Event.Source)
	require.Equal(t, 56.0, event.Event.Data["cpu01"])
}
