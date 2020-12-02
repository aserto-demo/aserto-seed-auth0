SHELL 	   := $(shell which bash)

## BOF define block

BINARIES   := aserto-seed-auth0
BINARY     = $(word 1, $@)

PLATFORMS  := linux darwin windows
PLATFORM   = $(word 1, $@)

ROOT_DIR   := $(shell git rev-parse --show-toplevel)
BIN_DIR    := $(ROOT_DIR)/bin
REL_DIR    := $(ROOT_DIR)/release
SRC_DIR    := $(ROOT_DIR)/cmd
PPROF_DIR  := $(ROOT_DIR)/pprof

VERSION    :=`git describe --tags 2>/dev/null`
COMMIT     :=`git rev-parse --short HEAD 2>/dev/null`
DATE       :=`date "+%FT%T%z"`

LDBASE     := github.com/aserto-demo/aserto-seed-auth0/pkg/version
DEV_LDFLAGS:= -ldflags "-X $(LDBASE).ver=${VERSION} -X $(LDBASE).date=${DATE} -X $(LDBASE).commit=${COMMIT}"
REL_LDFLAGS:= -ldflags "-w -s -X $(LDBASE).ver=${VERSION} -X $(LDBASE).date=${DATE} -X $(LDBASE).commit=${COMMIT}"

GOARCH     ?= amd64
GOOS       := $(shell go env GOOS)

LINTER     := $(BIN_DIR)/golangci-lint
LINTVERSION:= v1.32.2

TESTRUNNER := $(BIN_DIR)/gotestsum
TESTVERSION:= v0.6.0

NO_COLOR   :=\033[0m
OK_COLOR   :=\033[32;01m
ERR_COLOR  :=\033[31;01m
WARN_COLOR :=\033[36;01m
ATTN_COLOR :=\033[33;01m

GRPC_SVCS  := 
GTW_SVCS   := edge

## EOF define block

.PHONY: all
all: deps build test

deps:
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@GO111MODULE=on go mod download

.PHONY: dobuild
dobuild:
	@echo -e "$(ATTN_COLOR)==> $@ $(B) GOOS=$(P) GOARCH=$(GOARCH) VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) $(NO_COLOR)"
	@GOOS=$(P) GOARCH=$(GOARCH) GO111MODULE=on go build $(DEV_LDFLAGS) -o $(T)/$(P)-$(GOARCH)/$(B)$(if $(findstring $(P),windows),".exe","") $(SRC_DIR)/$(B)
ifneq ($(P),windows)
	@chmod +x $(T)/$(P)-$(GOARCH)/$(B)
endif

.PHONY: build 
build: $(BIN_DIR)
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@for b in ${BINARIES}; 									\
	do 														\
		$(MAKE) dobuild B=$${b} P=${GOOS} T=${BIN_DIR}; 	\
	done 													

.PHONY: doinstall
doinstall:
	@echo -e "$(ATTN_COLOR)==> $@ GOOS=$(P) GOARCH=$(GOARCH) VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) $(NO_COLOR)"
	@GOOS=$(P) GOARCH=$(GOARCH) GO111MODULE=on go install $(DEV_LDFLAGS) $(SRC_DIR)/$(B)

.PHONY: install
install: 
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@for b in ${BINARIES}; 									\
	do 														\
		$(MAKE) doinstall B=$${b} P=${GOOS}; 			 	\
	done 													

.PHONY: dorelease
dorelease:
	@echo -e "$(ATTN_COLOR)==> $@ $(B) GOOS=$(P) GOARCH=$(GOARCH) VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) $(NO_COLOR)"
	@GOOS=$(P) GOARCH=$(GOARCH) GO111MODULE=on go build $(REL_LDFLAGS) -o $(T)/$(P)-$(GOARCH)/$(B)$(if $(findstring $(P),windows),".exe","") $(SRC_DIR)/$(B)
ifneq ($(P),windows)
	@chmod +x $(T)/$(P)-$(GOARCH)/$(B)
endif

.PHONY: release
release: $(REL_DIR)
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@for b in ${BINARIES}; 									\
	do 														\
		for p in ${PLATFORMS};								\
		do 													\
			$(MAKE) dorelease B=$${b} P=$${p} T=${REL_DIR}; \
		done;												\
	done 													\

$(REL_DIR):
	@echo -e "$(ATTN_COLOR)==> create REL_DIR $(REL_DIR) $(NO_COLOR)"
	@mkdir -p $(REL_DIR)

$(BIN_DIR):
	@echo -e "$(ATTN_COLOR)==> create BIN_DIR $(BIN_DIR) $(NO_COLOR)"
	@mkdir -p $(BIN_DIR)

$(PPROF_DIR):
	@echo -e "$(ATTN_COLOR)==> create PPROF_DIR $(PPROF_DIR) $(NO_COLOR)"
	@mkdir -p $(PPROF_DIR)

$(TESTRUNNER):
	@echo -e "$(ATTN_COLOR)==> get $@  $(NO_COLOR)"
	@GOBIN=$(BIN_DIR) go get -u gotest.tools/gotestsum@$(TESTVERSION)

.PHONY: test 
test: $(TESTRUNNER)
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@CGO_ENABLED=0 $(BIN_DIR)/gotestsum --format short-verbose -- -count=1 -v $(ROOT_DIR)/...

$(LINTER):
	@echo -e "$(ATTN_COLOR)==> get $@  $(NO_COLOR)"
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(LINTVERSION)
 
.PHONY: lint
lint: $(LINTER)
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@CGO_ENABLED=0 $(LINTER) run 
	@echo -e "$(NO_COLOR)\c"

.PHONY: clean
clean:
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@rm -rf $(BIN_DIR)
	@rm -rf $(REL_DIR)
	@rm -rf $(PPROF_DIR)
	@go clean

.PHONY: bench
bench: $(PPROF_DIR)
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@go test -run='^\$$' -bench=. -cpuprofile=$(PPROF_DIR)/cpu.pprof -benchmem -memprofile=$(PPROF_DIR)/mem.pprof -trace $(PPROF_DIR)/trace.out ./cmd
	@go tool trace -pprof=net $(PPROF_DIR)/trace.out > $(PPROF_DIR)/net.pprof
	@go tool trace -pprof=sync $(PPROF_DIR)/trace.out > $(PPROF_DIR)/sync.pprof
	@go tool trace -pprof=syscall $(PPROF_DIR)/trace.out > $(PPROF_DIR)/syscall.pprof
	@go tool trace -pprof=sched $(PPROF_DIR)/trace.out > $(PPROF_DIR)/sched.pprof

.PHONY: memprofile
memprofile:
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@go tool pprof -http=:8080 ./pprof/mem.pprof

.PHONY: cpuprofile
cpuprofile:
	@echo -e "$(ATTN_COLOR)==> $@ $(NO_COLOR)"
	@go tool pprof -http=:8080 ./pprof/cpu.pprof
