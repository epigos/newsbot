.PHONY: install lint test test-cover build deploy-dev deploy-prod


setup:
	go get -v -t `go list ./... | grep newsbot`
	go get github.com/jstemmer/go-junit-report

lint: setup
	go vet -v `go list ./... | grep newsbot`

test: lint
	go test -v -coverprofile=./test-reports/cover.out `go list ./... | grep newsbot/`

test-cover: lint
	mkdir -pv ./test-reports
	go test -v -coverprofile=./test-reports/cover.out `go list ./... | grep newsbot/` 2>&1 | tee go-junit-report > ./test-reports/junit.xml

build: lint
	go install newsbot

build-docker:
	bash -c "source ./deploy/vars.sh dev && ./deploy/image.sh"

dev:
	bash -c "source ./deploy/vars.sh dev && ./deploy/deploy.sh"

tag:
	bash -c "source ./deploy/vars.sh prod && ./deploy/tag.sh"

release:
	bash -c "source ./deploy/vars.sh prod && ./deploy/deploy.sh"