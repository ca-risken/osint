package main

import (
	"context"
)

func main() {
	//	conf, err := newBackendConfig()
	//	if err != nil {
	//		panic(err)
	//	}

	ctx := context.Background()
	consumer := newSQSConsumer()
	appLogger.Info("Start the intrigue SQS consumer server...")
	consumer.Start(ctx, newHandler())
}
