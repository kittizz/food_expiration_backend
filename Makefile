run:
	go run main.go

migrate:
	go run cmd/migrate/main.go
dev:
	air --build.cmd "go build -o bin/backend main.go" --build.bin "./bin/backend"
build:
	docker build --platform linux/amd64 -t registry.xver.cloud/food_expiration/food_expiration_backend .
	docker push registry.xver.cloud/food_expiration/food_expiration_backend