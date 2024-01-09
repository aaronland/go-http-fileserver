GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/fileserver cmd/fileserver/main.go

dist-build:
	OS=darwin make dist-os
	OS=windows make dist-os
	OS=linux make dist-os

dist-os:
	@echo "build tools for $(OS)"
	mkdir -p dist/$(OS)
	GOOS=$(OS) GOARCH=386 go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o dist/$(OS)/fileserver cmd/fileserver/main.go
	chmod +x dist/$(OS)/fileserver
