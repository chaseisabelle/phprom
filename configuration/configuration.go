package configuration

import "flag"

type Configuration struct {
	Host string
}

func FromFlags() *Configuration {
	host := flag.String("host", ":3333", "the tcp server host")

	flag.Parse()

	return &Configuration{
		Host: *host,
	}
}
