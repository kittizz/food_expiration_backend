run:
	go run main.go
dev:
	air --build.cmd "go build -o bin/backend main.go" --build.bin "./bin/backend"