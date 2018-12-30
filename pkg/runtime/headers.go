package runtime

// LambdaRuntimeTraceID custom type to avoid context collision
type LambdaRuntimeTraceID string

// InvocationHeaders struct represents invocation response headers
type InvocationHeaders struct {
	LambdaRuntimeAwsRequestID       string `json:"Lambda-Runtime-Aws-Request-Id"`
	LambdaRuntimeDeadlineMs         string `json:"Lambda-Runtime-Deadline-Ms"`
	LambdaRuntimeInvokedFunctionArn string `json:"Lambda-Runtime-Invoked-Function-Arn"`
	LambdaRuntimeTraceID            string `json:"Lambda-Runtime-Trace-Id"`
}

// NewInvocationHeaders creates a new instance of invocation struct from response headers
func NewInvocationHeaders(headers map[string][]string) *InvocationHeaders {
	return &InvocationHeaders{
		LambdaRuntimeAwsRequestID:       headers["Lambda-Runtime-Aws-Request-Id"][0],
		LambdaRuntimeDeadlineMs:         headers["Lambda-Runtime-Deadline-Ms"][0],
		LambdaRuntimeInvokedFunctionArn: headers["Lambda-Runtime-Invoked-Function-Arn"][0],
		LambdaRuntimeTraceID:            headers["Lambda-Runtime-Trace-Id"][0],
	}
}
