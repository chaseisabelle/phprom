package v1

import (
	"context"
	"encoding/json"
	"fmt"
	phprom_v1 "github.com/chaseisabelle/phprom/pkg/v1"
	v1 "github.com/chaseisabelle/phprom/src/v1"
	"github.com/prometheus/common/log"
	"net/http"
)

type RESTServer struct {
	address string
	phprom  *v1.PHProm
}

func newREST(adr string) (*RESTServer, error) {
	php, err := v1.New()

	if err != nil {
		return nil, err
	}

	srv := &RESTServer{
		address: adr,
		phprom:  php,
	}

	http.HandleFunc("/metrics", srv.get)
	http.HandleFunc("/register/counter", srv.registerCounter)
	http.HandleFunc("/register/histogram", srv.registerHistogram)
	http.HandleFunc("/register/summary", srv.registerSummary)
	http.HandleFunc("/register/gauge", srv.registerGauge)
	http.HandleFunc("/record/counter", srv.recordCounter)
	http.HandleFunc("/record/histogram", srv.recordHistogram)
	http.HandleFunc("/record/summary", srv.recordSummary)
	http.HandleFunc("/record/gauge", srv.recordGauge)

	return srv, nil
}

func (r *RESTServer) Serve() error {
	return http.ListenAndServe(r.address, nil)
}

func (r *RESTServer) get(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodGet) {
		return
	}

	gr, err := r.phprom.Get(context.Background(), &phprom_v1.GetRequest{})

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	r.respond(res, []byte(gr.GetMetrics()))
}

func (r *RESTServer) registerCounter(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RegisterCounterRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RegisterCounter(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) registerHistogram(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RegisterHistogramRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RegisterHistogram(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) registerSummary(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RegisterSummaryRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RegisterSummary(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) registerGauge(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RegisterGaugeRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RegisterGauge(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) recordCounter(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RecordCounterRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RecordCounter(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) recordHistogram(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RecordHistogramRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RecordHistogram(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) recordSummary(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RecordSummaryRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RecordSummary(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) recordGauge(res http.ResponseWriter, req *http.Request) {
	if !r.allowed(req, res, http.MethodPost) {
		return
	}

	rrq := &phprom_v1.RecordGaugeRequest{}
	err := json.NewDecoder(req.Body).Decode(rrq)

	if err != nil {
		r.bad(res, err)

		return
	}

	rrr, err := r.phprom.RecordGauge(context.Background(), rrq)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.marshal(res, rrr)
}

func (r *RESTServer) allowed(req *http.Request, res http.ResponseWriter, mth string) bool {
	ok := req.Method == mth

	if !ok {
		http.Error(res, fmt.Sprintf("method not allowed: %s", req.Method), http.StatusMethodNotAllowed)
	}

	return ok
}

func (r *RESTServer) marshal(res http.ResponseWriter, raw interface{}) {
	enc, err := json.Marshal(raw)

	if err != nil {
		r.failure(res, err)

		return
	}

	r.respond(res, enc)
}

func (r *RESTServer) respond(res http.ResponseWriter, bod []byte) {
	_, err := res.Write(bod)

	if err != nil {
		log.Error(err)
	}
}

func (r *RESTServer) bad(res http.ResponseWriter, err error) {
	http.Error(res, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadRequest)
}

func (r *RESTServer) failure(res http.ResponseWriter, err error) {
	http.Error(res, err.Error(), http.StatusInternalServerError)
}
