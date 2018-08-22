package main

import (
	"context"
	"encoding/json"
	"errors"
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

// implements Ser/Des
//func (Response) ToJson() (string, error) {
//	panic("implement me")
//}
//
//func (Response) ToBytes() ([]byte, error) {
//	panic("implement me")
//}

func HandleRequest(ctx context.Context, request *Request) (response *Response, err error) {
	if len(request.Name) == 0 {
		return nil, errors.New("name is empty")
	}
	resp := new(Response)
	resp.Result = "Hello, " + request.Name
	return resp, nil
}

func TestInitHandler(t *testing.T) {
	m := InitHandler()
	result, err := m["demo"].JsonCall(`{"name":"1"}`, context.Background())
	t.Log(result, len(result))
	t.Log(err)
}

func BenchmarkInitHandler(b *testing.B) {
	m := InitHandler()
	handler := m["demo"]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.JsonCall(`{"name":"1"}`, context.Background())
	}
}

func BenchmarkInitHandler2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var request = new(Request)
		err := json.Unmarshal([]byte(`{"name":"1"}`), request)
		if err != nil {
			panic(err)
		}
		HandleRequest(context.Background(), request)
	}
}
