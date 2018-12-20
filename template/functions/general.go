package functions

import (
	"fmt"
	"reflect"
	"strings"
)

var FuncOptional = Function{
	Description: "Asks the given <holder> if a property of given <name> exists and returns it. Otherwise the result is <nil>.",
	Parameters: Parameters{{
		Name: "defaultValue",
	}, {
		Name: "given",
	}},
}.MustWithFunc(func(defaultValue interface{}, given interface{}) interface{} {
	if empty(given) {
		return defaultValue
	}
	return given
})

var FuncEmpty = Function{
	Description: "Checks the given <argument> if it is empty or not.",
	Parameters: Parameters{{
		Name: "argument",
	}},
}.MustWithFunc(empty)

var FuncIsNotEmpty = Function{
	Description: "Checks the given <argument> if it is not empty.",
	Parameters: Parameters{{
		Name: "argument",
	}},
}.MustWithFunc(func(given interface{}) bool {
	return !empty(given)
})

var FuncDefault = Function{
	Description: "Asks the given <holder> if a property of given <name> exists and returns it. Otherwise the result is <nil>.",
	Parameters: Parameters{{
		Name: "name",
	}, {
		Name: "holder",
	}},
}.MustWithFunc(func(name string, holder interface{}) (interface{}, error) {
	v := reflect.ValueOf(holder)
	t := v.Type()
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsNil() {
		return nil, nil
	}
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if field, ok := t.FieldByName(name); ok {
			fieldValue := v.FieldByIndex(field.Index)
			return fieldValue.Interface(), nil
		}
		return nil, nil
	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(name))
		if !value.IsValid() {
			return nil, nil
		}
		return value.Interface(), nil
	}
	return nil, fmt.Errorf("cannot get value for '%s' because cannot handle values of type %v", name, t)
})

var FuncContains = Function{
	Description: `Checks if the given <input> string contains the given <search> string.`,
	Parameters: Parameters{{
		Name: "search",
	}, {
		Name: "input",
	}},
}.MustWithFunc(func(search interface{}, input interface{}) (bool, error) {
	if s, ok := input.(string); ok {
		return strings.Contains(s, fmt.Sprint(search)), nil
	} else if s, ok := input.(*string); ok {
		return strings.Contains(*s, fmt.Sprint(search)), nil
	}
	return false, fmt.Errorf("currently contains only supports strings but got: %v", reflect.TypeOf(input))
})

var FuncsGeneral = Functions{
	"optional":   FuncOptional,
	"empty":      FuncEmpty,
	"isNotEmpty": FuncIsNotEmpty,
	"default":    FuncDefault,
	"contains":   FuncContains,
}
var CategoryGeneral = Category{
	Functions: FuncsGeneral,
}

func empty(given interface{}) bool {
	g := reflect.ValueOf(given)
	if !g.IsValid() {
		return true
	}

	// Basically adapted from text/template.isTrue
	switch g.Kind() {
	default:
		return g.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return g.Len() == 0
	case reflect.Bool:
		return g.Bool() == false
	case reflect.Complex64, reflect.Complex128:
		return g.Complex() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return g.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return g.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return g.Float() == 0
	case reflect.Struct:
		return false
	}
}