package main

import (
	"context"

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

func MyHandler(ctx context.Context, e Event) error {
	logger.Info("in my handler", zap.Any("event", e))
	return nil
}

func main() {
	lambda.Start(MyHandler)
}
