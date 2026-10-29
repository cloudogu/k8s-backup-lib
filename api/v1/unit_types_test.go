package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupStatus_Units(t *testing.T) {
	t.Run("should return no units for empty status", func(t *testing.T) {
		assert.Empty(t, (&BackupStatus{}).Units())
	})

	t.Run("should return all units with identifiers", func(t *testing.T) {
		status := BackupStatus{
			GlobalConfig: &GlobalConfigStatus{
				State: UnitStateSucceeded, ProviderRef: &ProviderRef{Name: "backup-global-config"},
				FQDN: "ces.example.com", Domain: "example.com",
			},
			V2Dogus: &V2DogusUnitStatus{
				State: UnitStateSucceeded, ProviderRef: &ProviderRef{Name: "backup-doguv2"},
				Dogus: []DoguVersion{{Name: "redmine", Version: "5.1.0-1"}},
			},
			Dogus: []DoguUnitStatus{
				{Name: "scm", Version: "3.0.0-1", State: UnitStateFailed, Reason: "ProviderBackupFailed"},
				{Name: "jenkins", Version: "2.0.0-1", UnitStatus: UnitStatus{}},
			},
			Misc: &UnitStatus{State: UnitStateSucceeded, ProviderRef: &ProviderRef{Name: "backup-misc"}},
		}

		expected := []IdentifiedUnitStatus{
			{Type: UnitTypeGlobalConfig, UnitStatus: status.GlobalConfig.UnitStatus},
			{Type: UnitTypeDoguV2, UnitStatus: status.V2Dogus.UnitStatus},
			{Type: UnitTypeDogu, Name: "scm", UnitStatus: status.Dogus[0].UnitStatus},
			{Type: UnitTypeDogu, Name: "jenkins", UnitStatus: status.Dogus[1].UnitStatus},
			{Type: UnitTypeMisc, UnitStatus: *status.Misc},
		}
		assert.Equal(t, expected, status.Units())
	})
}

func TestBackupStatus_FailedUnits(t *testing.T) {
	status := BackupStatus{
		GlobalConfig: &GlobalConfigStatus{State: UnitStateMissing, Reason: "NotStarted"},
		Dogus: []DoguUnitStatus{
			{Name: "scm", Version: "3.0.0-1", State: UnitStateFailed},
			{Name: "jenkins", Version: "2.0.0-1", State: UnitStateSucceeded},
			{Name: "ldap", Version: "2.6.0-1", UnitStatus: UnitStatus{}},
		},
	}

	actual := status.FailedUnits()

	require.Len(t, actual, 2)
	assert.Equal(t, UnitIdentifier{Type: UnitTypeGlobalConfig}, actual[0].UnitIdentifier)
	assert.Equal(t, UnitStateMissing, actual[0].State)
	assert.Equal(t, "NotStarted", actual[0].Reason)
	assert.Equal(t, UnitIdentifier{Type: UnitTypeDogu, Name: "scm"}, actual[1].UnitIdentifier)
}

func TestRestoreStatus_FailedUnits(t *testing.T) {
	status := RestoreStatus{Units: []IdentifiedUnitStatus{
		{Type: UnitTypeMisc, State: UnitStateSucceeded},
		{Type: UnitTypeDogu, Name: "scm", State: UnitStateFailed},
		{Type: UnitTypeDoguV2, State: UnitStateMissing},
	}}

	actual := status.FailedUnits()

	require.Len(t, actual, 2)
	assert.Equal(t, UnitIdentifier{Type: UnitTypeDogu, Name: "scm"}, actual[0].UnitIdentifier)
	assert.Equal(t, UnitIdentifier{Type: UnitTypeDoguV2}, actual[1].UnitIdentifier)
}
