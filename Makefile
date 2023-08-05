run:
	go run cmd/backend/main.go
dev:
	air --build.cmd "go build -o bin/backend cmd/backend/main.go" --build.bin "./bin/backend"