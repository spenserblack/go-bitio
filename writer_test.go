package bitio

import (
	"bytes"
	"io"
	"testing"
)

// TestNewWriter tests that a new writer is created with the proper values.
func TestNewWriter(t *testing.T) {
	w := NewWriter(&bytes.Buffer{}, 2)
	if _, ok := w.w.(*bytes.Buffer); !ok {
		t.Fatalf(`w.w is %T, want *bytes.Buffer`, w.w)
	}
	if l := len(w.bytes); l != 2 {
		t.Fatalf(`len(w.bytes) = %v, want 2`, l)
	}
	if w.index != 0 {
		t.Fatalf(`w.index = %v, want 0`, w.index)
	}
}

// TestWriteAndCommitBit checks that bits can be written.
func TestWriteAndCommitBit(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b, 2)
	bits := []Bit{
		0, 1, 0, 1, 1, 0, 1, 0,
		1, 1, 1, 1,
	}
	for _, bit := range bits {
		if err := w.WriteBit(bit); err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}
	}
	if err := w.Commit(); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	expected := []byte{0b01011010, 0b11110000}

	for i, actual := range b.Bytes() {
		if want := expected[i]; actual != want {
			t.Errorf(`Byte %d = %08b, want %08b`, i, actual, want)
		}
	}
}

// TestWriteBitErr checks that WriteBit can return an error.
func TestWriteBitErr(t *testing.T) {
	lim := limitWriter(4)
	w := NewWriter(&lim, 2)
	for i := 0; i < 4*8; i++ {
		if err := w.WriteBit(Bit(i % 2)); err != nil {
			t.Fatalf(`Writing bit %d: err = %v, want nil`, i, err)
		}
	}
	for i := 0; i < 8; i++ {
		if err := w.WriteBit(Bit(i % 2)); err != nil {
			t.Fatalf(`Writing bit %d: err = %v, want nil`, i, err)
		}
	}
	if err := w.WriteBit(1); err != io.EOF {
		t.Fatalf(`err = %v, want io.EOF`, err)
	}
}

// TestWriterCommitErr checks that committing chunks can return an error.
func TestWriterCommitErr(t *testing.T) {
	lim := limitWriter(4)
	w := NewWriter(&lim, 5)
	if err := w.Commit(); err != io.EOF {
		t.Fatalf(`err = %v, want io.EOF`, err)
	}
}

type limitWriter int

func (w *limitWriter) Write(p []byte) (n int, err error) {
	for n = range p {
		if *w <= 0 {
			err = io.EOF
			return
		}
		*w--
	}
	return
}