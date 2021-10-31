SNAPSHOTDIR = ./bin
BINARY = inform
.PHONY: mac-binary
mac-binary: clean
	mkdir -p $(SNAPSHOTDIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o $(SNAPSHOTDIR)/$(BINARY) .

.PHONY: clean
clean:
	rm -rf $(SNAPSHOTDIR)
