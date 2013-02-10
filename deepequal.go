package gotest

import (
	"reflect"
	"testing"
)

func AssertDeepEqual(a, b interface{}, t *testing.T) {
	if a == nil || b == nil {
		if a != b {
			t.Errorf("Found that %v and %v are unequal because one of them is nil", a, b)
		}
		return
	}
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("reflect.DeepEqual found %v and %v to be inequal", a, b)
		if v1.Type() != v2.Type() {
			t.Errorf("%v and %v are different types", a, b)
			return
		}
		switch v1.Kind() {
		case reflect.Array:
			if v1.Len() != v2.Len() {
				t.Errorf("%v and %v are arrays of different length", a, b)
			}
			for i := 0; i < v1.Len(); i++ {
				AssertDeepEqual(v1.Index(i), v2.Index(1), t)
			}
			return
		case reflect.Slice:
			if v1.IsNil() != v2.IsNil() {
				t.Errorf("%v and %v are both slices but one is nil", v1, v2)
				return
			}
			if v1.Len() != v2.Len() {
				t.Errorf("%v and %v are both slices but are of different length (%d and %d respectively)", v1, v2, v1.Len(), v2.Len())
				return
			}
			for i := 0; i < v1.Len(); i++ {
				AssertDeepEqual(v1, v2, t)
			}
			return
		case reflect.Interface:
			if v1.IsNil() || v2.IsNil() {
				t.Errorf("%v and %v are interfaces but are unequal because one is nil", v1, v2)
				return
			}
			AssertDeepEqual(v1.Elem(), v2.Elem(), t)
			return
		case reflect.Ptr:
			AssertDeepEqual(v1.Elem(), v2.Elem(), t)
			return
		case reflect.Struct:
			for i, n := 0, v1.NumField(); i < n; i++ {
				AssertDeepEqual(v1.Elem(), v2.Elem(), t)
			}
			return
		case reflect.Map:
			if v1.IsNil() != v2.IsNil() {
				t.Errorf("%v and %v are maps but are unequal because one is nil", v1, v2)
				return
			}
			if v1.Len() != v2.Len() {
				t.Errorf("%v and %v are both maps but are of different length (%d and %d respectively)", v1, v2, v2.Len(), v2.Len())
				return
			}
			for _, v := range v1.MapKeys() {
				AssertDeepEqual(v1.MapIndex(v), v2.MapIndex(v), t)
			}
			return
		case reflect.Func:
			if !v1.IsNil() || !v2.IsNil() {
				t.Errorf("Equality check failed because %v and %v are function pointers and one or more of them are zero", v1, v2)
				return
			}
			return
		}
	}
}
