package metrics

import (
	"context"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"net/http"
	"strconv"
	"time"
)

var (
	latency_measure = stats.Int64("maps.googleapis.com/measure/client/latency", "Latency in msecs", stats.UnitMilliseconds)

	requestName = tag.MustNewKey("request_name")
	apiStatus   = tag.MustNewKey("api_status")
	httpCode    = tag.MustNewKey("http_code")
	metroArea   = tag.MustNewKey("metro_area")

	Count = &view.View{
		Name:        "maps.googleapis.com/client/count",
		Description: "Request Counts",
		TagKeys:     []tag.Key{requestName, apiStatus, httpCode, metroArea},
		Measure:     latency_measure,
		Aggregation: view.Count(),
	}

	Latency = &view.View{
		Name:        "maps.googleapis.com/client/request_latency",
		Description: "Total time between library method called and results returned",
		TagKeys:     []tag.Key{requestName, apiStatus, httpCode, metroArea},
		Measure:     latency_measure,
		Aggregation: view.Distribution(20.0, 25.2, 31.7, 40.0, 50.4, 63.5, 80.0, 100.8, 127.0, 160.0, 201.6, 254.0, 320.0, 403.2, 508.0, 640.0, 806.3, 1015.9, 1280.0, 1612.7, 2031.9, 2560.0, 3225.4, 4063.7),
	}
)

func RegisterViews() error {
	return view.Register(Latency, Count)
}

type OpenCensusReporter struct {
}

func (o OpenCensusReporter) NewRequest(name string) Request {
	return &openCensusRequest{
		name:  name,
		start: time.Now().UnixNano() / int64(time.Millisecond),
	}
}

type openCensusRequest struct {
	name  string
	start int64
}

func (o *openCensusRequest) EndRequest(ctx context.Context, err error, httpResp *http.Response, metro string) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	duration := now - o.start
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	httpCodeStr := ""
	if httpResp != nil {
		httpCodeStr = strconv.Itoa(httpResp.StatusCode)
	}
	stats.RecordWithTags(ctx, []tag.Mutator{
		tag.Upsert(requestName, o.name),
		tag.Upsert(apiStatus, errStr),
		tag.Upsert(httpCode, httpCodeStr),
		tag.Upsert(metroArea, metro),
	}, latency_measure.M(duration))
}
