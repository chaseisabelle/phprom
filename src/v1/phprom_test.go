package v1

import (
	"strings"
	"sync"
	"testing"
)

func Test_Counter_Success(t *testing.T) {
	ns := "namespace"
	nom := "counter"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	val := map[string]string{"a": "A", "b": "B", "c": "C"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res1, err := regCounter(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to register counter: %+v", err)
	}

	if res1 == nil {
		t.Errorf("bad register counter response: %+v", res1)
	}

	res2, err := recCounter(srv, ns, nom, val, 5)

	if err != nil {
		t.Errorf("failed to record counter: %+v", err)
	}

	if res2 == nil {
		t.Errorf("bad record counter response: %+v", err)
	}

	res3, err := srv.Get(nil, &GetRequest{})

	if err != nil {
		t.Errorf("failed to get counter metrics: %+v", err)
	}

	sub := "# HELP namespace_counter who cares?\n"
	sub += "# TYPE namespace_counter counter\n"
	sub += "namespace_counter{a=\"A\",b=\"B\",c=\"C\"} 5\n"

	if !strings.Contains(res3.Metrics, sub) {
		t.Errorf("failed to detect counter metrics: %+v", res3)
	}
}

func Test_Histogram_Success(t *testing.T) {
	ns := "namespace"
	nom := "histo"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	val := map[string]string{"a": "A", "b": "B", "c": "C"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res1, err := regHisto(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to register histogram: %+v", err)
	}

	if res1 == nil {
		t.Errorf("bad register histogram response: %+v", res1)
	}

	res2, err := recHisto(srv, ns, nom, val, 2)

	if err != nil {
		t.Errorf("failed to record histogram: %+v", err)
	}

	if res2 == nil {
		t.Errorf("bad record histogram response: %+v", err)
	}

	res3, err := srv.Get(nil, &GetRequest{})

	if err != nil {
		t.Errorf("failed to get histogram metrics: %+v", err)
	}

	sub := "# HELP namespace_histo who cares?\n"
	sub += "# TYPE namespace_histo histogram\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"0.005\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"0.01\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"0.025\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"0.05\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"0.1\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"0.25\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"0.5\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"1\"} 0\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"2.5\"} 1\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"5\"} 1\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"10\"} 1\n"
	sub += "namespace_histo_bucket{a=\"A\",b=\"B\",c=\"C\",le=\"+Inf\"} 1\n"
	sub += "namespace_histo_sum{a=\"A\",b=\"B\",c=\"C\"} 2\n"
	sub += "namespace_histo_count{a=\"A\",b=\"B\",c=\"C\"} 1\n"

	if !strings.Contains(res3.Metrics, sub) {
		t.Errorf("failed to detect histogram metrics: %+v", res3)
	}
}

func Test_Summary_Success(t *testing.T) {
	ns := "namespace"
	nom := "summary"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	val := map[string]string{"a": "A", "b": "B", "c": "C"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res1, err := regSumm(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to register summary: %+v", err)
	}

	if res1 == nil {
		t.Errorf("bad register summary response: %+v", res1)
	}

	res2, err := recSumm(srv, ns, nom, val, 1)

	if err != nil {
		t.Errorf("failed to record summary: %+v", err)
	}

	if res2 == nil {
		t.Errorf("bad record summary response: %+v", err)
	}

	res3, err := srv.Get(nil, &GetRequest{})

	if err != nil {
		t.Errorf("failed to get summary metrics: %+v", err)
	}

	sub := "# HELP namespace_summary who cares?\n"
	sub += "# TYPE namespace_summary summary\n"
	sub += "namespace_summary_sum{a=\"A\",b=\"B\",c=\"C\"} 1\n"
	sub += "namespace_summary_count{a=\"A\",b=\"B\",c=\"C\"} 1\n"

	if !strings.Contains(res3.Metrics, sub) {
		t.Errorf("failed to detect summary metrics: %+v", res3)
	}
}

func Test_Gauge_Success(t *testing.T) {
	ns := "namespace"
	nom := "gauge"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	val := map[string]string{"a": "A", "b": "B", "c": "C"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res1, err := regGauge(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to register gauge: %+v", err)
	}

	if res1 == nil {
		t.Errorf("bad register gauge response: %+v", res1)
	}

	res2, err := recGauge(srv, ns, nom, val, 5)

	if err != nil {
		t.Errorf("failed to record gauge: %+v", err)
	}

	if res2 == nil {
		t.Errorf("bad record gauge response: %+v", err)
	}

	res3, err := srv.Get(nil, &GetRequest{})

	if err != nil {
		t.Errorf("failed to get gauge metrics: %+v", err)
	}

	sub := "# HELP namespace_gauge who cares?\n"
	sub += "# TYPE namespace_gauge gauge\n"
	sub += "namespace_gauge{a=\"A\",b=\"B\",c=\"C\"} 5\n"

	if !strings.Contains(res3.Metrics, sub) {
		t.Errorf("failed to detect summary metrics: %+v", res3)
	}
}

func Test_RegisterCounter_Failure(t *testing.T) {
	ns := "namespace"
	nom := "counter"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res, err := regCounter(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to do first counter register: %+v", err)
	}

	if res == nil {
		t.Errorf("bad reg counter resp: %+v", res)
	}

	res, err = regCounter(srv, ns, nom, des, []string{"pee"})

	if err == nil {
		t.Errorf("expected error")
	}
}

func Test_RegisterHistogram_Failure(t *testing.T) {
	ns := "namespace"
	nom := "histo"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res, err := regHisto(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to do first histogram register: %+v", err)
	}

	if res == nil {
		t.Errorf("bad reg histogram resp: %+v", res)
	}

	res, err = regHisto(srv, ns, nom, des, []string{"pee"})

	if err == nil {
		t.Errorf("expected error")
	}
}

func Test_RegisterSummary_Failure(t *testing.T) {
	ns := "namespace"
	nom := "summ"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res, err := regSumm(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to do first summary register: %+v", err)
	}

	if res == nil || res.Registered {
		t.Errorf("bad reg summary resp: %+v", res)
	}

	res, err = regSumm(srv, ns, nom, des, []string{"pee"})

	if err == nil {
		t.Errorf("expected error")
	}
}

func Test_RegisterGauge_Failure(t *testing.T) {
	ns := "namespace"
	nom := "gauge"
	des := "who cares?"
	lab := []string{"a", "b", "c"}
	srv, err := New()

	if err != nil {
		t.Errorf("failed to get instance: %+v", err)
	}

	res, err := regGauge(srv, ns, nom, des, lab)

	if err != nil {
		t.Errorf("failed to do first gauge register: %+v", err)
	}

	if res == nil {
		t.Errorf("bad reg gauge resp: %+v", res)
	}

	res, err = regGauge(srv, ns, nom, des, []string{"pee"})

	if err == nil {
		t.Errorf("expected error")
	}
}

func Test_Race_Success(t *testing.T) {
	max := 1000
	ns := "test"
	cn := "cnt"
	hn := "histo"
	sn := "summ"
	gn := "gauge"
	lab := []string{"a", "b", "c"}
	lvs := map[string]string{"a":"A", "b":"B", "c":"C"}
	des := "poop"
	gos := 100
	wg := sync.WaitGroup{}
	srv, err := New()

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < gos; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < max; i++ {
				_, err := regCounter(srv, ns, cn, des, lab)

				if err != nil {
					t.Error(err)
				}

				_, err = recCounter(srv, ns, cn, lvs, 1)

				if err != nil {
					t.Error(err)
				}

				_, err = regHisto(srv, ns, hn, des, lab)

				if err != nil {
					t.Error(err)
				}

				_, err = recHisto(srv, ns, hn, lvs, 1)

				if err != nil {
					t.Error(err)
				}

				_, err = regSumm(srv, ns, sn, des, lab)

				if err != nil {
					t.Error(err)
				}

				_, err = recSumm(srv, ns, sn, lvs, 1)

				if err != nil {
					t.Error(err)
				}

				_, err = regGauge(srv, ns, gn, des, lab)

				if err != nil {
					t.Error(err)
				}

				_, err = recGauge(srv, ns, gn, lvs, 1)

				if err != nil {
					t.Error(err)
				}
			}
		}()
	}

	wg.Wait()
}

// helpers

func regCounter(s *PHProm, ns string, n string, d string, l []string) (*RegisterResponse, error) {
	return s.RegisterCounter(nil, &RegisterCounterRequest{
		Namespace:   ns,
		Name:        n,
		Description: d,
		Labels:      l,
	})
}

func recCounter(s *PHProm, ns string, n string, l map[string]string, v float32) (*RecordResponse, error) {
	return s.RecordCounter(nil, &RecordCounterRequest{
		Namespace: ns,
		Name:      n,
		Labels:    l,
		Value:     v,
	})
}

func regHisto(s *PHProm, ns string, n string, d string, l []string) (*RegisterResponse, error) {
	return s.RegisterHistogram(nil, &RegisterHistogramRequest{
		Namespace:   ns,
		Name:        n,
		Description: d,
		Labels:      l,
	})
}

func recHisto(s *PHProm, ns string, n string, l map[string]string, v float32) (*RecordResponse, error) {
	return s.RecordHistogram(nil, &RecordHistogramRequest{
		Namespace: ns,
		Name:      n,
		Labels:    l,
		Value:     v,
	})
}

func regSumm(s *PHProm, ns string, n string, d string, l []string) (*RegisterResponse, error) {
	return s.RegisterSummary(nil, &RegisterSummaryRequest{
		Namespace:   ns,
		Name:        n,
		Description: d,
		Labels:      l,
	})
}

func recSumm(s *PHProm, ns string, n string, l map[string]string, v float32) (*RecordResponse, error) {
	return s.RecordSummary(nil, &RecordSummaryRequest{
		Namespace: ns,
		Name:      n,
		Labels:    l,
		Value:     v,
	})
}

func regGauge(s *PHProm, ns string, n string, d string, l []string) (*RegisterResponse, error) {
	return s.RegisterGauge(nil, &RegisterGaugeRequest{
		Namespace:   ns,
		Name:        n,
		Description: d,
		Labels:      l,
	})
}

func recGauge(s *PHProm, ns string, n string, l map[string]string, v float32) (*RecordResponse, error) {
	return s.RecordGauge(nil, &RecordGaugeRequest{
		Namespace: ns,
		Name:      n,
		Labels:    l,
		Value:     v,
	})
}
