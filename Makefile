tidy:
	go mod tidy

run:
	go run main.go

build:
	go build -o gorawler main.go

.PHONY: test
test:
	go test --cover ./...