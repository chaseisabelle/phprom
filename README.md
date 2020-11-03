# phprom
a prometheus metric datastore for php apps

---
### prerequisites
- install the [client](https://github.com/chaseisabelle/phprom-client)
    - `composer require chaseisabelle/phprom-client`

---
### usage
- from command line: `go run cmd/v1/main.go --address=0.0.0.0:3333 --api=grpc`
- get the latest image from [docker](https://hub.docker.com/repository/docker/chaseisabelle/phprom)
    - `docker run phprom ./phprom --address=0.0.0.0:3333`

##### cli options
```bash
$ go run cmd/v1/main.go --help
Usage of phprom:
  -address string
    	the host:port to listen on (default "0.0.0.0:3333")
  -api string
    	the api to use (grpc, rest, resp) (default "grpc")
```

---
### apis
- [grpc](https://grpc.io/)
- rest