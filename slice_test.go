package luar

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func Test_slice(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	things := []string{
		"cake",
		"wallet",
		"calendar",
		"phone",
		"speaker",
	}

	L.SetGlobal("things", New(L, things))

	testReturn(t, L, `return #things, things[1], things[5]`, "5", "cake", "speaker")
	if err := L.DoString(`things[1] = "cookie"`); err != nil {
		t.Fatal(err)
	}
	if things[0] != "cookie" {
		t.Fatalf(`expected things[0] = "cookie", got %s`, things[0])
	}
}

func Test_slice_2(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	items := make([]string, 0, 10)

	L.SetGlobal("items", New(L, items))

	testReturn(t, L, `return #items`, "0")
	testReturn(t, L, `return items:capacity()`, "10")
	testReturn(t, L, `items = items:append("hello", "world"); return #items`, "2")
	testReturn(t, L, `return items:capacity()`, "10")
	testReturn(t, L, `return items[1]`, "hello")
	testReturn(t, L, `return items[2]`, "world")
}

func Test_slice_iterator(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	s := []string{
		"hello",
		"there",
	}

	e := []string{}

	L.SetGlobal("s", New(L, s))
	L.SetGlobal("e", New(L, e))

	testReturn(t, L, `local itr = s(); local a, b = itr(); local c, d = itr(); return a, b, c, d`, "1", "hello", "2", "there")
	testReturn(t, L, `local itr = e(); local a, b = itr(); return a, b`, "nil", "nil")
}

func Test_slice_immutable_edit(t *testing.T) {
	// Attempting to modify an immutable slice - should error
	L := lua.NewState()
	defer L.Close()

	s := []string{"first", "second"}
	L.SetGlobal("s", New(L, s, ReflectOptions{Immutable: true}))

	testError(t, L, `s[1] = "hi"`, "invalid operation on immutable slice")
}

func Test_slice_immutable_append(t *testing.T) {
	// Attempting to append to an immutable slice - should error
	L := lua.NewState()

	s := []string{"first", "second"}
	L.SetGlobal("s", New(L, s, ReflectOptions{Immutable: true}))

	testError(t, L, `s = s:append("hi")`, "invalid operation on immutable slice")
}

func Test_slice_immutable_access(t *testing.T) {
	// Attempting to access a member of an immutable slice - should be fine
	L := lua.NewState()
	defer L.Close()

	s := []string{"first", "second"}
	L.SetGlobal("s", New(L, s, ReflectOptions{Immutable: true}))

	testReturn(t, L, `return s[1]`, "first")
}

func Test_slice_transparentptr_call(t *testing.T) {
	// Iterate over a slice via the call method. Access an undefined pointer
	// field on the returned object. Returned object should inherit the transparent
	// pointer behavior, and the field should be accessible without indirection.
	L := lua.NewState()
	defer L.Close()

	val := "foo"
	b := TestTransparentPtrAccessB{}
	b.Str = &val

	slice := []*TestTransparentPtrAccessB{&b}

	L.SetGlobal("slice", New(L, slice, ReflectOptions{TransparentPointers: true}))

	testReturn(
		t,
		L,
		`
		local result = ""
		for i, b in slice() do
			result = result .. i .. " " .. b.Str
		end
		return result
		`,
		"1 foo",
	)
}
