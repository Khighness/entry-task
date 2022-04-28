GO          = go
PRODUCT     = web-svc
GOARCH     := amd64
GO111MODULE = on

all: $(shell $(GO) env GOOS)

build:
	env GO111MODULE=${GO111MODULE} GOOS=${GOOS} GOARCH=$(GOARCH) $(GO) build $(EXTFLAGS) -o $(PRODUCT)$(EXT) .

linux: export GOOS=linux
linux: build

darwin: export GOOS=darwin
darwin: EXT=.mach
darwin: build

run:
	make
	./$(PRODUCT).mach

.PHONY: clean
clean:
	@rm -f $(PRODUCT) $(PRODUCT).elf $(PRODUCT).mach
