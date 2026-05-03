package bit

// IO represents a bit-level IO direction. Implementations of this interface are Reader and Writer.
// Reader reads bits from the input int32 into the pointers passed, whereas Writer writes the values
// the pointers point to into the output int32.
type IO interface {
	Pad(bits uint8)
	Uint8(x *uint8, bits uint8)
	Int32(x *int32, bits uint8)
	Bool(x *bool)
}
