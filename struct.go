package luar

import (
	"reflect"

	"github.com/yuin/gopher-lua"
)

func structIndex(L *lua.LState) int {
	ref, opts, mt, isPtr := check(L, 1, reflect.Struct)
	key := L.CheckString(2)

	if !isPtr {
		if fn := mt.method(key); fn != nil {
			L.Push(fn)
			return 1
		}
	}

	if fn := mt.ptrMethod(key); fn != nil {
		if opts.Immutable {
			L.RaiseError("cannot call pointer methods on immutable objects")
		}

		L.Push(fn)
		return 1
	}

	ref = reflect.Indirect(ref)
	index := mt.fieldIndex(key)
	if index == nil {
		return 0
	}
	field := ref.FieldByIndex(index)
	if !field.CanInterface() {
		L.RaiseError("cannot interface field " + key)
	}

	switch field.Kind() {
	case reflect.Ptr:
		if opts.TransparentPointers {
			// Initialize pointers on first access
			if !field.IsValid() || field.IsNil() {
				if !field.CanSet() {
					L.RaiseError("cannot transparently create pointer field " + key)
				}
				if !opts.Immutable {
					field.Set(reflect.New(field.Type().Elem()))
				}
			}

			// Return the value of the pointer
			if field.Elem().IsValid() {
				field = field.Elem()
			}
		}
	case reflect.Slice:
		if opts.TransparentPointers {
			// Initialize slices on first access
			if !field.IsValid() || field.IsNil() {
				if !field.CanSet() {
					L.RaiseError("cannot transparently create slice " + key)
				}
				if !opts.Immutable {
					field.Set(reflect.MakeSlice(field.Type(), 0, 10))
				}
			}
		}
	}

	if (field.Kind() == reflect.Struct || field.Kind() == reflect.Array) && field.CanAddr() {
		field = field.Addr()
	}

	L.Push(New(L, field.Interface(), opts))
	return 1
}

func structNewIndex(L *lua.LState) int {
	ref, opts, mt, isPtr := check(L, 1, reflect.Struct)

	if opts.Immutable {
		L.RaiseError("invalid operation on immutable struct")
	}

	if isPtr {
		ref = ref.Elem()
	}

	key := L.CheckString(2)
	value := L.CheckAny(3)

	index := mt.fieldIndex(key)
	if index == nil {
		L.RaiseError("unknown field " + key)
	}
	field := ref.FieldByIndex(index)

	if opts.TransparentPointers {
		// With transparent pointers, we are going to get passed the new value
		// by... value, not by reference. Thus, if the current field is a
		// pointer type, we need some extra work to reflect a new object for
		// assignment to the field.
		if field.Type().Kind() == reflect.Ptr {
			hint := field.Type().Elem()
			goValue := lValueToReflect(L, value, hint, nil)

			if !field.CanSet() {
				// Can happen if the field is on a struct reflected by
				// value instead of by reference.
				L.RaiseError("cannot set field " + key)
			}
			field.Set(reflect.New(goValue.Type()))
			field.Elem().Set(goValue)
			return 0
		}
	}

	if !field.CanSet() {
		L.RaiseError("cannot set field " + key)
	}
	field.Set(lValueToReflect(L, value, field.Type(), nil))
	return 0
}
