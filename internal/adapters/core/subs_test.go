package core

import (
	"testing"

	"github.com/josh-tracey/eventual-agent/internal/logging"
	"github.com/stretchr/testify/require"
)

func TestSubStruct(t *testing.T) {

	// Setup
	var (
		Sub *sub
	)

	// Test
	Sub = newSub("ChannelName")

	// Assert
	require.NotNil(t, Sub)
	require.Equal(t, "ChannelName", Sub.ID)

	client1 := "123"
	client2 := "456"

	id1 := Sub.AddClient(&client1)
	id2 := Sub.AddClient(&client2)

	require.Equal(t, []string{id1, id2}, Sub.GetClients())

	Sub.RemoveClientId(id1)

	require.Equal(t, []string{id2}, Sub.GetClients())
	Sub.RemoveClientId(id2)
	require.Equal(t, []string(nil), Sub.GetClients())
}

func TestAdapter(t *testing.T) {
	// Setup
	var (
		adapt *Adapter
	)

	// Test
	logger := logging.NewLogger()
	adapt = NewAdapter(logger)

	clientId1 := adapt.AddClient("10.0.0.5:8000", []string{"tempUpdates"})
	clientId2 := adapt.AddClient("10.0.0.10:8000", []string{"serverUpdates"})

	require.ElementsMatch(t, []string{clientId1, clientId2}, adapt.GetClients())

	adapt.RemoveClientId(clientId1)

	require.Equal(t, []string{clientId2}, adapt.GetClients())

	adapt.RemoveClientId(clientId2)

	require.Equal(t, []string(nil), adapt.GetClients())
}
