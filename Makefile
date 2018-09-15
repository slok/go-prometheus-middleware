
UNIT_TEST_CMD := go test -v
INTEGRATION_TEST_CMD := go test -v -tags='integration'

.PHONY: default
default: test

.PHONY: unit-test
unit-test:
	$(UNIT_TEST_CMD)
.PHONY: integration-test
integration-test:
	$(INTEGRATION_TEST_CMD)
.PHONY: test
test: integration-test

.PHONY: ci
ci: test

.PHONY: godoc
godoc: 
	godoc -http=":6060"
