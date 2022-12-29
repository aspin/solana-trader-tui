all:
	go build -o $(CURDIR)/bin/ $(CURDIR)/cmd/...

run:
	go run $(CURDIR)/cmd/app/main.go
