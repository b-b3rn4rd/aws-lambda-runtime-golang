package runtime

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

// InvocationError error response struct
type InvocationError struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorType    string `json:"errorType"`
}

// Invocation struct contains invocation data
type Invocation struct {
	Headers *InvocationHeaders
	Payload []byte
}

// NewInvocation creates a new invocation struct
func NewInvocation(headers *InvocationHeaders, payload []byte) *Invocation {
	return &Invocation{
		Headers: headers,
		Payload: payload,
	}
}

// Context creates context object from response headers
func (i *Invocation) Context() context.Context {
	invokeContext := context.Background()

	lc := &lambdacontext.LambdaContext{
		AwsRequestID:       i.Headers.LambdaRuntimeAwsRequestID,
		InvokedFunctionArn: i.Headers.LambdaRuntimeInvokedFunctionArn,
	}

	invokeContext = lambdacontext.NewContext(invokeContext, lc)

	return context.WithValue(invokeContext, LambdaRuntimeTraceID("x-amzn-trace-id"), i.Headers.LambdaRuntimeTraceID)
}
