before.build:
	go mod download && go mod vendor

build.JSextractor:
	@echo "build in ${PWD}";go build ./cmd/jse/jse.go