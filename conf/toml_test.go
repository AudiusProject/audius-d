package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevnetConfig(t *testing.T) {
	devnet := ReadTomlUnsafe("../devnet.toml")
	require.EqualValues(t, "devnet", devnet.Network.Name)
	require.EqualValues(t, 2, len(devnet.CreatorNodes))
	require.EqualValues(t, 1, len(devnet.DiscoveryNodes))
}

func TestSpConfig(t *testing.T) {
	sp := ReadTomlUnsafe("../sp.toml")
	require.EqualValues(t, "devnet", sp.Network.Name)
	require.EqualValues(t, 2, len(sp.CreatorNodes))
	require.EqualValues(t, 1, len(sp.DiscoveryNodes))
}
