package faas_demo

import (
	"context"
	"errors"

	"github.com/moetang-arch/faas-api"
)

func init() {
	faas.Register("demo", HandleRequest)
}

type Request struct {
	Name string
}

type Response struct {
	Result string
}

func HandleRequest(ctx context.Context, request *Request) (response *Response, err error) {
	if len(request.Name) == 0 {
		return nil, errors.New("name is empty")
	}
	resp := new(Response)
	resp.Result = "Hello, " + request.Name
	return resp, nil
}
