package v1

import "fmt"

type Server interface {
	Serve() error
}

func New(api API, adr string) (Server, error) {
	switch api {
	case GrpcApi:
		return newGRPCServer(adr)
	case RestApi:
		return newRESTServer(adr)
	default:
		break
	}

	return nil, fmt.Errorf("invalid api: %s", api)
}
