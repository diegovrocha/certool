PREFIX ?= /usr/local
BINARY = certui
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/diegovrocha/certui/internal/ui.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/certui

install: build
	@mkdir -p $(PREFIX)/bin
	@cp $(BINARY) $(PREFIX)/bin/$(BINARY)
	@chmod +x $(PREFIX)/bin/$(BINARY)
	@echo "✔ certui instalado em $(PREFIX)/bin/$(BINARY)"

uninstall:
	@rm -f $(PREFIX)/bin/$(BINARY)
	@echo "✔ certui removido"

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test ./... -count=1

.PHONY: build install uninstall run clean test
