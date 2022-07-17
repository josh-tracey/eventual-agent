package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertToStringSlice(t *testing.T) {

	input := []interface{}{"TestChannel", "TestChannel2"}
	output := ConvertToStringSlice(input)

	require.Equal(t, []string{"TestChannel", "TestChannel2"}, output)
}
