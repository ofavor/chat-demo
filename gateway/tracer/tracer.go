package tracer

import (
	"fmt"
	"gateway/log"

	"github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	// reporter = httpreporter.NewReporter("http://localhost:9411/api/v2/spans")

	tracer *zipkin.Tracer
)

func Init(addr string) {
	var err error
	reporter := httpreporter.NewReporter(fmt.Sprintf("http://%s/api/v2/spans", addr))
	tracer, err = zipkin.NewTracer(reporter)
	if err != nil {
		log.Error("Create tracer error:", err)
	}
}

func StartSpan(name string) zipkin.Span {
	return tracer.StartSpan(name)
}
