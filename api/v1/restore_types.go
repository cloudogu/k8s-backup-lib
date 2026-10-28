package v1

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deprecated: use Conditions instead.
const (
	RestoreStatusNew        = ""
	RestoreStatusInProgress = "in progress"
	RestoreStatusFailed     = "failed"
	RestoreStatusCompleted  = "completed"
	RestoreStatusDeleting   = "deleting"
)

const (
	ConditionWorkloadsRecovered = "WorkloadsRecovered"
)

const RestoreFinalizer = "cloudogu-restore-finalizer"

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RestoreSpec defines the desired state of Restore
type RestoreSpec struct {
	// BackupName references the backup that should be restored.
	BackupName string `json:"backupName,omitempty"`
	// Provider defines the backup provider which should be used for the restore.
	// +kubebuilder:validation:Enum=velero
	Provider Provider `json:"provider,omitempty"`
	// ExcludedDogus contains the names of V3 dogus which are not restored.
	// V2 dogus can only be excluded as a whole with Options.IgnoreV2Dogus.
	// +listType=set
	// +kubebuilder:validation:MaxItems=256
	// +kubebuilder:validation:items:MaxLength=63
	// +kubebuilder:validation:items:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	// +optional
	ExcludedDogus []string `json:"excludedDogus,omitempty"`
	// ExcludedComponents contains the names of components whose resources are not restored.
	// +listType=set
	// +kubebuilder:validation:MaxItems=256
	// +kubebuilder:validation:items:MaxLength=63
	// +kubebuilder:validation:items:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	// +optional
	ExcludedComponents []string `json:"excludedComponents,omitempty"`
	// Options defines which units are restored and how failed units are handled.
	// +optional
	Options UnitOptions `json:"options,omitempty"`
}

// RestoreStatus defines the observed state of Restore
type RestoreStatus struct {
	// Status represents the state of the backup.
	Status string `json:"status,omitempty"`
	// RequeueTimeNanos contains the time in nanoseconds to wait until the next requeue.
	RequeueTimeNanos time.Duration `json:"requeueTimeNanos,omitempty"`
	// Conditions represent the latest available observations of the Restore's state.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// FailedUnits contains the units whose provider restore failed. ProviderRef references the provider restore.
	// +optional
	FailedUnits []FailedUnit `json:"failedUnits,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels=app=ces;app.kubernetes.io/name=k8s-backup-operator;k8s.cloudogu.com/part-of=backup
// +kubebuilder:printcolumn:name="Backup name",type="string",JSONPath=".spec.backupName",description="The backup name for the restore"
// +kubebuilder:printcolumn:name="Succeeded",type="string",JSONPath=".status.conditions[?(@.type=='Succeeded')].status",description="Whether the restore is already completed successfully. False shows a failure."
// +kubebuilder:printcolumn:name="Prepared",type="string",JSONPath=".status.conditions[?(@.type=='Prepared')].status",description="Whether all preparations for the restore are completed"
// +kubebuilder:printcolumn:name="ProviderSucceeded",type="string",JSONPath=".status.conditions[?(@.type=='ProviderSucceeded')].status",description="Whether the backup provider finished the restore process successfully. False shows a failure."
// +kubebuilder:printcolumn:name="WorkloadsRecovered",type="string",JSONPath=".status.conditions[?(@.type=='WorkloadsRecovered')].status",description="Whether the post restore workload operations are completed"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of the resource"

// Restore is the Schema for the restores API
type Restore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of Restore
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="spec is immutable"
	Spec RestoreSpec `json:"spec,omitempty"`
	// Status defines the observed state of Restore
	Status RestoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RestoreList contains a list of Restore
type RestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Restore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Restore{}, &RestoreList{})
}

// GetFieldSelectorWithName return the field selector with the metadata.name field.
func (r *Restore) GetFieldSelectorWithName() string {
	return fmt.Sprintf("metadata.name=%s", r.Name)
}

// GetStatus return the requeueable status.
func (r *Restore) GetStatus() RequeueableStatus {
	return r.Status
}

// GetStatus return the status from the status object.
func (rs RestoreStatus) GetStatus() string {
	return rs.Status
}

// GetRequeueTimeNanos returns the requeue time in nano seconds.
func (rs RestoreStatus) GetRequeueTimeNanos() time.Duration {
	return rs.RequeueTimeNanos
}
