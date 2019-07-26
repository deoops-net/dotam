VERSION=v1.0.0

.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test -v -count=1

.PHONY: run-dev-build
run-dev-build:
	LOG_LEVEL=debug go run . build

.PHONY: run-dev-init
run-dev-init:
	LOG_LEVEL=debug go run . init