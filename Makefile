include .env
VARS:=$(shell sed -ne 's/ *\#.*$$//; /./ s/=.*$$// p' .env )
$(foreach v,$(VARS),$(eval $(shell echo export $(v)="$($(v))")))
ifndef PROJECT_NAME
	@echo "no"
endif

default: local

.PHONY: dep  prod docker-image local docker-run docker-up docker-down clear manifests
dep:
	@echo "download go mod dependencies"
	@go mod tidy
	@go mod download

pb:
	@echo "generate protobuf file in local machine"
	@sh ./build/pb_gen.sh


dk-pb:
	@echo "generate protobuf file in docker"
	@sh ./build/pb_gen.sh -d yes

prod:
	@echo "build prod version for CI or local"
	sh ./build/build.sh -a linux -e prod
	docker build -t feimumoke/$(PROJECT_NAME) \
		     --build-arg PROJECT_NAME=$(PROJECT_NAME) \
		     --build-arg GRPC_PORT=$(GRPC_PORT) \
		     --build-arg HTTP_PORT=$(HTTP_PORT) .

docker-image:
	@echo "build docker image local"
	sh ./build/build.sh -a linux -e dev
	docker build -t feimumoke/$(PROJECT_NAME) \
		     --build-arg PROJECT_NAME=$(PROJECT_NAME) \
		     --build-arg GRPC_PORT=$(GRPC_PORT) \
		     --build-arg HTTP_PORT=$(HTTP_PORT) .

local: build

ci:
	@echo "Set the image tag for manifests"
	sh ./build/manifests_update_image.sh "$(PROJECT_NAME)" "$(VERSION)"
	@echo "ci build in standard env"
	sh ./build/build.sh -v enable
	@echo 'Insert VERSION into generated deployment.yaml for each pipeline'
	sh ./build/manifests_version_subst.sh

manifests:
	@echo "Generate manifests"
	sh ./build/kustomize.sh
	sh ./build/manifests_rename.sh

unit-test:
	@go test -covermode=atomic -count=1  -short $$(go list ./...| grep -v /vendor/)

integration-test:docker-up
	@go test -v $$(go list ./...| grep -v /vendor/)
	@docker-compose down

docker-run: docker-image
	@echo "run docker images"
	docker run -p ${GRPC_PORT}:${GRPC_PORT} -p ${SWAGGER_PORT}:${SWAGGER_PORT} -p ${HTTP_PORT}:${HTTP_PORT} -p ${QCLOUDAPI_PORT}:${QCLOUDAPI_PORT} -it feimumoke/${PROJECT_NAME}

docker-up: docker-image
	@echo "run docker by docker-compose integrated with other services[redis,mysql,consul...]"
	@docker-compose up -d

clean:
	@docker-compose down

local-start: local
	@echo "start server local"
	@cd ./target/$(PROJECT_NAME)/ ; sh ./scripts/start.sh

local-integration-test:local
	@echo "run local integration test "
	@cd ./target/$(PROJECT_NAME)/;nohup ./scripts/start.sh > /tmp/$(PROJECT_NAME)_test.log &
	@cd .;go test -v $$(go list ./...| grep -v /vendor/)
	@cd ./target/$(PROJECT_NAME)/; ./scripts/stop.sh
	@rm /tmp/${PROJECT_NAME}_test.log

dk-local:
	@echo "build based on tools like protoc go-bindata proto-gen-doc in docker instead of local installation"
	sh ./build/build.sh -e dev -d yes

.PHONY: build
build:
	@echo "========= local build ========="
	@sh ./build/build.sh -e dev

build-c:
	@echo "========= local build ========="
	@sh ./build/build.sh -e dev -c yes

