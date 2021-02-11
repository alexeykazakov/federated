QUAY_NAMESPACE ?= ${GO_PACKAGE_ORG_NAME}
TARGET_REGISTRY := quay.io
IMAGE ?= ${TARGET_REGISTRY}/${QUAY_NAMESPACE}/${GO_PACKAGE_REPO_NAME}:${GIT_COMMIT_ID_SHORT}
QUAY_USERNAME ?= ${QUAY_NAMESPACE}
TIMESTAMP := $(shell date +%s)
IMAGE_DEV ?= ${TARGET_REGISTRY}/${QUAY_NAMESPACE}/${GO_PACKAGE_REPO_NAME}:${TIMESTAMP}

.PHONY: docker-image
## Build the docker image locally that can be deployed (only contains bare operator)
docker-image: build
	$(Q)docker build -f build/Dockerfile -t ${IMAGE} .

.PHONY: docker-image-dev
## Build the docker image locally that can be deployed to dev environment
docker-image-dev: build
	$(Q)docker build -f build/Dockerfile -t ${IMAGE_DEV} .


.PHONY: docker-push-dev
## Push the docker dev image to quay.io registry
docker-push-dev: docker-image
ifeq ($(QUAY_NAMESPACE),${GO_PACKAGE_ORG_NAME})
	@echo "#################################################### WARNING ####################################################"
	@echo you are going to push to $(QUAY_NAMESPACE) namespace, make sure you have set QUAY_NAMESPACE variable appropriately
	@echo "#################################################################################################################"
endif
	$(Q)docker push ${IMAGE_DEV}

.PHONY: docker-push
## Push the docker image to quay.io registry
docker-push: docker-image
ifeq ($(QUAY_NAMESPACE),${GO_PACKAGE_ORG_NAME})
	@echo "#################################################### WARNING ####################################################"
	@echo you are going to push to $(QUAY_NAMESPACE) namespace, make sure you have set QUAY_NAMESPACE variable appropriately
	@echo "#################################################################################################################"
endif
	$(Q)docker push ${IMAGE}

.PHONY: docker-login
docker-login:
	@echo "${DOCKER_PASSWORD}" | docker login quay.io -u "${QUAY_USERNAME}" --password-stdin
