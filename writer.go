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
func NewWriter(r io.Writer, chunkSize int) *Writer {
	return nil
}

// WriteBit writes a single bit.
func (w *Writer) WriteBit() error {
	return nil
}

// Commit commits the current bytes to the writer, even if a byte is only
// partially written. Partial bytes will be zero-filled. A commit will happen
// any time that a write would overflow the current chunk. A single commit when
// writing is finished is recommended.
func (w *Writer) Commit() error {
	return nil
}
