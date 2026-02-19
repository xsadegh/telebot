// Package assert provides test assertion helpers that mirror
// the stretchr/testify/assert API surface used in this project.
// Failures call t.Errorf (non-fatal).
package assert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Equal checks that expected == actual using reflect.DeepEqual.
func Equal(t testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	t.Errorf("%sExpected: %#v\n     Got: %#v", formatMsg(msgAndArgs), expected, actual)
	return false
}

// NotEqual checks that expected != actual.
func NotEqual(t testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		return true
	}
	t.Errorf("%sValues should not be equal: %#v", formatMsg(msgAndArgs), actual)
	return false
}

// Nil checks that obj is nil.
func Nil(t testing.TB, obj interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if isNil(obj) {
		return true
	}
	t.Errorf("%sExpected nil, got: %#v", formatMsg(msgAndArgs), obj)
	return false
}

// NotNil checks that obj is not nil.
func NotNil(t testing.TB, obj interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !isNil(obj) {
		return true
	}
	t.Errorf("%sExpected non-nil value", formatMsg(msgAndArgs))
	return false
}

// True checks that value is true.
func True(t testing.TB, value bool, msgAndArgs ...interface{}) bool {
	t.Helper()
	if value {
		return true
	}
	t.Errorf("%sExpected true", formatMsg(msgAndArgs))
	return false
}

// False checks that value is false.
func False(t testing.TB, value bool, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !value {
		return true
	}
	t.Errorf("%sExpected false", formatMsg(msgAndArgs))
	return false
}

// NoError checks that err is nil.
func NoError(t testing.TB, err error, msgAndArgs ...interface{}) bool {
	t.Helper()
	if err == nil {
		return true
	}
	t.Errorf("%sUnexpected error: %v", formatMsg(msgAndArgs), err)
	return false
}

// Error checks that err is not nil.
func Error(t testing.TB, err error, msgAndArgs ...interface{}) bool {
	t.Helper()
	if err != nil {
		return true
	}
	t.Errorf("%sExpected an error but got nil", formatMsg(msgAndArgs))
	return false
}

// EqualError checks that err.Error() == expectedErrMsg.
func EqualError(t testing.TB, err error, expectedErrMsg string, msgAndArgs ...interface{}) bool {
	t.Helper()
	if err == nil {
		t.Errorf("%sExpected error %q but got nil", formatMsg(msgAndArgs), expectedErrMsg)
		return false
	}
	if err.Error() != expectedErrMsg {
		t.Errorf("%sExpected error: %q\n     Got: %q", formatMsg(msgAndArgs), expectedErrMsg, err.Error())
		return false
	}
	return true
}

// Contains checks that s contains substr (for strings) or element (for slices/maps).
func Contains(t testing.TB, s, contains interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	ok, found := containsElement(s, contains)
	if !ok {
		t.Errorf("%sCannot check containment on type %T", formatMsg(msgAndArgs), s)
		return false
	}
	if found {
		return true
	}
	t.Errorf("%s%#v does not contain %#v", formatMsg(msgAndArgs), s, contains)
	return false
}

// NotContains checks that s does not contain contains.
func NotContains(t testing.TB, s, contains interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	ok, found := containsElement(s, contains)
	if !ok {
		t.Errorf("%sCannot check containment on type %T", formatMsg(msgAndArgs), s)
		return false
	}
	if !found {
		return true
	}
	t.Errorf("%s%#v should not contain %#v", formatMsg(msgAndArgs), s, contains)
	return false
}

// Len checks that the object has the expected length.
func Len(t testing.TB, obj interface{}, length int, msgAndArgs ...interface{}) bool {
	t.Helper()
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if v.Len() == length {
			return true
		}
		t.Errorf("%sExpected length %d, got %d", formatMsg(msgAndArgs), length, v.Len())
		return false
	default:
		t.Errorf("%sCannot get length of type %T", formatMsg(msgAndArgs), obj)
		return false
	}
}

// NotEmpty checks that the object is not empty (len > 0 for collections, non-zero for others).
func NotEmpty(t testing.TB, obj interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if v.Len() > 0 {
			return true
		}
	default:
		if !v.IsZero() {
			return true
		}
	}
	t.Errorf("%sExpected non-empty value", formatMsg(msgAndArgs))
	return false
}

// NotZero checks that the value is not the zero value for its type.
func NotZero(t testing.TB, obj interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if obj != nil && !reflect.DeepEqual(obj, reflect.Zero(reflect.TypeOf(obj)).Interface()) {
		return true
	}
	t.Errorf("%sExpected non-zero value", formatMsg(msgAndArgs))
	return false
}

// Panics checks that fn panics.
func Panics(t testing.TB, fn func(), msgAndArgs ...interface{}) bool {
	t.Helper()
	panicked := didPanic(fn)
	if panicked {
		return true
	}
	t.Errorf("%sExpected function to panic", formatMsg(msgAndArgs))
	return false
}

// NotPanics checks that fn does not panic.
func NotPanics(t testing.TB, fn func(), msgAndArgs ...interface{}) bool {
	t.Helper()
	panicked := didPanic(fn)
	if !panicked {
		return true
	}
	t.Errorf("%sExpected function not to panic", formatMsg(msgAndArgs))
	return false
}

// Implements checks that obj implements the specified interface.
func Implements(t testing.TB, iface interface{}, obj interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	ifaceType := reflect.TypeOf(iface).Elem()
	objType := reflect.TypeOf(obj)
	if objType.Implements(ifaceType) {
		return true
	}
	t.Errorf("%s%T does not implement %v", formatMsg(msgAndArgs), obj, ifaceType)
	return false
}

// ElementsMatch checks that two slices contain the same elements regardless of order.
func ElementsMatch(t testing.TB, listA, listB interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	a := reflect.ValueOf(listA)
	b := reflect.ValueOf(listB)
	if a.Kind() != reflect.Slice || b.Kind() != reflect.Slice {
		t.Errorf("%sBoth arguments must be slices", formatMsg(msgAndArgs))
		return false
	}
	if a.Len() != b.Len() {
		t.Errorf("%sLengths differ: %d vs %d\n  Left: %#v\n Right: %#v", formatMsg(msgAndArgs), a.Len(), b.Len(), listA, listB)
		return false
	}

	// Build a "used" tracker for b
	used := make([]bool, b.Len())
	for i := 0; i < a.Len(); i++ {
		found := false
		for j := 0; j < b.Len(); j++ {
			if used[j] {
				continue
			}
			if reflect.DeepEqual(a.Index(i).Interface(), b.Index(j).Interface()) {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%sElement %#v not found in other slice\n  Left: %#v\n Right: %#v",
				formatMsg(msgAndArgs), a.Index(i).Interface(), listA, listB)
			return false
		}
	}
	return true
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

func didPanic(fn func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	fn()
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
