.PHONY: clean
## cleans up
clean:
	@cd "$(GOPATH)/src/github.com/alexeykazakov/proxy" && \
		rm -rf $(COV_DIR) && \
		rm -rf $(OUT_DIR) && \
		rm -rf ${V_FLAG} ./vendor
	$(Q)go clean ${X_FLAG} ./...
