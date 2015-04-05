EXE=cacher.native
LIBS='core_extended yojson'

all: native

native:
	touch _tags
	corebuild -pkgs $(LIBS) -use-ocamlfind $(EXE)

clean:
	rm -rf $(EXE) _tags _build
