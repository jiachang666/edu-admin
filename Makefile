APP_NAME := edu-admin

.PHONY: run test tidy web-dev

run:
	go run ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

web-dev:
	cd web/admin && npm run dev
