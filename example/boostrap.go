package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// User lambda response struct
type User struct {
	Name string `json:"name"`
}

func main() {
	lambda.Start(func(ctx context.Context, payload User) (map[string]string, error) {
		return map[string]string{
			"name": payload.Name,
		}, nil
	})
}
