list:
	@echo Available targets\:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

run:
	go run main.go -file $(file)

run-race:
	go run -race main.go -file $(file)

test:
	go test -cover -v ./...

test-race:
	go test -cover -race -v ./...

lint:
	golangci-lint -E gofmt,golint run

.PHONY: list run run-race test test-race lint
	
