package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	l, _ := zap.NewProduction()
	logger = l
	defer logger.Sync() // flushes buffer, if any
}

type Event struct {
	Name string `json:"name"`
}

func MyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	//	logger.Info("in my handler", zap.Any("event", e))
	var res *events.APIGatewayProxyResponse

	logger.Info("recieved event", zap.Any("method", event.HTTPMethod), zap.Any("path", event.Path), zap.Any("body", event.Body))

	if event.Path == "/hello" {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "hello",
		}
	} else {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "default",
		}
	}

	return res, nil
}

func main() {
	lambda.Start(MyHandler)
}
