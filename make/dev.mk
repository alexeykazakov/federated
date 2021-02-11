MINISHIFT_IP?=$(shell minishift ip)
MINISHIFT_HOSTNAME=minishift.local
MINISHIFT_HOSTNAME_REGEX='minishift\.local'
ETC_HOSTS=/etc/hosts

APP_NAMESPACE ?= $(NAMESPACE)
NAMESPACE ?= "proxy"

.PHONY: create-namespace
## Create the test namespace
create-namespace:
	$(Q)-echo "Creating Namespace"
	$(Q)-oc new-project $(NAMESPACE)
	$(Q)-echo "Switching to the namespace $(NAMESPACE)"
	$(Q)-oc project $(NAMESPACE)

.PHONY: use-namespace
## Log in as system:admin and enter the test namespace
use-namespace: login-as-admin
	$(Q)-echo "Using to the namespace $(NAMESPACE)"
	$(Q)-oc project $(NAMESPACE)

.PHONY: clean-namespace
## Delete the test namespace
clean-namespace:
	$(Q)-echo "Deleting Namespace"
	$(Q)-oc delete project $(NAMESPACE)

.PHONY: deploy
## Deploy DevCluster
deploy: create-namespace build docker-image-dev docker-push-dev apply-resources print-route

.PHONY: apply-resources
## Apply DevCluster resources
apply-resources:
	$(Q)oc process -f ./deploy/proxy.yaml \
        -p IMAGE=${IMAGE_DEV} \
        -p NAMESPACE=${NAMESPACE} \
        | oc apply -f -

.PHONY: print-route
print-route:
	@echo "------------------------------------------------------------------"
	@echo "Deployment complete! Waiting for the proxy service route."
	@echo -n "."
	@while [[ -z `oc get routes proxy -n ${NAMESPACE} 2>/dev/null` ]]; do \
		if [[ $${NEXT_WAIT_TIME} -eq 100 ]]; then \
            echo ""; \
            echo "The timeout of waiting for the service route has been reached. Try to run 'make print-route' later or check the deployment logs"; \
            exit 1; \
		fi; \
		echo -n "."; \
		sleep 1; \
	done
	@echo ""
	$(eval ROUTE = $(shell oc get routes proxy -n ${NAMESPACE} -o=jsonpath='{.spec.host}'))
	@echo Access the Landing Page here: https://${ROUTE}
	@echo "------------------------------------------------------------------"