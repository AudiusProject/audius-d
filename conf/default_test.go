package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevDefaults(t *testing.T) {
	devDefaults := GetDevDefaults()
	require.EqualValues(t, "dev", devDefaults.Network.Name)
}