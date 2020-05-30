package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/chaseisabelle/histo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	adr := flag.String("address", ":8080", "server host")
	uri := flag.String("uri", "/metrics", "the metrics scrape uri")

	flag.Parse()

	mux := http.NewServeMux()

	mux.Handle(*uri, promhttp.Handler())

	mux.HandleFunc("/histogram", func(res http.ResponseWriter, req *http.Request) {
		raw := struct {
			Name    string            `json:"name"`
			Help    string            `json:"help"`
			Labels  map[string]string `json:"labels"`
			Buckets []float64         `json:"buckets"`
			Value   float64           `json:"value"`
		}{}

		err := json.NewDecoder(req.Body).Decode(&raw)

		if err == nil && raw.Name == "" {
			err = errors.New("no name")
		}

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)

			return
		}

		go func() {
			strs := make([]string, 0, len(raw.Labels))

			for str := range raw.Labels {
				strs = append(strs, str)
			}

			his, err := histo.New(raw.Name, raw.Help, strs, raw.Buckets)

			if err != nil {
				println(err)

				return
			}

			strs = make([]string, 0, len(raw.Labels))

			for _, str := range raw.Labels {
				strs = append(strs, str)
			}

			his.Observe(raw.Value, strs...)
		}()
	})

	err := http.ListenAndServe(*adr, mux)

	if err != nil {
		panic(err)
	}
}
