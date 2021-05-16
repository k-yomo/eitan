
.PHONY: setup
setup:
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
	xo mysql://root@localhost:13306/accountdb --int32-type int64 --uint32-type int64  --template-path xo_templates -o src/account_service/infra
	xo mysql://root@localhost:13306/eitandb --int32-type int64 --uint32-type int64  --template-path xo_templates -o src/eitan_service/infra

.PHONY: gen-graphql
gen-graphql:
	cd src/eitan_service; go generate ./...
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
