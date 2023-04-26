package otelawslambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func WrapAPIGatewayLambda(fn func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)) func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (resp events.APIGatewayV2HTTPResponse, err error) {
		tracer := otel.GetTracerProvider().Tracer("")
		ctx, span := tracer.Start(
			otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(req.Headers)),
			req.RequestContext.HTTP.Method,
			trace.WithSpanKind(trace.SpanKindServer),
		)

		span.SetAttributes(
			semconv.HTTPClientIP(req.RequestContext.HTTP.SourceIP),
			semconv.HTTPMethod(req.RequestContext.HTTP.Method),
			semconv.HTTPRequestContentLength(contentLength(req.Body, req.IsBase64Encoded)),
			semconv.HTTPRoute(req.RouteKey),
			semconv.HTTPScheme(req.RequestContext.HTTP.Protocol),
			semconv.HTTPUserAgent(req.RequestContext.HTTP.UserAgent),
			semconv.HTTPTarget(req.RawPath),
			semconv.NetHostName(req.RequestContext.DomainName),
		)

		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err, trace.WithStackTrace(true))
			} else {
				span.SetStatus(codes.Ok, "")
			}
			span.SetAttributes(
				semconv.HTTPStatusCode(resp.StatusCode),
				semconv.HTTPResponseContentLength(contentLength(resp.Body, resp.IsBase64Encoded)),
			)
		}()

		defer span.End()
		return fn(ctx, req)
	}
}
