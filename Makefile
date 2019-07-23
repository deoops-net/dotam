VERSION=v1.0.0

.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test -v -count=1

.PHONY: run-dev
run-dev:
	LOG_LEVEL=debug go run . build