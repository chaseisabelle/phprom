PROTO_DIR := ${PWD}/api/proto/v1
PB_OUT := ${PWD}/pkg/v1

gen-pb:
	rm -Rf ${PB_OUT} && \
	mkdir ${PB_OUT} && \
	docker pull jaegertracing/protobuf && \
	docker run --rm  \
		-v${PROTO_DIR}:${PROTO_DIR} \
		-v${PB_OUT}:${PB_OUT} \
		-w${PB_OUT} \
		jaegertracing/protobuf:latest \
		--proto_path=${PROTO_DIR} \
        --go_out=plugins=grpc:${PB_OUT} \
        -I/usr/include/github.com/gogo/protobuf \
        ${PROTO_DIR}/service.proto

build:
	docker build --no-cache --tag phprom:latest .

run:
	docker run --rm --name phprom phprom:latest

up:
	make build && make run

rmi:
	docker rmi phprom:latest