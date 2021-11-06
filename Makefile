PROTO_DIR := ${PWD}/api/proto/v1
PB_OUT := ${PWD}/pkg/v1
NETWORK := phprom
IMAGE := chaseisabelle/phprom:latest
CONTAINER := phprom
REQUESTS := 2000
CONCURRENTS := 20

genpb:
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

netup:
	docker network create ${NETWORK}

netrm:
	docker network rm ${NETWORK}

build:
	docker build --no-cache --tag ${IMAGE} .

run:
	docker run -it -d --rm --net ${NETWORK} --name ${CONTAINER} ${IMAGE}

up:
	make netup
	make build
	make run

logs:
	docker logs -f -n 100 ${CONTAINER}

kill:
	docker kill ${CONTAINER}

rmi:
	docker rmi ${IMAGE}

nuke:
	make kill || true
	make rmi || true
	make netrm || true

reup:
	make nuke
	make up

ghzget:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.Get \
		-d '{}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzregcounter:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RegisterCounter \
		-d '{"namespace":"test","name":"counter","description":"test counter","labels":["foo"]}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzreghisto:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RegisterHistogram \
		-d '{"namespace":"test","name":"histo","description":"test histo","labels":["foo"],"buckets":[1,2,3,4,5]}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzregsumm:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RegisterSummary \
		-d '{"namespace":"test","name":"sum","description":"test sum","labels":["foo"],"objectives":[{"key":1.1,"value":1.1}],"maxAge":20,"ageBuckets":100,"bufCap":100}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzreggauge:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RegisterGauge \
		-d '{"namespace":"test","name":"gauge","description":"test gauge","labels":["foo"]}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzreccounter:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RecordCounter \
		-d '{"namespace":"test","name":"counter","value":1,"labels":{"foo":"bar"}}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzrechisto:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RecordHistogram \
		-d '{"namespace":"test","name":"histo","value":5,"labels":{"foo":"bar"}}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzrecsumm:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RecordSummary \
		-d '{"namespace":"test","name":"sum","value":2,"labels":{"foo":"bar"}}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

ghzrecgauge:
	docker run --rm --net ${NETWORK} -v $(shell pwd)/api/proto:/proto obvionaoe/ghz \
		--insecure \
		--proto /proto/v1/service.proto \
		--call PHProm.v1.Service.RecordGauge \
		-d '{"namespace":"test","name":"gauge","value":1.1,"labels":{"foo":"bar"}}' \
		-n ${REQUESTS} \
		-c ${CONCURRENTS} \
		${CONTAINER}:3333

