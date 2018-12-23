package runtime

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
)
type Handler func()
type Runtime struct {
	handler lambda.Handler
	endpoint string
}

func NewRuntime(endpoint string)*Runtime  {
	return &Runtime{
		endpoint: endpoint,
	}
}

func (r *Runtime) next()  error {
	resp, err := http.Get(fmt.Sprintf("%s/2018-06-01/runtime/invocation/next", r.endpoint))
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
func (r *Runtime) Run()  {
	for {
		r.next()
	}
}