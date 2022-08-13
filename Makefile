.PHONY: build
build:
	@echo Building...
	@echo
	@CGO_ENABLED=0 go build -ldflags '-w -s -extldflags "-static"' -tags netgo -o bin/irsa-emu-webhook ./cmd

.PHONY: ginkgo
ginkgo:
ifeq (,$(shell which ginkgo 2>/dev/null))
	@{ \
	set -e ;\
	go install github.com/onsi/ginkgo/ginkgo@v1.16.5 ;\
	}
endif

.PHONY: test-suite
test-suite: ginkgo
ifeq (,$(pkg))
	@echo specify a package. e.g. pkg=pkg/runner
	@false
endif
ifeq (,$(wildcard $(pkg)))
	@echo package ${pkg} does not exist.
else
	@{ \
	set -e ;\
	cd ${pkg} ;\
	ginkgo bootstrap -internal ;\
	}
endif

.PHONY: test-template
test-template: ginkgo
ifeq (,$(src))
	@echo specify a source file. e.g. src=pkg/webhook/server.go
	@false
endif
ifeq (,$(wildcard $(src)))
	@echo source ${src} does not exist.
else
	@{ \
	set -e ;\
	cd $$(dirname ${src}) ;\
	ginkgo generate -internal $$(basename ${src}) ;\
	}
endif

.PHONY: mockgen
mockgen:
ifeq (,$(shell which mockgen 2>/dev/null))
	@{ \
	set -e ;\
	go install github.com/golang/mock/mockgen@v1.6.0 ;\
	}
endif

.PHONY: mock
mock: model mockgen
	@echo Generating mocks...
	@echo
	@for GO_FILE in $$(find ./pkg -name "*.go" -not -name "*_test.go" -not -name "doc.go"); do\
		MOCK_DIR=mock/$$(dirname $$(dirname $$GO_FILE))/mock_$$(basename $$(dirname $$GO_FILE)) ;\
		mkdir -p $${MOCK_DIR} ;\
		mockgen -source=$$GO_FILE -destination $${MOCK_DIR}/$$(basename $$GO_FILE) ;\
	done
	@# Remove emply mock files.
	@rm -f $$(find mock/ -name "*.go" | xargs grep -iL "func ")

.PHONY: test mock
test: ginkgo mock
	@echo Running unit tests...
	@echo
	@ginkgo -r -cover
