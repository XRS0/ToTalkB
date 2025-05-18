NOTIFY_PATH=./notify
EVENT_PATH=./event_manager
GEN_PROTO_PATH=./proto
GO_OUT_PATH_NOTIFY=${NOTIFY_PATH}/internal/domain
GO_GRPC_OUT_PATH_NOTIFY=${NOTIFY_PATH}/internal/domain
GO_OUT_PATH_EVENT=${EVENT_PATH}/internal/domain
GO_GRPC_OUT_PATH_EVENT=${EVENT_PATH}/internal/domain

setpath:
	export "PATH=$PATH:$(go env GOPATH)/bin"

gnfe:
	protoc --go_out=${GO_OUT_PATH_NOTIFY} \
	--go-grpc_out=${GO_GRPC_OUT_PATH_NOTIFY} \
	${GEN_PROTO_PATH}/event.proto

gnfn:
	protoc --go_out=${GO_OUT_PATH_EVENT} \
	--go-grpc_out=${GO_GRPC_OUT_PATH_EVENT} \
	${GEN_PROTO_PATH}/notification.proto

generate-notify: gnfe
	protoc --go_out=${GO_OUT_PATH_NOTIFY} \
	--go-grpc_out=${GO_GRPC_OUT_PATH_NOTIFY} \
	${GEN_PROTO_PATH}/notification.proto

generate-event: gnfn
	protoc --go_out=${GO_OUT_PATH_EVENT} \
	--go-grpc_out=${GO_GRPC_OUT_PATH_EVENT} \
	${GEN_PROTO_PATH}/event.proto

generate-all-proto:
	make generate-event
	make generate-notify

build-notify:
	cd notify && go build -o notify main.go

build-event:
	cd event_manager && go build -o event_manager main.go

build-all:
	make build-notify
	make build-event

run-notify:
	cd notify && ./notify

run-event:
	cd event_manager && ./event_manager

run-test-client:
	cd test_client && ./test_client

create-tables:
	docker cp ./schema/000001_init.up.sql totalk_db:/tmp/init.sql
	docker exec -i totalk_db psql -U totalkadmin -d totalk_db -f /tmp/init.sql

run-all:
	make run-notify
	make run-event
