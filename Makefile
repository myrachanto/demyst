build:
	@go build -o demyst

run:
	@go run .

test:
	@go test -v ./...

testCover:
	@go test -v ./... -cover

swagger:
	@"$HOME/go/bin/swag init -g ./src/routes/routes.go"

dockerize:
	@docker build -t demyst:latest .

dockerrun:
	@docker run --name demyst -p 5000:5000 demyst:latest