# noobaa-operator git tag for OB/OBC CRDs (keep in sync with EmbeddedNoobaaOperatorVersion in pkg/noobaa/crds.go).
NOOBAA_OPERATOR_VERSION ?= v5.21.0

build: gen-noobaa-crds
	find . -type f -name '*.go' -print0 | xargs -0 gofmt -w
	@echo "=> Building bin/odf (GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH))"
	env GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o bin/odf cmd/odf/main.go
	@echo "=> Build finished: bin/odf"

gen-noobaa-crds:
	@mkdir -p pkg/noobaa/crds
	@curl -fsSL https://raw.githubusercontent.com/noobaa/noobaa-operator/$(NOOBAA_OPERATOR_VERSION)/deploy/obc/objectbucket.io_objectbuckets_crd.yaml \
		-o pkg/noobaa/crds/objectbucket.io_objectbuckets_crd.yaml
	@curl -fsSL https://raw.githubusercontent.com/noobaa/noobaa-operator/$(NOOBAA_OPERATOR_VERSION)/deploy/obc/objectbucket.io_objectbucketclaims_crd.yaml \
		-o pkg/noobaa/crds/objectbucket.io_objectbucketclaims_crd.yaml

clean:
	@rm -f bin/odf

test:
	@echo "running unit tests"
	go test ./...

help:
	@echo "build : Create go binary."
	@echo "test  : Runs unit tests"
	@echo "clean : Remove go binary file."
	@echo "gen-noobaa-crds : Download OB/OBC CRD YAML from noobaa-operator."
