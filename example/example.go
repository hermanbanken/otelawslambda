package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hermanbanken/otelawslambda"
)

func HandleRequest(ctx context.Context, name events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		Body:    "hello world",
		Headers: map[string]string{"Content-Type": "text/plain"},
	}, nil
}

func main() {
	lambda.Start(otelawslambda.WrapAPIGatewayLambda(HandleRequest))
}
