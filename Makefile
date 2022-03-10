BINARYDIR = ./bin
BINARY = inform

.PHONY: linux-binary
linux-binary: clean
	mkdir -p $(BINARYDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(BINARYDIR)/$(BINARY) .

.PHONY: mac-binary
mac-binary: clean
	mkdir -p $(BINARYDIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o $(BINARYDIR)/$(BINARY) .

.PHONY: clean
clean:
	rm -rf $(BINARYDIR)
