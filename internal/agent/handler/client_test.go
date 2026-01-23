package handler

import (
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/config"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/service/mock"
	"net"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
)

func asSorted(in []string) []string {
	cp := append([]string(nil), in...)
	sort.Strings(cp)
	return cp
}

func TestClient_GetRequests(t *testing.T) {
	mn := &mock.MockManager{
		M: map[string]*model.Stat{
			"Alloc": {Type: "gauge", Value: 10.5},
			"Poll":  {Type: "counter", Value: 7},
		},
	}

	conf := config.AgentConfig{
		Port:           "127.0.0.1:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}

	c := NewClientResty(mn, conf)

	got := asSorted(c.GetRequests())
	want := asSorted([]string{
		"http://127.0.0.1:8080/update/gauge/Alloc/10.5",
		"http://127.0.0.1:8080/update/counter/Poll/7",
	})

	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected requests:\n got:  %v\n want: %v", got, want)
	}
}

func TestClient_SendRequest_CollectsErrorsOnNon200(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Skipf("port 8080 is busy on this machine: %v", err)
	}
	defer ln.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/ok/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/update/gauge/bad/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	srv := httptest.NewUnstartedServer(mux)
	srv.Listener = ln
	srv.Start()
	defer srv.Close()

	mn := &mock.MockManager{
		M: map[string]*model.Stat{
			"ok":  {Type: "gauge", Value: 1},
			"bad": {Type: "gauge", Value: 2},
		},
	}

	conf := config.AgentConfig{
		Port:           "127.0.0.1:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}

	c := NewClientResty(mn, conf)

	errs := c.SendRequest()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
}
