all:
	go build -o $(CURDIR)/bin/ $(CURDIR)/cmd/...

run:
	go run $(CURDIR)/cmd/app/main.go

run-example:
	go run $(CURDIR)/cmd/example/main.go
