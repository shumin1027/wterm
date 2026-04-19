PWD=$(shell pwd)
DIST=$(shell pwd)/bin
DATE=$(shell date --iso-8601=seconds)

.PHONY: build-frontend-react
build-frontend-react:
	@echo ">> building frontend whit react"
	cd frontend/React/ && pnpm install && pnpm run build

.PHONY: build-frontend-vanilla
build-frontend-vanilla:
	@echo ">> building frontend whit vanilla"
	cd frontend/Vanilla/  && pnpm install && pnpm run build

.PHONY: build
build: build-frontend-react
	go mod tidy
	go mod vendor
	go build -mod=vendor -o ${DIST}/wterm
