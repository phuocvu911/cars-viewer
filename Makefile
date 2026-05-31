.PHONY: run

api-run:
	cd cars-api && node main.js

run:
	cd go-backend && go run main.go

