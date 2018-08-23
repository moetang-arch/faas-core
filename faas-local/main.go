package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/moetang-arch/faas-api"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	handlers := InitHandler()
	serviceNamespace := faas.GetGlobalServiceNameSpace()

	mux := http.NewServeMux()

	var uriPrefix = "/"
	if len(serviceNamespace) > 0 {
		uriPrefix = uriPrefix + serviceNamespace + "/"
	}

	var listenAddr = "0.0.0.0:8001"
	fmt.Println("registered services:")
	// register http handlers
	for k, v := range handlers {
		h := HttpHandler{}
		h.handlerStruct = v
		uri := uriPrefix + k
		mux.Handle(uri, h)
		// print services
		fmt.Println(fmt.Sprintf("http://%s%s", listenAddr, uri))
	}
	fmt.Println() // new line

	handler := mux

	fmt.Println("Listening HTTP requests on: 0.0.0.0:8001")

	err := http.ListenAndServe(listenAddr, handler)
	if err != nil {
		panic(err)
	}
}

type HttpHandler struct {
	handlerStruct *HandlerStruct
}

func (this HttpHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		logger.Println(fmt.Sprintf("%dms %s %s", time.Now().Sub(start).Nanoseconds()/1000000, req.Method, req.RequestURI))
		p := recover()
		if p != nil {
			dealWithPanic(resp, req, p)
			return
		}
	}()
	switch req.Method {
	case "GET":
		fallthrough
	case "POST":
		err := req.ParseForm()
		if err != nil {
			dealWithError(resp, req, err)
			return
		}
		param := req.Form.Get("param")
		if len(param) == 0 {
			dealWithError(resp, req, errors.New("param is empty"))
			return
		}
		result, err := this.handlerStruct.JsonCall(param, context.Background())
		if err != nil {
			dealWithError(resp, req, err)
			return
		}
		if len(result) == 0 {
			dealWithError(resp, req, errors.New("empty result"))
			return
		}
		resp.Header().Add("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(result))
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}

type ErrResult struct {
	Msg string `json:"msg"`
}

func dealWithError(resp http.ResponseWriter, req *http.Request, err error) {
	resp.Header().Add("Content-Type", "application/json")
	resp.WriteHeader(http.StatusBadRequest)
	e := ErrResult{
		Msg: err.Error(),
	}
	data, _ := json.Marshal(e)
	resp.Write(data)
}

func dealWithPanic(resp http.ResponseWriter, req *http.Request, p interface{}) {
	resp.Header().Add("Content-Type", "application/json")
	resp.WriteHeader(http.StatusBadRequest)
	e := ErrResult{
		Msg: fmt.Sprint(p),
	}
	data, _ := json.Marshal(e)
	resp.Write(data)
}
