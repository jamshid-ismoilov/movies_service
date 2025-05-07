APP_NAME = movies_service

# build application
build:
	go build -o $(APP_NAME) .

# running migrations
migrate-up:
	sql-migrate up -config=dbconfig.yml -env=development

# generate swagger documentation
swagger:
	swag init -g main.go -o docs

# running unit tests
test:
	go test ./...

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run --rm -p 8080:8080 \
		-e DB_HOST=localhost \
		-e DB_PORT=5432 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=postgres \
		-e DB_NAME=movies_db \
		-e JWT_SECRET=secret \
		$(APP_NAME)
