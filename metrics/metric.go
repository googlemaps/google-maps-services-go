package metrics

import (
	"context"
	"net/http"
)

type Reporter interface {
	NewRequest(name string) Request
}

type Request interface {
	EndRequest(ctx context.Context, err error, httpResp *http.Response, metro string)
}

type NoOpReporter struct {
}

func (n NoOpReporter) NewRequest(name string) Request {
	return noOpRequest{}
}

type noOpRequest struct {
}

func (n noOpRequest) EndRequest(ctx context.Context, err error, httpResp *http.Response, metro string) {
}
