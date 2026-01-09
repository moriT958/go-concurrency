.PHONY: bench
bench:
	go test -bench . ./exercise/...


.PHONY: dev
dev: cmd/api/main.go
	go run cmd/api/main.go

.PHONY: test
test:
	go test -short -v ./internal/...