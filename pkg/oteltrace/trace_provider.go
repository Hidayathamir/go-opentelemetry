package oteltrace

import (
	"context"

	"github.com/erajayatech/go-opentelemetry/pkg/config"
	"github.com/erajayatech/go-opentelemetry/pkg/errutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// source: https://opentelemetry.io/docs/languages/go/instrumentation/#traces
// source: https://opentelemetry.io/docs/languages/go/exporters/#otlp-traces-over-grpc

// NewTraceProvider return opentelemetry trace provider.
func NewTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	fail := func(err error) (*sdktrace.TracerProvider, error) {
		return nil, errutil.AddFuncName(err, errutil.WithSkip(1))
	}

	err := validateConfig()
	if err != nil {
		return fail(err)
	}

	sampler := sdktrace.NeverSample()
	isAlwaysSample, err := config.GetOtelAlwaysSample()
	if err != nil {
		return fail(err)
	}

	if isAlwaysSample {
		sampler = sdktrace.AlwaysSample()
	}

	opt, err := getNROption()
	if err != nil {
		return fail(err)
	}

	exporter, err := otlptracegrpc.New(ctx, opt...)
	if err != nil {
		return fail(err)
	}

	serviceName, err := config.GetServiceName()
	if err != nil {
		return fail(err)
	}

	appVersion, err := config.GetAppVersion()
	if err != nil {
		return fail(err)
	}

	appEnv, err := config.GetAppEnvironment()
	if err != nil {
		return fail(err)
	}

	_resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(appVersion),
		attribute.String("environment", appEnv),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(_resource),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

func validateConfig() error {
	errlist := errutil.NewErrList()

	var err error
	_, err = config.GetOtelAlwaysSample()
	errlist.AddIfErr(err)
	_, err = config.GetServiceName()
	errlist.AddIfErr(err)
	_, err = config.GetAppVersion()
	errlist.AddIfErr(err)
	_, err = config.GetAppEnvironment()
	errlist.AddIfErr(err)
	_, err = config.GetOtelOTLPNewrelicHost()
	errlist.AddIfErr(err)
	_, err = config.GetOtelOTLPNewrelicHeaderAPIKey()
	errlist.AddIfErr(err)

	if errlist.IsErr() {
		err := errlist.Err()
		return errutil.AddFuncName(err)
	}

	return nil
}

func getNROption() ([]otlptracegrpc.Option, error) {
	otelNRHost, err := config.GetOtelOTLPNewrelicHost()
	if err != nil {
		return nil, errutil.AddFuncName(err)
	}
	otelNRHeaderAPIKey, err := config.GetOtelOTLPNewrelicHeaderAPIKey()
	if err != nil {
		return nil, errutil.AddFuncName(err)
	}
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(otelNRHost),
		otlptracegrpc.WithHeaders(map[string]string{"api-key": otelNRHeaderAPIKey}),
		otlptracegrpc.WithCompressor("gzip"),
	}
	return opts, nil
}

func NewInMemoryTracer() trace.Tracer {
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(tracetest.NewInMemoryExporter()),
		sdktrace.WithResource(resource.Default()),
	).Tracer("")
}
