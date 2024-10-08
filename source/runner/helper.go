package runner

import (
	"fmt"
	"math"
	"meow/source/runner/object"
)

func nativeBoolToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	}
	return true
}

func isAllTruthy(obj []object.Object) bool {
	for _, o := range obj {
		if !isTruthy(o) {
			return false
		}
	}
	return true
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func IsError(obj object.Object) bool {
	if _, ok := obj.(*object.Error); ok {
		return true
	}
	return false
}

func checkForDefault(name string) bool {
	switch name {
	case "meow":
		return true
	case "len":
		return true
	case "tail":
		return true
	}
	return false
}

func isWhole(x float64) bool {
	return math.Ceil(x) == x
}
