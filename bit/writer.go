package bit

// Writer implements bit-level writing operations into an int32.
// Writer panics on invalid data or out-of-bounds bit writes.
type Writer struct {
	out *int32
	pos uint8
}

// NewWriter creates a new Writer using the int32 as the underlying destination to write bits into.
func NewWriter(out *int32) *Writer {
	*out = 0
	return &Writer{out: out, pos: 0}
}

// Pad advances the write cursor by the given number of bits.
// Skipped bits remain zero.
func (w *Writer) Pad(bits uint8) {
	w.pos += bits
}

// Uint8 writes a fixed number of bits from the uint8 value pointed to by x.
// bits specifies the number of bits to write, which must be in [1, 8].
func (w *Writer) Uint8(x *uint8, bits uint8) {
	mask := uint8(1<<bits - 1)
	v := int32(*x & mask)
	*w.out |= v << w.pos
	w.pos += bits
}

// Int32 writes a fixed number of bits from the int32 value pointed to by x.
// bits specifies the number of bits to write, which must be in [1, 32].
func (w *Writer) Int32(x *int32, bits uint8) {
	mask := int32(1<<bits - 1)
	v := *x & mask
	*w.out |= v << w.pos
	w.pos += bits
}

// Bool writes a single bit from the bool value pointed to by x.
func (w *Writer) Bool(x *bool) {
	if *x {
		*w.out |= 1 << w.pos
	}
	w.pos++
}
