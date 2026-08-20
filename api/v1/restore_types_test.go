package v1

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRestoreConditionTypes(t *testing.T) {
	assert.Equal(t, "Successful", ConditionSuccessful)
	assert.Equal(t, "Prepared", ConditionPrepared)
	assert.Equal(t, "ProviderRestoreSuccessful", ConditionProviderRestoreSuccessful)
	assert.Equal(t, "WorkloadsRecovered", ConditionWorkloadsRecovered)
}

func TestRestoreStatus_ConditionSerializationRoundTrip(t *testing.T) {
	expected := RestoreStatus{
		Conditions: []metav1.Condition{{
			Type:               ConditionProviderRestoreSuccessful,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: 7,
			LastTransitionTime: metav1.NewTime(time.Date(2026, time.July, 23, 12, 0, 0, 0, time.UTC)),
			Reason:             "VeleroRestoreCompleted",
			Message:            "The provider restore completed successfully.",
		}},
	}

	serialized, err := json.Marshal(expected)
	require.NoError(t, err)

	var actual RestoreStatus
	err = json.Unmarshal(serialized, &actual)

	require.NoError(t, err)
	require.Len(t, actual.Conditions, 1)
	assert.True(t, expected.Conditions[0].LastTransitionTime.Time.Equal(actual.Conditions[0].LastTransitionTime.Time))
	actual.Conditions[0].LastTransitionTime = expected.Conditions[0].LastTransitionTime
	assert.Equal(t, expected, actual)
}

func TestRestore_DeepCopyDoesNotShareConditions(t *testing.T) {
	original := &Restore{
		Status: RestoreStatus{
			Conditions: []metav1.Condition{{
				Type:    ConditionCompleted,
				Status:  metav1.ConditionFalse,
				Reason:  "Pending",
				Message: "Restore is pending.",
			}},
		},
	}

	copied := original.DeepCopy()
	copied.Status.Conditions[0].Message = "changed"
	copied.Status.Conditions = append(copied.Status.Conditions, metav1.Condition{
		Type:   ConditionPrepared,
		Status: metav1.ConditionTrue,
		Reason: "Prepared",
	})

	require.Len(t, original.Status.Conditions, 1)
	assert.Equal(t, "Restore is pending.", original.Status.Conditions[0].Message)
	assert.NotSame(t, &original.Status.Conditions[0], &copied.Status.Conditions[0])
}
