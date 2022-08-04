build:
	go build -o ./movie_app cmd/api/*.go
start:
	migrate -database ${POSTGRES_URI} -path ./migrations up
	./movie_app 