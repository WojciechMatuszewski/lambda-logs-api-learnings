package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) error {
	fmt.Println("it works")

	fmt.Println("it works 2")

	return nil
}
