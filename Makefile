
.PHONY: run
run:
	docker-compose up

.PHONY: gen-model
gen-model:
	rm -f src/account_service/infra/*.xo.go
	xo mysql://root@localhost:13306/accountdb --int32-type int64 --uint32-type int64  --template-path xo_templates -o src/account_service/infra

.PHONY: gen-graphql
gen-graphql:
	cd src/eitan_service; go generate ./...

.PHONY: gen_proto
gen-proto:
	rm -f src/internal/pb/eitan/*
	protoc -I defs/proto defs/proto/*.proto \
	--experimental_allow_proto3_optional \
	--go_out=plugins=grpc:src/internal/pb/eitan

.PHONY: reset-db
reset-db:
	docker-compose down --volume