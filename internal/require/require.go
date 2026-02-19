// Package require provides test assertion helpers that mirror
// the stretchr/testify/require API surface used in this project.
// Failures call t.FailNow (fatal).
package require

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Equal checks that expected == actual using reflect.DeepEqual. Fatal on failure.
func Equal(t testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		return
	}
	t.Fatalf("%sExpected: %#v\n     Got: %#v", formatMsg(msgAndArgs), expected, actual)
}

// NotNil checks that obj is not nil. Fatal on failure.
func NotNil(t testing.TB, obj interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !isNil(obj) {
		return
	}
	t.Fatalf("%sExpected non-nil value", formatMsg(msgAndArgs))
}

// Nil checks that obj is nil. Fatal on failure.
func Nil(t testing.TB, obj interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if isNil(obj) {
		return
	}
	t.Fatalf("%sExpected nil, got: %#v", formatMsg(msgAndArgs), obj)
}

// NoError checks that err is nil. Fatal on failure.
func NoError(t testing.TB, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		return
	}
	t.Fatalf("%sUnexpected error: %v", formatMsg(msgAndArgs), err)
}

// Error checks that err is not nil. Fatal on failure.
func Error(t testing.TB, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		return
	}
	t.Fatalf("%sExpected an error but got nil", formatMsg(msgAndArgs))
}

// --- helpers ---

func isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func containsElement(list interface{}, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	switch listValue.Kind() {
	case reflect.String:
		elemStr, ok := element.(string)
		if !ok {
			return false, false
		}
		return true, strings.Contains(listValue.String(), elemStr)
	case reflect.Slice, reflect.Array:
		for i := 0; i < listValue.Len(); i++ {
			if reflect.DeepEqual(listValue.Index(i).Interface(), element) {
				return true, true
			}
		}
		return true, false
	case reflect.Map:
		keys := listValue.MapKeys()
		for _, k := range keys {
			if reflect.DeepEqual(k.Interface(), element) {
				return true, true
			}
		}
		return true, false
	}
	return false, false
}

func formatMsg(msgAndArgs []interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	if len(msgAndArgs) == 1 {
		if msg, ok := msgAndArgs[0].(string); ok {
			return msg + ": "
		}
		return fmt.Sprintf("%v: ", msgAndArgs[0])
	}
	if msg, ok := msgAndArgs[0].(string); ok {
		return fmt.Sprintf(msg, msgAndArgs[1:]...) + ": "
	}
	return ""
}
