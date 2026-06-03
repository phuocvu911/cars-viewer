.PHONY: run

api:
	cd cars-api && make build && make run

run:
	cd go-backend && go run main.go

