package v1

import "fmt"

type Server interface {
	Serve() error
}

func New(api API, adr string) (Server, error) {
	switch api {
	case GrpcApi:
		return newGRPC(adr)
	case RestApi:
		return newREST(adr)
	default:
		break
	}

	return nil, fmt.Errorf("invalid api: %s", api)
}
