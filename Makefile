BIN := dossier
CMD := ./cmd/dossier

.PHONY: build install clean

build:
	go build -o $(BIN) $(CMD)

install:
	go install $(CMD)

clean:
	rm -f $(BIN)
