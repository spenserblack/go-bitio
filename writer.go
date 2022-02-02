package bitio

import "io"

// Writer writes bits to an io.Writer
//
// Assuming bytes are little-endian, writing occurs from left to right.
type Writer struct {
	// W is the underlying writer that writes bytes.
	w io.Writer
	// Bytes stores the bytes that are written to the reader. When the last
	// available bit is written, Bytes will be written and then emptied.
	bytes []byte
	// Index is the index of the bit being written.
	index int
}

// NewWriter creates a new bit writer. The amount of full bytes to write at a
// time is set by chunkSize.
func NewWriter(w io.Writer, chunkSize int) *Writer {
	return &Writer{
		w:     w,
		bytes: make([]byte, chunkSize),
	}
}

// WriteBit writes a single bit.
//
// The number of bits written will be returned, which will be 0 until a chunk
// is filled.
func (w *Writer) WriteBit(b Bit) (written int, err error) {
	w.bytes[w.byteIndex()] |= b << (7 - w.bitIndex())
	// NOTE Just defining bool in var for clarity of purpose
	if wrapped := w.incIndex(); wrapped {
		written, err = w.Commit()
	}
	return
}

// WriteBits writes multiple bits from an int.
//
// Bits will be interpreted from the left to the right bit (assuming
// the int is little-endian).
//
// Length is used to specify the number of bits to write, to remove ambiguity
// between an "empty" set of bits and a long string of 0s. Length specifies
// the left-most bit.
//
// The number of bits written will be returned, which will be 0 if a chunk
// wasn't filled.
func (w *Writer) WriteBits(bits Bits, length int) (written int, err error) {
	for i := 1; i <= length; i++ {
		var writtenBits int
		shift := length - i
		b := (bits >> shift) & 1
		writtenBits, err = w.WriteBit(Bit(b))
		written += writtenBits
		if err != nil {
			return
		}
	}
	return
}

// Commit commits the current bytes to the writer, even if a byte is only
// partially written. Partial bytes will be zero-filled. A commit will happen
// any time that a write would overflow the current chunk.
//
// The number of bits written are returned, and any error that occurred when
// writing.
func (w *Writer) Commit() (written int, err error) {
	written, err = w.w.Write(w.bytes)
	written *= byteSize
	w.bytes = make([]byte, len(w.bytes))
	w.index = 0
	return
}

// CommitPending is a helper function that commits the current written bits
// only if the byte chunk is partially written. Does nothing if the byte chunk
// is empty.
//
// If it is unknown if the number of bits written will completely fill all
// chunks, then it is recommended to execute this once to conclude writing.
func (w *Writer) CommitPending() (written int, err error) {
	if w.HasPendingBits() {
		return w.Commit()
	}
	return
}

// HasPendingBits returns true if there are any bits that are pending, but have
// not yet been written to the underlying writer. For example, if half of a
// byte is written, then there are 4 pending bits.
func (w Writer) HasPendingBits() bool {
	return w.index != 0
}

// ByteIndex gets the index of the current byte to write bits to.
func (w Writer) byteIndex() int {
	return w.index / byteSize
}

// BitIndex gets the index of the current bit in the byte to write to.
func (w Writer) bitIndex() int {
	return w.index % byteSize
}

// IncIndex increments the index and wraps it. Returns if the number was
// wrapped.
func (w *Writer) incIndex() bool {
	w.index++
	w.index %= len(w.bytes) * byteSize
	return w.index == 0
}
