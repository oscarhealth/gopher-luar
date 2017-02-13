# gopher-luar [![GoDoc](https://godoc.org/github.com/oscarhealth/gopher-luar?status.svg)](https://godoc.org/github.com/oscarhealth/gopher-luar)

custom type reflection for [gopher-lua](https://github.com/yuin/gopher-lua).

## oscarhealth fork

This fork includes additional features targeted at modifying the behavior of reflected objects. In particular,
the following are supported:

* **Immutability**: Various types can be reflected as immutable, preventing any modification in Lua. This includes
modification of struct fields, slices, maps, pointers, etc. A channel is prevented from being closed.
* **Transparent pointers**: Objects can be reflected as if all pointer fields were plain value fields - removing
  the need for the `^` and `-` operators that gopher-luar typically requires for manipulating pointers.
* **Automatic population** On top of transparent pointers, objects can optionally have nil fields automatically
  created on read or write. This allows the lua code to act as if there were no pointers in the way, and they just
  had the zero values.

See the documentation for usage.

## License

MPL 2.0
