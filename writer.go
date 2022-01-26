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
// The number if bits written will be returned, which will be 0 until a chunk
// is filled.
func (w *Writer) WriteBit(b Bit) (written int, err error) {
	w.bytes[w.byteIndex()] |= b << (7 - w.bitIndex())
	// NOTE Just defining bool in var for clarity of purpose
	if wrapped := w.incIndex(); wrapped {
		written, err = w.Commit()
	}
	return
}

// Commit commits the current bytes to the writer, even if a byte is only
// partially written. Partial bytes will be zero-filled. A commit will happen
// any time that a write would overflow the current chunk. A single commit when
// writing is finished is recommended.
//
// The number of bits written are returned, and any error that occurred when
// writing.
func (w *Writer) Commit() (written int, err error) {
	written, err = w.Write(w.bytes)
	w.bytes = make([]byte, len(w.bytes))
	return
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
