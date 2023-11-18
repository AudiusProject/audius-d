package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	devDefaults := GetDevnetDefaults()
	require.EqualValues(t, "devnet", devDefaults.Network.Name)

	stageDefaults := GetTestnetDefaults()
	require.EqualValues(t, "testnet", stageDefaults.Network.Name)

	prodDefaults := GetMainnetDefaults()
	require.EqualValues(t, "mainnet", prodDefaults.Network.Name)
}