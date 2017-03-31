package luar

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func Test_struct(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tim := &StructTestPerson{
		Name: "Tim",
		Age:  30,
	}

	john := StructTestPerson{
		Name: "John",
		Age:  40,
	}

	L.SetGlobal("user1", New(L, tim))
	L.SetGlobal("user2", New(L, john))

	testReturn(t, L, `return user1.Name`, "Tim")
	testReturn(t, L, `return user1.Age`, "30")
	testReturn(t, L, `return user1:Hello()`, "Hello, Tim")

	testReturn(t, L, `return user2.Name`, "John")
	testReturn(t, L, `return user2.Age`, "40")
	testReturn(t, L, `local hello = user2.Hello; return hello(user2)`, "Hello, John")
}

func Test_struct_tostring(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	p1 := StructTestPerson{
		Name: "Tim",
		Age:  99,
	}
	p2 := StructTestPerson{
		Name: "John",
		Age:  2,
	}

	L.SetGlobal("p1", New(L, &p1))
	L.SetGlobal("p2", New(L, &p2))

	testReturn(t, L, `return tostring(p1)`, `Tim (99)`)
	testReturn(t, L, `return tostring(p2)`, `John (2)`)
}

func Test_struct_pointers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	p1 := StructTestPerson{
		Name: "Tim",
	}
	p2 := StructTestPerson{
		Name: "John",
	}

	L.SetGlobal("p1", New(L, &p1))
	L.SetGlobal("p1_alias", New(L, &p1))
	L.SetGlobal("p2", New(L, &p2))

	testReturn(t, L, `return -p1 == -p1`, "true")
	testReturn(t, L, `return -p1 == -p1_alias`, "true")
	testReturn(t, L, `return p1 == p1`, "true")
	testReturn(t, L, `return p1 == p1_alias`, "true")
	testReturn(t, L, `return p1 == p2`, "false")
}

func Test_struct_lstate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	p := StructTestPerson{
		Name: "Tim",
	}

	L.SetGlobal("p", New(L, &p))

	testReturn(t, L, `return p:AddNumbers(1, 2, 3, 4, 5)`, "Tim counts: 15")
}

type StructTestHidden struct {
	Name   string `luar:"name"`
	Name2  string `luar:"Name"`
	Str    string
	Hidden bool `luar:"-"`
}

func Test_struct_hiddenfields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	a := &StructTestHidden{
		Name:   "tim",
		Name2:  "bob",
		Str:    "asd123",
		Hidden: true,
	}

	L.SetGlobal("a", New(L, a))

	testReturn(t, L, `return a.name`, "tim")
	testReturn(t, L, `return a.Name`, "bob")
	testReturn(t, L, `return a.str`, "asd123")
	testReturn(t, L, `return a.Str`, "asd123")
	testReturn(t, L, `return a.Hidden`, "nil")
	testReturn(t, L, `return a.hidden`, "nil")
}

func Test_struct_method(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	p := StructTestPerson{
		Name: "Tim",
		Age:  66,
	}

	L.SetGlobal("p", New(L, &p))

	testReturn(t, L, `return p:hello()`, "Hello, Tim")
	testReturn(t, L, `return p.age`, "66")
}

type NestedPointer struct {
	B NestedPointerChild
}

type NestedPointerChild struct {
}

func (*NestedPointerChild) Test() string {
	return "Pointer test"
}

func Test_struct_nestedptrmethod(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	a := NestedPointer{}
	L.SetGlobal("a", New(L, &a))

	testReturn(t, L, `return a.b:Test()`, "Pointer test")
}

type TestStructEmbeddedType struct {
	TestStructEmbeddedTypeString
}

type TestStructEmbeddedTypeString string

func Test_struct_embeddedtype(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	a := TestStructEmbeddedType{
		TestStructEmbeddedTypeString: "hello",
	}

	L.SetGlobal("a", New(L, &a))

	testReturn(t, L, `a.TestStructEmbeddedTypeString = "world"`)

	if val := a.TestStructEmbeddedTypeString; val != "world" {
		t.Fatalf("expecting %s, got %s", "world", val)
	}
}

type TestStructEmbedded struct {
	StructTestPerson
	P  StructTestPerson
	P2 StructTestPerson `luar:"other"`
}

func Test_struct_embedded(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	e := &TestStructEmbedded{}
	L.SetGlobal("e", New(L, e))

	testReturn(
		t,
		L,
		`
		e.StructTestPerson = {
			Name = "Bill",
			Age = 33
		}
		e.P = {
			Name = "Tim",
			Age = 94,
			Friend = {
				Name = "Bob",
				Age = 77
			}
		}
		e.other = {
			Name = "Dale",
			Age = 26
		}
		`,
	)

	{
		expected := StructTestPerson{
			Name: "Bill",
			Age:  33,
		}
		if e.StructTestPerson != expected {
			t.Fatalf("expected %#v, got %#v", expected, e.StructTestPerson)
		}
	}

	{
		expected := StructTestPerson{
			Name: "Bob",
			Age:  77,
		}
		if *(e.P.Friend) != expected {
			t.Fatalf("expected %#v, got %#v", expected, *e.P.Friend)
		}
	}

	{
		expected := StructTestPerson{
			Name: "Dale",
			Age:  26,
		}
		if e.P2 != expected {
			t.Fatalf("expected %#v, got %#v", expected, e.P2)
		}
	}
}

type TestPointerReplaceHidden struct {
	A string `luar:"q"`
	B int    `luar:"other"`
	C int    `luar:"-"`
}

func Test_struct_pointerreplacehidden(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	e := &TestPointerReplaceHidden{}
	L.SetGlobal("e", New(L, e))

	testReturn(
		t,
		L,
		`
		_ = e ^ {
			q = "Cat",
			other = 675
		}
		`,
	)

	expected := TestPointerReplaceHidden{
		A: "Cat",
		B: 675,
	}

	if *e != expected {
		t.Fatalf("expected %v, got %v", expected, *e)
	}

	testError(
		t,
		L,
		`
		_ = e ^ {
			C = 333
		}
		`,
		"unable to set pointer value",
	)
}

func Test_struct_immutable_Edit(t *testing.T) {
	// Modifying a field on an immutable struct - should error
	L := lua.NewState()
	defer L.Close()

	p := StructTestPerson{
		Name: "Tim",
		Age:  66,
	}

	L.SetGlobal("p", New(L, p, ReflectOptions{Immutable: true}))

	testError(t, L, `p.Name = "Tom"`, "invalid operation on immutable struct")
}

func Test_struct_immutable_ptrfunc(t *testing.T) {
	// Calling a pointer function on an immutable struct - should error
	L := lua.NewState()
	defer L.Close()

	p := StructTestPerson{
		Name: "Tim",
		Age:  66,
	}

	L.SetGlobal("p", New(L, &p, ReflectOptions{Immutable: true}))

	testError(t, L, `p:IncreaseAge()`, "cannot call pointer methods on immutable objects")
}

func Test_struct_immutable_access(t *testing.T) {
	// Accessing a field and calling a regular function on an immutable
	// struct - should be fine
	L := lua.NewState()
	defer L.Close()

	p := StructTestPerson{
		Name: "Tim",
		Age:  66,
	}

	L.SetGlobal("p", New(L, p, ReflectOptions{Immutable: true}))

	testReturn(t, L, `return p:Hello()`, "Hello, Tim")
	testReturn(t, L, `return p.Name`, "Tim")
}


func Test_struct_immutable_nested(t *testing.T) {
	// Attempt to modify a nested field on an immutable struct - should error
	L := lua.NewState()
	defer L.Close()

	f := StructTestFamily{
		Mother: StructTestPerson{
			Name: "Luara",
		},
		Father: StructTestPerson{
			Name: "Tim",
		},
	}

	L.SetGlobal("f", New(L, f, ReflectOptions{Immutable: true}))

	testError(t, L, `f.Mother.Name = "Laura"`, "invalid operation on immutable struct")
}

func Test_struct_immutable_nestedvar(t *testing.T) {
	// Assign a nested struct field to a variable - should inherit
	// parent's immutable setting and cause error
	L := lua.NewState()
	defer L.Close()

	f := StructTestFamily{
		Mother: StructTestPerson{
			Name: "Luara",
		},
		Father: StructTestPerson{
			Name: "Tim",
		},
	}

	L.SetGlobal("f", New(L, f, ReflectOptions{Immutable: true}))

	testError(
		t,
		L,
		`
		mother = f.Mother
		mother.Name = "Laura"
		`,
		"invalid operation on immutable struct",
	)
}

func Test_struct_immutable_nestedslice(t *testing.T) {
	// Attempt to modify a nested field in a nested slice, on an immutable
	// struct - should error
	L := lua.NewState()
	defer L.Close()

	f := StructTestFamily{
		Mother: StructTestPerson{
			Name: "Luara",
		},
		Father: StructTestPerson{
			Name: "Tim",
		},
		Children: []StructTestPerson{
			{Name: "Bill"},
		},
	}

	L.SetGlobal("f", New(L, f, ReflectOptions{Immutable: true}))

	testError(t, L, `f.Children[1].Name = "Bill"`, "invalid operation on immutable struct")
}

func Test_struct_immutable_ptrnestedslice(t *testing.T) {
	// Attempt to modify a nested field in a nested slice, on an immutable
	// struct pointer - should error
	L := lua.NewState()
	defer L.Close()

	f := StructTestFamily{
		Mother: StructTestPerson{
			Name: "Luara",
		},
		Father: StructTestPerson{
			Name: "Tim",
		},
		Children: []StructTestPerson{
			{Name: "Bill"},
		},
	}

	L.SetGlobal("f", New(L, &f, ReflectOptions{Immutable: true}))

	testError(t, L, `f.Children[1].Name = "Bill"`, "invalid operation on immutable struct")
}


type TestTransparentPtrAccessB struct {
	Str *string
}

type TestTransparentPtrAccessA struct {
	B *TestTransparentPtrAccessB
}

func Test_struct_tranparentptr_access(t *testing.T) {
	// Access an undefined pointer field - should auto populate with zero
	// value as if a non-pointer object
	L := lua.NewState()
	defer L.Close()

	val := "foo"
	b := TestTransparentPtrAccessB{}
	b.Str = &val

	L.SetGlobal("b", New(L, &b, ReflectOptions{TransparentPointers: true}))

	testReturn(t, L, `return b.Str`, "foo")
}

func Test_struct_transparentptr_varassign(t *testing.T) {
	// Assign one pointer value to another, with the left side
	// transparent - requires indirection of the right side since
	// the left behaves like a non-pointer field. They should
	// also be separate objects at that point - no shared address.
	// This is distinct from regular pointer assignment, where
	// modifying a value would change it for both references.
	L := lua.NewState()
	defer L.Close()

	val := "assigned ptr value"
	a := TestTransparentPtrAccessA{}
	b := TestTransparentPtrAccessB{
		Str: &val,
	}
	L.SetGlobal("a", New(L, &a, ReflectOptions{TransparentPointers: true}))
	L.SetGlobal("b", New(L, &b, ReflectOptions{TransparentPointers: true}))

	testReturn(
		t,
		L,
		`
		a.B = -b
		return a.B.Str
		`,
		"assigned ptr value",
	)
	testReturn(
		t,
		L,
		`
		b.Str = "new value"
		return b.Str
		`,
		"new value",
	)
	testReturn(t, L, `return a.B.Str`, "assigned ptr value")
}

func Test_struct_transparentptr_assign(t *testing.T) {
	// Assign a non-pointer struct value to a pointer field -
	// should be fine
	L := lua.NewState()
	defer L.Close()

	val := "assigned ptr value"
	a := TestTransparentPtrAccessA{}
	b := TestTransparentPtrAccessB{
		Str: &val,
	}
	L.SetGlobal("a", New(L, &a, ReflectOptions{TransparentPointers: true}))
	// Non-pointer
	L.SetGlobal("b", New(L, b, ReflectOptions{TransparentPointers: true}))

	testReturn(
		t,
		L,
		`
		a.B = b
		return a.B.Str
		`,
		"assigned ptr value",
	)
	testReturn(t, L, `return b.Str`, "assigned ptr value")
}

func Test_struct_transparentptr_assigninvalidtype(t *testing.T) {
	// Assign an invalid type to a given field - should return
	// a meaningful error
	L := lua.NewState()
	defer L.Close()

	b1 := &TestTransparentPtrAccessB{}
	b2 := &TestTransparentPtrAccessB{}

	L.SetGlobal("b1", New(L, b1, ReflectOptions{TransparentPointers: true}))
	L.SetGlobal("b2", New(L, b2, ReflectOptions{TransparentPointers: true}))

	testError(t, L, `b1.Str = b2`, "could not set field Str: expected type string")
}

func Test_struct_transparentptr_nonptrvar(t *testing.T) {
	// Attempt to access a nil pointer field on a transparent pointer
	// struct that was reflected by value. Since we can't actually set
	// values back to a struct that was reflected by value (as
	// opposed to by reference), an error will result.
	L := lua.NewState()
	defer L.Close()

	b := TestTransparentPtrAccessB{}

	// Non-pointer
	L.SetGlobal("b", New(L, b, ReflectOptions{TransparentPointers: true}))

	testError(t, L, `print(b.Str)`, "cannot transparently create pointer field Str")
}

func Test_struct_transparentptr_autopop(t *testing.T) {
	// Access an undefined nested pointer field - should auto populate
	// with zero values as if a non-pointer object
	L := lua.NewState()
	defer L.Close()

	a := TestTransparentPtrAccessA{}

	L.SetGlobal("a", New(L, &a, ReflectOptions{
		TransparentPointers: true, AutoPopulate: true}))

	testReturn(t, L, `return a.B.Str`, "")
}

func Test_struct_transparentptr_nestedassign(t *testing.T) {
	// Set an undefined nested pointer field - should get assigned like
	// a regular non-pointer field
	L := lua.NewState()
	defer L.Close()

	a := TestTransparentPtrAccessA{}

	L.SetGlobal("a", New(L, &a, ReflectOptions{
		TransparentPointers: true, AutoPopulate: true}))

	testReturn(
		t,
		L,
		`
		a.B.Str = "hello, world!"
		return a.B.Str
		`,
		"hello, world!",
	)
}

func Test_struct_transparentptr_equality(t *testing.T) {
	// Check equality on a pointer field - should act like a plain field
	L := lua.NewState()
	defer L.Close()

	b := TestTransparentPtrAccessB{}
	val := "foo"
	b.Str = &val

	L.SetGlobal("b", New(L, &b, ReflectOptions{TransparentPointers: true}))

	testReturn(t, L, `return b.Str == "foo"`, "true")
}

func Test_struct_transparentptr_pow(t *testing.T) {
	// Access a pointer field in the normal pointer way - should error
	L := lua.NewState()
	defer L.Close()

	b := TestTransparentPtrAccessB{}
	val := "foo"
	b.Str = &val

	L.SetGlobal("b", New(L, &b, ReflectOptions{TransparentPointers: true}))

	testError(t, L, `_ = b.Str ^ "hello"`, "cannot perform pow operation between string and string")
}

type TestTransparentStructSliceFieldA struct {
	List []string
}

func Test_struct_transparentptr_sliceautopop(t *testing.T) {
	// Access an undefined slice field - should be automatically created
	L := lua.NewState()
	defer L.Close()

	a := TestTransparentStructSliceFieldA{}

	L.SetGlobal("a", New(L, &a, ReflectOptions{
		TransparentPointers: true, AutoPopulate: true}))

	testReturn(t, L, `return #a.List`, "0")
}

func Test_struct_transparentptr_sliceautopopappend(t *testing.T) {
	// Append to an undefined slice field - should be fine
	L := lua.NewState()
	defer L.Close()

	a := TestTransparentStructSliceFieldA{}

	L.SetGlobal("a", New(L, &a, ReflectOptions{
		TransparentPointers: true, AutoPopulate: true}))

	testReturn(
		t,
		L,
		`
		a.List = a.List:append("hi")
		return #a.List
		`,
		"1",
	)
}

func Test_struct_transparentptr_var(t *testing.T) {
	// Assign the value of a pointer field to a variable - variable should
	// inherit the transparent reflect options
	L := lua.NewState()
	defer L.Close()

	val := "hello, world!"
	a := TestTransparentPtrAccessA{
		&TestTransparentPtrAccessB{&val},
	}

	L.SetGlobal("a", New(L, &a, ReflectOptions{TransparentPointers: true}))

	testReturn(
		t,
		L,
		`
		b = a.B
		return b.Str
		`,
		"hello, world!",
	)
}

func Test_struct_transparentptr_sliceval(t *testing.T) {
	// Assign a slice value to a variable - variable should inherit the
	// transparent reflect options
	val := "hello, world!"
	L := lua.NewState()
	defer L.Close()

	list := []TestTransparentPtrAccessA{
		{&TestTransparentPtrAccessB{&val}},
	}

	L.SetGlobal("list", New(L, &list, ReflectOptions{TransparentPointers: true}))

	testReturn(
		t,
		L,
		`
		a = list[1]
		return a.B.Str
		`,
		"hello, world!",
	)
}

func Test_struct_immutable_transparentptraccess(t *testing.T) {
	// Access a transparent pointer field on an immutable
	// struct - should be fine
	L := lua.NewState()
	defer L.Close()

	val := "foo"
	b := TestTransparentPtrAccessB{}
	b.Str = &val

	L.SetGlobal("b", New(L, &b, ReflectOptions{Immutable: true, TransparentPointers: true}))

	testReturn(t, L, `return b.Str`, "foo")
}

func Test_struct_immutable_edit(t *testing.T) {
	// Attempt to modify a transparent pointer field on an
	// immutable struct - should error
	L := lua.NewState()
	defer L.Close()

	val := "foo"
	b := TestTransparentPtrAccessB{}
	b.Str = &val

	L.SetGlobal("b", New(L, &b, ReflectOptions{Immutable: true, TransparentPointers: true}))

	testError(t, L, `b.Str = "bar"`, "invalid operation on immutable struct")
}

func Test_struct_immutable_nonautopop(t *testing.T) {
	// Attempt to read from nil on an immutable struct with transparent
	// pointers.  This should still return nil, since we should not modify the
	// immutable struct.
	L := lua.NewState()
	defer L.Close()

	x := struct {
		Ptr  *int
		List []int
	}{}
	L.SetGlobal("x", New(L, &x, ReflectOptions{TransparentPointers: true}))

	testReturn(t, L, `return x.Ptr`, "nil")
	testReturn(t, L, `return x.List`, "nil")

	if x.Ptr != nil || x.List != nil {
		t.Error("immutable struct was edited", x)
	}
}

func Test_struct_differentreflectopts(t *testing.T) {
	// Set up two structs with different ReflectOptions -
	// should retain their own behaviors
	L := lua.NewState()
	defer L.Close()

	val := "hello"
	a := TestTransparentPtrAccessB{Str: &val}
	b := TestTransparentPtrAccessB{Str: &val}

	L.SetGlobal("a", New(L, &a, ReflectOptions{Immutable: true}))
	L.SetGlobal("b", New(L, &b, ReflectOptions{TransparentPointers: true}))

	testReturn(t, L, `return -a.Str`, "hello")
	testReturn(t, L, `return b.Str`, "hello")
	testReturn(
		t,
		L,
		`
		b.Str = "world"
		return b.Str
		`,
		"world",
	)
	testReturn(t, L, `return -a.Str`, "hello")
}
