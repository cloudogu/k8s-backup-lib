package v1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackup_GetTruncatedName(t *testing.T) {
	t.Run("short names stay the same", func(t *testing.T) {
		name := strings.Repeat("a", 63)
		b := &Backup{Name: name}
		assert.Equal(t, name, b.GetTruncatedName())
	})
	t.Run("long names are truncated", func(t *testing.T) {
		backupName1 := strings.Repeat("a", 64)
		backupName2 := fmt.Sprintf("%sb", strings.Repeat("a", 63))
		backup1 := &Backup{Name: backupName1}
		backup2 := &Backup{Name: backupName2}

		truncated1 := backup1.GetTruncatedName()
		truncated2 := backup2.GetTruncatedName()

		assert.Len(t, truncated1, maxLabelValueLength)
		assert.Len(t, truncated2, maxLabelValueLength)
		assert.NotEqual(t, truncated1, truncated2)
		assert.True(t, strings.HasPrefix(truncated1, backupName1[:54]))
		assert.True(t, strings.HasPrefix(truncated2, backupName2[:54]))
	})
}
