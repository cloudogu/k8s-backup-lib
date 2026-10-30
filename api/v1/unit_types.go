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
	// FQDNAnnotationKey contains the FQDN of the system at the time of the backup.
	// It is only set on the provider backup of the global-config unit.
	FQDNAnnotationKey = "k8s.cloudogu.com/backup-fqdn"
	// DomainAnnotationKey contains the domain of the system at the time of the backup.
	// It is only set on the provider backup of the global-config unit.
	DomainAnnotationKey = "k8s.cloudogu.com/backup-domain"
	// ComponentsAnnotationKey contains the serialized list of components with name and version that are part of the misc unit.
	ComponentsAnnotationKey = "k8s.cloudogu.com/backup-components"
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

// BackupCoverage defines whether a backup contains the whole system.
type BackupCoverage string

const (
	// BackupCoverageFull marks a backup of the whole system.
	BackupCoverageFull BackupCoverage = "full"
	// BackupCoveragePartial marks a backup which excludes parts of the system.
	BackupCoveragePartial BackupCoverage = "partial"
)

// UnitOptions defines which units are processed.
type UnitOptions struct {
	// IgnoreGlobalConfig excludes the global-config unit.
	IgnoreGlobalConfig bool `json:"ignoreGlobalConfig,omitempty"`
	// IgnoreV2Dogus excludes the unit containing all V2 dogus.
	IgnoreV2Dogus bool `json:"ignoreV2Dogus,omitempty"`
	// IgnoreMisc excludes the misc unit containing resources with the label k8s.cloudogu.com/backup-scope.
	IgnoreMisc bool `json:"ignoreMisc,omitempty"`
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

// UnitState defines the terminal state of a single unit.
type UnitState string

const (
	// UnitStateSucceeded marks a unit whose provider operation succeeded.
	UnitStateSucceeded UnitState = "Succeeded"
	// UnitStateFailed marks a unit whose provider operation failed terminally.
	UnitStateFailed UnitState = "Failed"
	// UnitStateMissing marks an expected unit without a provider object, e.g. because it was never started,
	// was deleted after a cancellation or did not arrive during a sync.
	UnitStateMissing UnitState = "Missing"
)

// UnitStatus describes the state of a unit.
type UnitStatus struct {
	// State of the unit. It is empty as long as the unit has no terminal result.
	// +kubebuilder:validation:Enum=Succeeded;Failed;Missing
	// +optional
	State UnitState `json:"state,omitempty"`
	// Reason contains a programmatic identifier in CamelCase for the state, e.g. ProviderBackupFailed.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message contains a human-readable explanation of the state.
	// +optional
	Message string `json:"message,omitempty"`
	// ProviderRef references the provider object of the unit (e.g. a Velero backup or restore).
	// It is set for the states Succeeded and Failed and empty for the state Missing.
	// +optional
	ProviderRef *ProviderRef `json:"providerRef,omitempty"`
}

// IdentifiedUnitStatus describes the state of a unit together with its identifier.
type IdentifiedUnitStatus struct {
	UnitIdentifier `json:",inline"`
	UnitStatus     `json:",inline"`
}

// GlobalConfigStatus describes the global-config unit.
type GlobalConfigStatus struct {
	UnitStatus `json:",inline"`
	// FQDN of the system at the time of the backup.
	FQDN string `json:"fqdn"`
	// Domain of the system at the time of the backup.
	Domain string `json:"domain"`
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
	UnitStatus  `json:",inline"`
}

// V2DogusUnitStatus describes the unit containing all V2 dogus.
type V2DogusUnitStatus struct {
	UnitStatus `json:",inline"`
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

func failedUnits(units []IdentifiedUnitStatus) []IdentifiedUnitStatus {
	var failed []IdentifiedUnitStatus
	for _, unit := range units {
		if unit.State == UnitStateFailed || unit.State == UnitStateMissing {
			failed = append(failed, unit)
		}
	}
	return failed
}
