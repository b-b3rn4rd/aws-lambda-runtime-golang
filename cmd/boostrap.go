package main

import (
	"github.com/b-b3rn4rd/aws-lambda-runtime-golang/pkg/runtime"
	"os"
)

func main()  {
	r := runtime.NewRuntime(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	r.Run()
}
