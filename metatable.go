package luar

import (
	"reflect"

	"github.com/yuin/gopher-lua"
)

// Metatable holds the Lua metatable for a Go type.
type Metatable struct {
	*lua.LTable
}

// MT returns the metatable for value's type. nil is returned if value's type
// does not use a custom metatable.
func MT(L *lua.LState, value interface{}) *Metatable {
	val := reflect.ValueOf(value)
	switch val.Type().Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct:
		return &Metatable{
			LTable: getMetatableFromValue(L, val),
		}
	default:
		return nil
	}
}

func (m *Metatable) method(name string) lua.LValue {
	methods := m.RawGetString("methods").(*lua.LTable)
	if fn := methods.RawGetString(name); fn != lua.LNil {
		return fn
	}
	return nil
}

func (m *Metatable) ptrMethod(name string) lua.LValue {
	methods := m.RawGetString("ptr_methods").(*lua.LTable)
	if fn := methods.RawGetString(name); fn != lua.LNil {
		return fn
	}
	return nil
}

func (m *Metatable) fieldIndex(name string) []int {
	fields := m.RawGetString("fields").(*lua.LTable)
	if index := fields.RawGetString(name); index != lua.LNil {
		return index.(*lua.LUserData).Value.([]int)
	}
	return nil
}
