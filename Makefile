
.PHONY: setup
setup:
	brew list golangci-lint &>/dev/null || brew install golangci-lint
	brew list pubsub_cli &>/dev/null || (brew tap k-yomo/pubsub_cli && brew install pubsub_cli)
	GO111MODULE=off go get -u github.com/cosmtrek/air
	GO111MODULE=off go get -u github.com/mattn/goreman
	go mod download
	cd src/web_client; yarn


.PHONY: run
run:
	docker-compose up -d db pubsub datastore redis
	./scripts/create_local_pubsub_resources.sh
	open http://local.eitan-flash.com:3000
	goreman -set-ports=false start

.PHONY: run-dc
run-dc:
	docker-compose up -d db pubsub
	./scripts/create_local_pubsub_resources.sh
	docker-compose up
	open http://local.eitan-flash.com:3000

.PHONY: gen-model
gen-model:
	rm src/account_service/infra/*.xo.go src/eitan_service/infra/*.xo.go
	xo mysql://root@localhost:13306/accountdb --int32-type int64 --uint32-type int64  --template-path xo_templates -o src/account_service/internal/infra
	xo mysql://root@localhost:13306/eitandb --int32-type int64 --uint32-type int64  --template-path xo_templates -o src/eitan_service/internal/infra

test-account:
	go test ./src/account_service/... -v $(TESTARGS) -coverprofile=account_service.coverage.out

test-eitan:
	go test ./src/eitan_service/... -v $(TESTARGS) -coverprofile=eitan_service.coverage.out

test-notification:
	go test ./src/notification_service/... -v $(TESTARGS) -coverprofile=notification_service.coverage.out

test-pubsub-publisher:
	go test ./src/pubsub_publisher/... -v $(TESTARGS) -coverprofile=pubsub_publisher.coverage.out

lint:
	@golangci-lint run

.PHONY: gen-graphql
gen-graphql:
	cd src/eitan_service; go run github.com/99designs/gqlgen
	cd src/web_client; yarn codegen

.PHONY: gen_proto
gen-proto:
	rm -f src/internal/pb/eitan/*
	protoc -I defs/proto defs/proto/*.proto \
	--experimental_allow_proto3_optional \
	--go_out=plugins=grpc:src/internal/pb/eitan

.PHONY: reset-db
reset-db:
	docker-compose stop db
	docker-compose rm -f db
	docker volume rm eitan_db_data
	docker-compose up -d db

.PHONY: tf-symlink
tf-symlink:
	#cd ./terraform/dev && ln -sf ../shared/* .
	cd ./terraform/prod && ln -sf ../shared/* .
