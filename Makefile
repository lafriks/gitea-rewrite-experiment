DIST := dist
IMPORT := code.gitea.io/gitea

GO ?= go
SED_INPLACE := sed -i

ifeq ($(OS), Windows_NT)
	EXECUTABLE := gitea.exe
else
	EXECUTABLE := gitea
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		SED_INPLACE := sed -i ''
	endif
endif

BINDATA := modules/{options,public,templates}/bindata.go
GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*" ! -path "*/bindata.go")
GOFMT ?= gofmt -s

GOFLAGS := -i -v
EXTRA_GOFLAGS ?=

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
	GITEA_VERSION := $(VERSION)
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst release/v,,$(DRONE_BRANCH))
	else
		VERSION ?= master
	endif
	GITEA_VERSION := $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')
endif

LDFLAGS := -X "main.Version=$(GITEA_VERSION)" -X "main.Tags=$(TAGS)"

PACKAGES ?= $(filter-out code.gitea.io/gitea/integrations,$(shell $(GO) list ./... | grep -v /vendor/))
SOURCES ?= $(shell find . -name "*.go" -type f)

TAGS ?=

TMPDIR := $(shell mktemp -d 2>/dev/null || mktemp -d -t 'gitea-temp')

SWAGGER_SPEC := templates/swagger/v1_json.tmpl
SWAGGER_SPEC_S_TMPL := s|"basePath":\s*"/api/v1"|"basePath": "{{AppSubUrl}}/api/v1"|g
SWAGGER_SPEC_S_JSON := s|"basePath":\s*"{{AppSubUrl}}/api/v1"|"basePath": "/api/v1"|g

TEST_MYSQL_HOST ?= mysql:3306
TEST_MYSQL_DBNAME ?= testgitea
TEST_MYSQL_USERNAME ?= root
TEST_MYSQL_PASSWORD ?=
TEST_PGSQL_HOST ?= pgsql:5432
TEST_PGSQL_DBNAME ?= testgitea
TEST_PGSQL_USERNAME ?= postgres
TEST_PGSQL_PASSWORD ?= postgres

ifeq ($(OS), Windows_NT)
	EXECUTABLE := gitea.exe
else
	EXECUTABLE := gitea
endif

# $(call strip-suffix,filename)
strip-suffix = $(firstword $(subst ., ,$(1)))

.PHONY: all
all: build

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: generate
generate:
	$(GO) generate $(PACKAGES)

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	$(GO) build $(GOFLAGS) $(EXTRA_GOFLAGS) -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)' -o $@
