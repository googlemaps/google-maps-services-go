package metrics_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robin-samuel/maps"
	"github.com/robin-samuel/maps/metrics"
)

type testReporter struct {
	start, end int
}

func (t *testReporter) NewRequest(name string) metrics.Request {
	t.start++
	return &testMetric{reporter: t}
}

type testMetric struct {
	reporter *testReporter
}

func (t *testMetric) EndRequest(ctx context.Context, err error, httpResp *http.Response, metro string) {
	t.reporter.end++
}

func mockServer(codes []int, body string) *httptest.Server {
	i := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(codes[i])
		i++
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, body)
	}))
	return server
}

func TestClientWithMetricReporter(t *testing.T) {
	server := mockServer([]int{200}, `{"results" : [], "status" : "OK"}`)
	defer server.Close()
	reporter := &testReporter{}
	c, err := maps.NewClient(
		maps.WithAPIKey("AIza-Maps-API-Key"),
		maps.WithBaseURL(server.URL),
		maps.WithMetricReporter(reporter))
	if err != nil {
		t.Errorf("Unable to create client with MetricReporter")
	}
	r := &maps.ElevationRequest{
		Locations: []maps.LatLng{
			{
				Lat: 39.73915360,
				Lng: -104.9847034,
			},
		},
	}
	_, err = c.Elevation(context.Background(), r)
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}
	if reporter.start != 1 {
		t.Errorf("expected one start call")
	}
	if reporter.end != 1 {
		t.Errorf("expected one end call")
	}
}
