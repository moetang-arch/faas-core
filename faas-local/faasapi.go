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

	// out 0
	out0 := funcType.Out(0)
	switch out0.Kind() {
	case reflect.Ptr:
		h.resultKind = reflect.Ptr
		_, ok := reflect.New(out0.Elem()).Interface().(faas.Deserializable)
		if ok {
			h.resultDeserialization = true
		} else {
			h.resultDeserialization = false
		}
	case reflect.Struct:
		h.resultKind = reflect.Struct
		_, ok := reflect.New(out0).Elem().Interface().(faas.Deserializable)
		if ok {
			h.resultDeserialization = true
		} else {
			h.resultDeserialization = false
		}
	default:
		panic(errors.New("unsupported result type"))
	}

	return h
}

func (this *HandlerStruct) JsonCall(jsonStr string, ctx context.Context) (result string, err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = errors.New("panic occurs. information: " + fmt.Sprint(p))
		}
	}()

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
	var values []reflect.Value
	if this.paramKind == reflect.Ptr {
		// ptr
		values = this.function.Call([]reflect.Value{reflect.ValueOf(ctx), paramValue})
	} else {
		// struct
		values = this.function.Call([]reflect.Value{reflect.ValueOf(ctx), paramValue.Elem()})
	}

	errValue := values[1]
	if !errValue.IsNil() {
		return "", errValue.Interface().(error)
	}

	valueValue := values[0]
	if this.resultKind == reflect.Ptr && valueValue.IsNil() {
		return "", nil
	}
	if this.resultDeserialization {
		return valueValue.Interface().(faas.Deserializable).ToJson()
	} else {
		data, err := json.Marshal(valueValue.Interface())
		if err != nil {
			return "", err
		} else {
			return string(data), nil
		}
	}
}

func (this *HandlerStruct) BinaryCall(data []byte, ctx context.Context) ([]byte, error) {
	//TODO
	return nil, nil
}

func (this *HandlerStruct) SupportBinarySerdes() bool {
	return this.paramSerialization && this.resultDeserialization
}

func InitHandler() map[string]*HandlerStruct {
	serviceMap := faas.GetServiceMap()

	result := make(map[string]*HandlerStruct)
	for k, v := range serviceMap {
		result[k] = NewHandlerStruct(v)
	}

	return result
}
