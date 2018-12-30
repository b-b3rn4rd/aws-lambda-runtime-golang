package main

import (
	"context"

	"github.com/b-b3rn4rd/aws-lambda-runtime-golang/pkg/runtime"
)

// User lambda response struct
type User struct {
	Name string `json:"name"`
}

func main() {
	runtime.Start(func(ctx context.Context, payload User) (map[string]string, error) {
		return map[string]string{
			"name": payload.Name,
		}, nil
	})
}
