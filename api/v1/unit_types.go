package v1

// Labels on provider backups (e.g. Velero backups) which assign them to a backup unit and to the Cloudogu backup.
const (
	// BackupLabelKey contains the name of the Cloudogu backup which the provider backup belongs to.
	BackupLabelKey = "k8s.cloudogu.com/backup"
	// BackupTypeLabelKey contains the UnitType of the provider backup.
	BackupTypeLabelKey = "k8s.cloudogu.com/backup-type"
	// BackupCoverageLabelKey contains the BackupCoverage of the Cloudogu backup.
	BackupCoverageLabelKey = "k8s.cloudogu.com/backup-coverage"
	// DoguNameLabelKey contains the name of the dogu of a unit of type UnitTypeDogu.
	DoguNameLabelKey = "k8s.cloudogu.com/dogu.name"
	// DoguVersionLabelKey contains the version of the dogu of a unit of type UnitTypeDogu.
	DoguVersionLabelKey = "k8s.cloudogu.com/dogu.version"
)

// Labels on provider restores (e.g. Velero restores) which reference the Cloudogu restore.
const (
	RestoreSourceNameLabelKey = "k8s.cloudogu.com/restore-source-name"
	RestoreSourceUIDLabelKey  = "k8s.cloudogu.com/restore-source-uid"
)

// Annotations on provider backups.
const (
	// InstalledDogusV2AnnotationKey contains the serialized list of V2 dogus with their versions in the V2 unit.
	InstalledDogusV2AnnotationKey = "k8s.cloudogu.com/backup-installed-dogus-v2"
	// BlueprintIDAnnotationKey contains the ID of the blueprint which was installed at the time of the backup.
	BlueprintIDAnnotationKey = "backup.cloudogu.com/blueprintId"
	// UnitsAnnotationKey contains the serialized list of all expected units of the Cloudogu backup.
	// It is used to detect missing units when the Cloudogu backup is reconstructed from the provider backups.
	UnitsAnnotationKey = "k8s.cloudogu.com/backup-units"
	// FailureModeAnnotationKey contains the FailureMode of the Cloudogu backup.
	FailureModeAnnotationKey = "k8s.cloudogu.com/backup-failure-mode"
)

// UnitType defines the type of backup unit. Each unit is backed up by its own provider backup.
type UnitType string

const (
	// UnitTypeDogu contains all resources of a single V3 dogu.
	UnitTypeDogu UnitType = "dogu"
	// UnitTypeDoguV2 contains all V2 dogus.
	UnitTypeDoguV2 UnitType = "doguv2"
	// UnitTypeGlobalConfig contains the global config.
	UnitTypeGlobalConfig UnitType = "global-config"
	// UnitTypeMisc contains all resources with the label k8s.cloudogu.com/backup-scope which belong to no other unit.
	UnitTypeMisc UnitType = "misc"
)

// FailureMode defines how failed units are handled.
type FailureMode string

const (
	// FailureModeFailFast aborts the whole operation as soon as one unit fails.
	FailureModeFailFast FailureMode = "failfast"
	// FailureModePartiallyAllowed continues with the remaining units if one unit fails.
	FailureModePartiallyAllowed FailureMode = "partiallyAllowed"
)

// BackupCoverage defines whether a backup contains the whole system.
type BackupCoverage string

const (
	// BackupCoverageFull marks a backup of the whole system.
	BackupCoverageFull BackupCoverage = "full"
	// BackupCoveragePartial marks a backup which excludes parts of the system.
	BackupCoveragePartial BackupCoverage = "partial"
)

// UnitOptions defines which units are processed and how failures are handled.
type UnitOptions struct {
	// IgnoreGlobalConfig excludes the global-config unit.
	IgnoreGlobalConfig bool `json:"ignoreGlobalConfig,omitempty"`
	// IgnoreV2Dogus excludes the unit containing all V2 dogus.
	IgnoreV2Dogus bool `json:"ignoreV2Dogus,omitempty"`
	// IgnoreMisc excludes the misc unit containing resources with the label k8s.cloudogu.com/backup-scope.
	IgnoreMisc bool `json:"ignoreMisc,omitempty"`
	// FailureMode defines how failed units are handled.
	// +kubebuilder:validation:Enum=failfast;partiallyAllowed
	// +kubebuilder:default=partiallyAllowed
	FailureMode FailureMode `json:"failureMode,omitempty"`
}

// UnitIdentifier identifies a unit of a Cloudogu backup. A list of it is serialized into UnitsAnnotationKey.
type UnitIdentifier struct {
	// Type of the unit.
	Type UnitType `json:"type"`
	// Name of the dogu. Only set for units of type dogu.
	// +optional
	Name string `json:"name,omitempty"`
}

// ProviderRef references the object of the backup provider (e.g. a Velero backup or restore).
type ProviderRef struct {
	// Name of the provider object.
	Name string `json:"name"`
}

// UnitStatus describes a unit which is not specific to a dogu.
type UnitStatus struct {
	// ProviderRef references the provider backup of the unit.
	ProviderRef ProviderRef `json:"providerRef"`
}

// DoguVersion identifies a dogu in a specific version.
type DoguVersion struct {
	// Name of the dogu without namespace.
	Name string `json:"name"`
	// Version of the dogu.
	Version string `json:"version"`
}

// DoguUnitStatus describes the unit of a single V3 dogu.
type DoguUnitStatus struct {
	DoguVersion `json:",inline"`
	// ProviderRef references the provider backup of the unit.
	ProviderRef ProviderRef `json:"providerRef"`
}

// V2DogusUnitStatus describes the unit containing all V2 dogus.
type V2DogusUnitStatus struct {
	// ProviderRef references the provider backup of the unit.
	ProviderRef ProviderRef `json:"providerRef"`
	// Dogus contains all V2 dogus with their versions in this unit.
	// +listType=map
	// +listMapKey=name
	// +optional
	Dogus []DoguVersion `json:"dogus,omitempty"`
}

// ComponentStatus describes a component at the time of the backup.
type ComponentStatus struct {
	// Name of the component.
	Name string `json:"name"`
	// Version of the component.
	Version string `json:"version"`
}

// FailedUnit describes a unit whose provider operation failed.
type FailedUnit struct {
	// Type of the unit.
	Type UnitType `json:"type"`
	// Name of the dogu. Only set for units of type dogu.
	// +optional
	Name string `json:"name,omitempty"`
	// ProviderRef references the failed provider object.
	ProviderRef ProviderRef `json:"providerRef"`
}
