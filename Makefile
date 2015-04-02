GO=gom
GOFLAGS=-race
SOURCES=*.go
EXE=vx-cache-tool

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
# use the rest as arguments for "run"
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# ...and turn them into do-nothing targets
$(eval $(RUN_ARGS):;@:)
endif

all: bundle $(EXE)

bundle:
	gom install

$(EXE): $(SOURCES)
	$(GO) build $(GOFLAGS)

clean:
	rm -f $(EXE)

clean_bundle:
	rm -rf _vendor

run: $(EXE)
	./$(EXE) $(RUN_ARGS)

gom:
	go get github.com/mattn/gom
