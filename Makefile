BIN := dossier
CMD := ./cmd/dossier

.PHONY: build install test lint fmt clean

build:
	go build -o $(BIN) $(CMD)

install:
	go install $(CMD)

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

fmt:
	goimports -w .

clean:
	rm -f $(BIN)
