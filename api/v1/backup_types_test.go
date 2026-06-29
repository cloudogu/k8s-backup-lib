package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupStatus_EnsureConditions(t *testing.T) {
	// given
	sut := BackupStatus{}

	// when
	sut.EnsureConditions()

	// then
	if assert.NotNil(t, sut.Conditions) {
		assert.Empty(t, sut.Conditions)
	}

	bytes, err := json.Marshal(sut)
	require.NoError(t, err)
	assert.Contains(t, string(bytes), `"conditions":[]`)
}
