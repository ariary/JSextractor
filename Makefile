before.build:
	go mod download && go mod vendor

build.js-extractor:
	@echo "build in ${PWD}";go build js-extractor.go