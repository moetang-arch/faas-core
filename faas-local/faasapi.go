package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/moetang-arch/faas-api"
)

type HandlerStruct struct {
	function reflect.Value

	paramKind          reflect.Kind
	paramSerialization bool
	paramType          reflect.Type

	resultKind            reflect.Kind
	resultDeserialization bool
	resultType            reflect.Type
}

func NewHandlerStruct(f interface{}) *HandlerStruct {
	h := new(HandlerStruct)
	h.function = reflect.ValueOf(f)

	funcType := reflect.TypeOf(f)

	// in 1
	in1 := funcType.In(1)
	h.paramKind = in1.Kind()
	switch h.paramKind {
	case reflect.Ptr:
		h.paramType = in1.Elem()
		_, ok := reflect.New(h.paramType).Interface().(faas.Serializable)
		if ok {
			h.paramSerialization = true
		} else {
			h.paramSerialization = false
		}
	case reflect.Struct:
		h.paramType = in1
		_, ok := reflect.New(h.paramType).Interface().(faas.Serializable)
		if ok {
			h.paramSerialization = true
		} else {
			h.paramSerialization = false
		}
	default:
		panic(errors.New("unsupported param type"))
	}

	//TODO out 0

	return h
}

func (this *HandlerStruct) JsonCall(jsonStr string, ctx context.Context) (string, error) {
	var paramValue reflect.Value
	if this.paramSerialization {
		v := reflect.New(this.paramType)
		err := v.Interface().(faas.Serializable).FromJson(jsonStr)
		if err != nil {
			return "", err
		}
		paramValue = v
	} else {
		v := reflect.New(this.paramType)
		err := json.Unmarshal([]byte(jsonStr), v.Interface())
		if err != nil {
			return "", err
		}
		paramValue = v
	}
	if this.paramKind == reflect.Ptr {
		_ = this.function.Call([]reflect.Value{reflect.ValueOf(ctx), paramValue})
	} else {
		_ = this.function.Call([]reflect.Value{reflect.ValueOf(ctx), paramValue.Elem()})
	}
	//TODO
	return "", nil
}

func (this *HandlerStruct) BinaryCall(data []byte, ctx context.Context) ([]byte, error) {
	//TODO
	return nil, nil
}

func (this *HandlerStruct) SupportBinarySerdes() bool {
	return this.paramSerialization && this.resultDeserialization
}

func InitHandler() map[string]*HandlerStruct {
	globalNameSpace := faas.GetGlobalServiceNameSpace()
	fmt.Println(globalNameSpace)
	serviceMap := faas.GetServiceMap()

	result := make(map[string]*HandlerStruct)
	for k, v := range serviceMap {
		result[k] = NewHandlerStruct(v)
	}

	return result
}
