run:
	@trap 'kill 0' EXIT INT TERM; \
	echo "Building and starting API (:3000)..."; \
	(cd cars-api && make build && make run) & \
	echo "Waiting for API to be ready..."; \
	while ! (echo > /dev/tcp/localhost/3000) 2>/dev/null; do sleep 0.5; done; \
	echo "API is up. Starting backend (:8080)..."; \
	(cd go-backend && go run main.go) & \
	wait