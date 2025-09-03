PROJECT_NAME=k8s-backup-lib
ARTIFACT_ID=k8s-backup-operator-crd
APPEND_CRD_SUFFIX=false
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}

VERSION=1.5.0
GOTAG=1.24.1
MAKEFILES_VERSION=10.2.0

include build/make/variables.mk
include build/make/self-update.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/k8s-controller.mk
include build/make/release.mk

CRD_BACKUP_SOURCE = ${HELM_CRD_SOURCE_DIR}/templates/k8s.cloudogu.com_backups.yaml
CRD_RESTORE_SOURCE = ${HELM_CRD_SOURCE_DIR}/templates/k8s.cloudogu.com_restores.yaml
CRD_SCHEDULE_SOURCE = ${HELM_CRD_SOURCE_DIR}/templates/k8s.cloudogu.com_backupschedules.yaml
CRD_POST_MANIFEST_TARGETS = crd-add-labels crd-add-backup-labels

.PHONY: crd-add-backup-labels
crd-add-backup-labels: $(BINARY_YQ)
	@echo "Adding backup label to CRDs..."
	@for file in ${HELM_CRD_SOURCE_DIR}/templates/*.yaml ; do \
		$(BINARY_YQ) -i e ".metadata.labels.\"k8s.cloudogu.com/part-of\" = \"backup\"" $${file} ;\
	done

