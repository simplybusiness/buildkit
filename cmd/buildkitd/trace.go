package main

import (
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var tracer opentracing.Tracer
var closeTracer io.Closer

type DDCloser struct{
	closerFunction func()
}

func (d DDCloser) Close() error {
	d.closerFunction()
	return nil
}

func init() {

	tracer = opentracing.NoopTracer{}

	if traceAddr := os.Getenv("JAEGER_TRACE"); traceAddr != "" {
		tr, err := jaeger.NewUDPTransport(traceAddr, 0)
		if err != nil {
			panic(err)
		}

		tracer, closeTracer = jaeger.NewTracer(
			"buildkitd",
			jaeger.NewConstSampler(true),
			jaeger.NewRemoteReporter(tr),
		)
	}

	if traceAddr := os.Getenv( "DD_AGENT_TRACE"); traceAddr != "" {
		tracer = opentracer.New(ddtracer.WithService("buildkitd"))

		// Stop it using the regular Stop call for the tracer package.
		closeTracer = DDCloser{ddtracer.Stop}

		// Set the global OpenTracing tracer.
		//opentracing.SetGlobalTracer(t)
	}

}
