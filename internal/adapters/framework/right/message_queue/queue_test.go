package message_queue

import (
	"testing"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/stretchr/testify/require"
)

func TestMessageQueue(t *testing.T) {

	t.Run("Enqueue", func(t *testing.T) {

		t.Run("should enqueue a message", func(t *testing.T) {
			adapter := NewAdapter()
			message := core.CloudEvent{
				ID:     "123",
				Type:   "test",
				Source: "test",
				Time:   "test",
			}
			adapter.Enqueue("test", message)
			value, _ := adapter.cache.Get("test")
			adapter.Enqueue("test", message)
			require.Equal(t, 1, len(adapter.cache.Items()))
			require.Equal(t, []core.CloudEvent{message, message}, value)
		})

		t.Run("should return an error if the channel is not found", func(t *testing.T) {
			adapter := NewAdapter()
			message := core.CloudEvent{
				ID:     "123",
				Type:   "test",
				Source: "test",
				Time:   "test",
			}
			adapter.Enqueue("test", message)
			_, err := adapter.Dequeue("test2", false)
			require.Error(t, err)
		})

	})

	t.Run("Dequeue", func(t *testing.T) {

		t.Run("should dequeue a message", func(t *testing.T) {
			adapter := NewAdapter()
			message := core.CloudEvent{
				ID:     "123",
				Type:   "test",
				Source: "test",
				Time:   "test",
			}
			adapter.Enqueue("test", message)
			value, _ := adapter.Dequeue("test", false)
			require.Equal(t, message, value)
		})

		t.Run("should return an error if the channel is not found", func(t *testing.T) {
			adapter := NewAdapter()
			_, err := adapter.Dequeue("test", false)
			require.Error(t, err)
		})

	})

	t.Run("Iter", func(t *testing.T) {

		t.Run("should iterate over a message", func(t *testing.T) {
			adapter := NewAdapter()
			message := core.CloudEvent{
				ID:     "123",
				Type:   "test",
				Source: "test",
				Time:   "test",
			}
			adapter.Enqueue("test", message)
			value, _ := adapter.Iter("test", false)
			require.Equal(t, message, <-value)
		})

		t.Run("should return an error if the channel is not found", func(t *testing.T) {
			adapter := NewAdapter()
			message := core.CloudEvent{
				ID:     "123",
				Type:   "test",
				Source: "test",
				Time:   "test",
			}

			adapter.Enqueue("test", message)
			adapter.Enqueue("test", message)
			adapter.Enqueue("test", message)

			iter, err := adapter.Iter("test", true) // true means we want to consume all messages

			if err != nil {
				t.Error(err)
			}

			for event := range iter {
				require.Equal(t, message, event)
			}

			require.NoError(t, err)

			_, err2 := adapter.Iter("test", false)

			require.Error(t, err2)
		})

	})

}
