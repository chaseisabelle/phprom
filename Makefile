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
	docker build --no-cache --tag chaseisabelle/phprom:latest .

run:
	docker run -it -d --rm --name chaseisabelle-phprom-latest chaseisabelle/phprom:latest

up:
	make build && make run

logs:
	docker logs -f -n 100 chaseisabelle-phprom-latest

kill:
	docker kill chaseisabelle-phprom-latest

rmi:
	docker rmi chaseisabelle/phprom:latest

nuke:
	make kill || true
	make rmi || true

reup:
	make nuke
	make up

# todo put these in an isolated container
ghz-get:
	ghz --insecure --proto api/proto/v1/service.proto --call PHProm.v1.Service.Get -d '{}' -n 2000 -c 20 127.0.0.1:3333