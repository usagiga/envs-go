xxx:
	@echo "Please select optimal option."

build:
	@go build -o envs-go .

clean:
	@rm -f ./envs-go

run:
	@go run .

test:
	@go test -v "./..."
