package bit

// Reader implements bit-level reading operations from an int32.
// Reader panics on invalid data or out-of-bounds bit reads.
type Reader struct {
	in  *int32
	pos uint8
}

// NewReader creates a new Reader using the int32 as the underlying source to read bits from.
func NewReader(in *int32) *Reader {
	return &Reader{in: in, pos: 0}
}

// Pad skips the given number of bits.
func (r *Reader) Pad(bits uint8) {
	r.pos += bits
}

// Uint8 reads a fixed number of bits into the uint8 value pointed to by x.
// bits specifies the number of bits to read, which must be in [1, 8].
func (r *Reader) Uint8(x *uint8, bits uint8) {
	mask := uint8(1<<bits - 1)
	*x = uint8(*r.in>>r.pos) & mask
	r.pos += bits
}

// Int32 reads a fixed number of bits into the int32 value pointed to by x.
// bits specifies the number of bits to read, which must be in [1, 32].
func (r *Reader) Int32(x *int32, bits uint8) {
	mask := int32(1<<bits - 1)
	*x = (*r.in >> r.pos) & mask
	r.pos += bits
}

// Bool reads a single bit as a bool into the value pointed to by x.
func (r *Reader) Bool(x *bool) {
	*x = ((*r.in >> r.pos) & 1) != 0
	r.pos++
}
