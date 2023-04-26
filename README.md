AWS Lambda APIGateway OpenTelemetry instrumentation
------

Wrap your handler with `otelawslambda.WrapAPIGatewayLambda` to instrument and record OpenTelemetry HTTP spans.

```go
func HandleRequest(ctx context.Context, name events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		Body:    "hello world",
		Headers: map[string]string{"Content-Type": "text/plain"},
	}, nil
}

func main() {
	lambda.Start(otelawslambda.WrapAPIGatewayLambda(HandleRequest))
}
```