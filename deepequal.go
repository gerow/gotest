// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Modified from reflect/deepequal.go

// Deep equality test via reflection

package gotest

import (
    "reflect"
    "testing"
)

// During deepValueEqual, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited are stored in a map indexed by 17 * a1 + a2;
type visit struct {
	a1   uintptr
	a2   uintptr
	typ  reflect.Type
	next *visit
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueEqual(v1, v2 reflect.Value, visited map[uintptr]*visit, depth int, t *testing.T) (b bool) {
        v1i := v1.Interface()
        v2i := v2.Interface()
	if !v1.IsValid() || !v2.IsValid() {
                if v1.IsValid() != v2.IsValid() {
                  t.Errorf("%v and %v are unequal because one is invalid and the other is not", v1, v2)
                }
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
                t.Errorf("%v and %v are unequal because they are of different types", v1, v2)
		return false
	}

	// if depth > 10 { panic("deepValueEqual") }	// for debugging

	if v1.CanAddr() && v2.CanAddr() {
		addr1 := v1.UnsafeAddr()
		addr2 := v2.UnsafeAddr()
		if addr1 > addr2 {
			// Canonicalize order to reduce number of entries in visited.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are identical ...
		if addr1 == addr2 {
			return true
		}

		// ... or already seen
		h := 17*addr1 + addr2
		seen := visited[h]
		typ := v1.Type()
		for p := seen; p != nil; p = p.next {
			if p.a1 == addr1 && p.a2 == addr2 && p.typ == typ {
				return true
			}
		}

		// Remember for later.
		visited[h] = &visit{addr1, addr2, typ, seen}
	}

	switch v1.Kind() {
	case reflect.Array:
		if v1.Len() != v2.Len() {
                        t.Errorf("%v and %v are Arrays of different length (%d and %d)", v1i, v2i, v1.Len(), v2.Len())
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1, t) {
                                t.Errorf("%v and %v are arrays with differing elements", v1i, v2i)
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
                        t.Errorf("%v and %v Slices but one is nil", v1i, v2i)
			return false
		}
		if v1.Len() != v2.Len() {
                        t.Errorf("%v and %v are Slices of different length (%d and %d)", v1i, v2i, v1.Len(), v2.Len())
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1, t) {
                                t.Errorf("%v and %v are slices with differing elements", v1i, v2i)
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
                        if v1.IsNil() != v2.IsNil() {
                          t.Errorf("%v and %v are Interfaces but one is nil", v1i, v2i)
                        }
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1, t)
	case reflect.Ptr:
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1, t)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepValueEqual(v1.Field(i), v2.Field(i), visited, depth+1, t) {
                                t.Errorf("%v and %v are structs with differing elements", v1i, v2i)
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
                        t.Errorf("%v and %v are Maps but one is nil", v1i, v2i)
			return false
		}
		if v1.Len() != v2.Len() {
                        t.Errorf("%v and %v are Maps of different length (%d and %d)", v1i, v2i, v1.Len(), v2.Len())
			return false
		}
		for _, k := range v1.MapKeys() {
			if !deepValueEqual(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1, t) {
                                t.Errorf("%v and %v are maps with differing elements", v1i, v2i)
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
                t.Errorf("%v and %v are functions but they aren't both nil", v1i, v2i)
		return false
	default:
		// Normal equality suffices
		return v1.Interface() == v2.Interface()
	}

	panic("Not reached")
}

// DeepEqual tests for deep equality. It uses normal == equality where possible
// but will scan members of arrays, slices, maps, and fields of structs. It correctly
// handles recursive types. Functions are equal only if they are both nil.
//
// This has been modified from reflect/deepequal.go to be more useful for testing by
// firing off a t.Errorf when it discovers something that makes the two objects unequal
func AssertDeepEqual(a1, a2 interface{}, t *testing.T) bool {
	if a1 == nil || a2 == nil {
                if a1 != a2 {
                  t.Errorf("%v and %v are unequal because one is nil and the other isn't", a1, a2)
                }
		return a1 == a2
	}
	v1 := reflect.ValueOf(a1)
	v2 := reflect.ValueOf(a2)
	if v1.Type() != v2.Type() {
                t.Errorf("%v and %v are unequal because they are different types", a1, a2)
		return false
	}
	return deepValueEqual(v1, v2, make(map[uintptr]*visit), 0, t)
}
