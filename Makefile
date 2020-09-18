PROTO_DIR := ${PWD}/api/proto/v1
PHP_OUT := ${PWD}/api/php/v1
GO_OUT := ${PWD}/api/go/v1
CP_SERVICE = ${PWD}/src/v1

generate-client:
	rm -Rf ${PHP_OUT} && \
	mkdir ${PHP_OUT} && \
	docker pull jaegertracing/protobuf && \
	docker run --rm  \
		-v${PROTO_DIR}:${PROTO_DIR} \
		-v${PHP_OUT}:${PHP_OUT} \
		-w${PHP_OUT} \
		jaegertracing/protobuf:latest \
		--proto_path=${PROTO_DIR} \
        --php_out=${PHP_OUT} \
        -I/usr/include/github.com/gogo/protobuf \
        ${PROTO_DIR}/service.proto

generate-server:
	rm -Rf ${GO_OUT} && \
	mkdir ${GO_OUT} && \
	docker pull jaegertracing/protobuf && \
	docker run --rm  \
		-v${PROTO_DIR}:${PROTO_DIR} \
		-v${GO_OUT}:${GO_OUT} \
		-w${GO_OUT} \
		jaegertracing/protobuf:latest \
		--proto_path=${PROTO_DIR} \
        --go_out=plugins=grpc:${GO_OUT} \
        -I/usr/include/github.com/gogo/protobuf \
        ${PROTO_DIR}/service.proto

update-server: generate-server
	cp ${GO_OUT}/service.pb.go pkg/v1/

update-client: generate-client
