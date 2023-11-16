package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevnetConfig(t *testing.T) {
	devnet := ReadTomlUnsafe("../templates/devnet.toml")
	require.EqualValues(t, "devnet", devnet.Network.Name)
	require.EqualValues(t, 2, len(devnet.CreatorNodes))
	require.EqualValues(t, 1, len(devnet.DiscoveryNodes))
}

func TestOperatorConfig(t *testing.T) {
	op := ReadTomlUnsafe("../templates/operator.toml")
	require.EqualValues(t, "stage", op.Network.Name)
	require.EqualValues(t, 1, len(op.CreatorNodes))
	require.EqualValues(t, 0, len(op.DiscoveryNodes))
}
