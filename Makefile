GO=gom
GOFLAGS=-race
SOURCES=*.go
EXE=vx-cache-tool

all: bundle $(EXE)

bundle:
	gom install

$(EXE): $(SOURCES)
	$(GO) build $(GOFLAGS)

clean:
	rm -f $(EXE)

clean_bundle:
	rm -rf _vendor

gom:
	go get github.com/mattn/gom
