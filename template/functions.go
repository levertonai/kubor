package template

import (
	"fmt"
	nt "text/template"
)

type FunctionProvider interface {
	GetFunctions() (Functions, error)
}

type Function interface {
	Execute(context ExecutionContext, args ...interface{}) (interface{}, error)
}

type Functions map[string]Function

func (instance Functions) CreateFuncMap(context ExecutionContext) (nt.FuncMap, error) {
	result := nt.FuncMap{}
	for name, function := range instance {
		result[name] = instance.methodCaller(context, function)
	}
	return result, nil
}

func (instance Functions) CreateDummyFuncMap() (nt.FuncMap, error) {
	result := nt.FuncMap{}
	for name := range instance {
		result[name] = instance.dummyMethod
	}
	return result, nil
}

func (instance Functions) methodCaller(context ExecutionContext, function Function) func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		return function.Execute(context, args...)
	}
}

func (instance Functions) dummyMethod(args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("method not initialzed")
}

type ExecutionContext interface {
	GetTemplate() Template
	GetFactory() Factory
	GetData() interface{}
}

type ExecutionContextImpl struct {
	Template Template
	Factory  Factory
	Data     interface{}
}

func (instance *ExecutionContextImpl) GetTemplate() Template {
	return instance.Template
}

func (instance *ExecutionContextImpl) GetFactory() Factory {
	return instance.Factory
}

func (instance *ExecutionContextImpl) GetData() interface{} {
	return instance.Data
}
