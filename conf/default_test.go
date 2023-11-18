package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	devDefaults := GetDevDefaults()
	require.EqualValues(t, "dev", devDefaults.Network.Name)

	stageDefaults := GetStageDefaults()
	require.EqualValues(t, "stage", stageDefaults.Network.Name)

	prodDefaults := GetProdDefaults()
	require.EqualValues(t, "prod", prodDefaults.Network.Name)
}