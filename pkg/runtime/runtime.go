package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"
)

// Runtime runtime mandatory functions
type Runtime interface {
	Response(requestID string, payload []byte) error
	Error(requestID string, err error) error
	InitError(err error) error
	Next() (*Invocation, error)
}

// GoRuntime runtime implementation in golang
type GoRuntime struct {
	endpoint    string
	version     string
	handlerFunc interface{}
	client      *http.Client
}

// Start creates and runs golang runtime
func Start(handlerFunc interface{}) {
	r := &GoRuntime{
		endpoint:    os.Getenv("AWS_LAMBDA_RUNTIME_API"),
		version:     "2018-06-01",
		handlerFunc: handlerFunc,
		client:      &http.Client{},
	}

	r.Run()
}

func (r *GoRuntime) invocationError(err error) *InvocationError {
	var errorName string
	if errorType := reflect.TypeOf(err); errorType.Kind() == reflect.Ptr {
		errorName = errorType.Elem().Name()
	} else {
		errorName = errorType.Name()
	}

	return &InvocationError{
		ErrorMessage: err.Error(),
		ErrorType:    errorName,
	}
}

// Response sends lambdas non-error response
func (r *GoRuntime) Response(requestID string, payload []byte) {
	url := fmt.Sprintf("http://%s/%s/runtime/invocation/%s/response", r.endpoint, r.version, requestID)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		panic(err)
	}
}

// Error sends lambdas error response
func (r *GoRuntime) Error(requestID string, err error) {
	url := fmt.Sprintf("http://%s/%s/runtime/invocation/%s/error", r.endpoint, r.version, requestID)

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(r.invocationError(err))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		panic(err)
	}
}

// InitError sends init error response
func (r *GoRuntime) InitError(err error) {
	url := fmt.Sprintf("http://%s/%s/runtime/init/error", r.endpoint, r.version)

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(r.invocationError(err))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		panic(err)
	}
}

// Next retrieves information about next invocation
func (r *GoRuntime) Next() (*Invocation, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://%s/%s/runtime/invocation/next", r.endpoint, r.version),
		nil,
	)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := NewInvocationHeaders(resp.Header)

	return NewInvocation(headers, body), nil
}

// Run runtime process
func (r *GoRuntime) Run() {
	handler := lambda.NewHandler(r.handlerFunc)

	for {
		invocation, err := r.Next()
		if err != nil {
			r.InitError(err)
			return
		}

		ctx := invocation.Context()
		os.Setenv("_X_AMZN_TRACE_ID", invocation.Headers.LambdaRuntimeTraceID)

		payload, err := handler.Invoke(ctx, invocation.Payload)
		if err != nil {
			r.Error(invocation.Headers.LambdaRuntimeAwsRequestID, err)
			return
		}

		r.Response(invocation.Headers.LambdaRuntimeAwsRequestID, payload)
	}
}
