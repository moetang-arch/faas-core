package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/moetang-arch/faas-api"
)

func init() {
	faas.Register("demo", HandleRequest)
}

type Request struct {
	Name string `json:"name"`
}

// implements Ser/Des
//func (this *Request) FromJson(json string) error {
//	panic("implement me")
//}
//
//func (this *Request) FromBytes(data []byte) error {
//	panic("implement me")
//}

type Response struct {
	Result string
}

func HandleRequest(ctx context.Context, request *Request) (response *Response, err error) {
	fmt.Println(request)
	if len(request.Name) == 0 {
		return nil, errors.New("name is empty")
	}
	resp := new(Response)
	resp.Result = "Hello, " + request.Name
	return resp, nil
}

func TestInitHandler(t *testing.T) {
	m := InitHandler()
	m["demo"].JsonCall(`{"name":"1"}`, context.Background())
}
