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
	for i, bit := range bits {
		var wantWritten int
		if i%16 == 15 {
			wantWritten = 16
		}
		written, err := w.WriteBit(bit)
		if written != wantWritten {
			t.Errorf(`Writing bit %d: written = %v, want %v`, i, written, wantWritten)
		}
		if err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}
	}
	written, err := w.Commit()
	if written != 16 {
		t.Errorf(`written = %v, want 16`, written)
	}
	if err != nil {
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
	lim := limitWriter(5)
	w := NewWriter(&lim, 2)
	for i := 0; i < 4*8; i++ {
		if _, err := w.WriteBit(Bit(i % 2)); err != nil {
			t.Fatalf(`Writing bit %d: err = %v, want nil`, i, err)
		}
	}
	t.Logf(`limit writer before overwrite: %d`, lim)
	for i := 0; i < 8; i++ {
		if _, err := w.WriteBit(Bit(i % 2)); err != nil {
			t.Fatalf(`Writing bit %d: err = %v, want nil`, i, err)
		}
	}
	w.WriteBit(1)
	written, err := w.Commit()
	t.Logf(`limit writer after overwrite: %d`, lim)
	if written != 8 {
		t.Errorf(`written = %v, want 8`, written)
	}
	if err != io.EOF {
		t.Errorf(`err = %v, want io.EOF`, err)
	}
}

// TestWriterCommitErr checks that committing chunks can return an error.
func TestWriterCommitErr(t *testing.T) {
	lim := limitWriter(4)
	w := NewWriter(&lim, 5)
	written, err := w.Commit()
	if written != 32 {
		t.Errorf(`written = %v, want 32`, written)
	}
	if err != io.EOF {
		t.Fatalf(`err = %v, want io.EOF`, err)
	}
}

// TestWriteBits checks that multiple bits can be written.
func TestWriteBits(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b, 2)
	var bits Bits = 0xA5F

	if _, err := w.WriteBits(bits, 12); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if n, err := w.Commit(); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	} else if n != 16 {
		t.Fatalf(`n = %v, want 16`, n)
	}

	tests := []byte{0xA5, 0xF0}

	for i, actual := range b.Bytes() {
		if want := tests[i]; actual != want {
			t.Errorf(`byte %d = %02X, want %02X`, i, actual, want)
		}
	}
}

// TestCommitPending checks that bits will only be written if they are pending.
func TestCommitPending(t *testing.T) {
	tests := []struct {
		bits        Bits
		len         int
		wantWritten int
		wantByte    byte
	}{
		{0xA, 4, 8, 0xA0},
		{0xA0, 8, 0, 0xA0},
	}

	for i, tt := range tests {
		var b bytes.Buffer

		t.Logf(`TestCommitPending %d`, i)
		w := NewWriter(&b, 1)

		w.WriteBits(tt.bits, tt.len)

		written, err := w.CommitPending()

		if err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}

		if written != tt.wantWritten {
			t.Errorf(`written = %v, want %v`, written, tt.wantWritten)
		}
		if actual, _ := b.ReadByte(); actual != tt.wantByte {
			t.Errorf(`written byte = %02X, want %02X`, actual, tt.wantByte)
		}
	}
}

// TestHasPendingBits checks that it returns true if bits are pending, false
// otherwise.
func TestHasPendingBits(t *testing.T) {
	var b bytes.Buffer

	w := NewWriter(&b, 1)

	tests := []struct {
		test func()
		want bool
	}{
		{func() {}, false},
		{func() { w.WriteBit(1) }, true},
		{func() { w.WriteBits(0, 3) }, true},
		{func() { w.Commit() }, false},
	}

	for i, tt := range tests {
		t.Logf(`TestHasPendingBits %d`, i)
		tt.test()
		if actual := w.HasPendingBits(); actual != tt.want {
			t.Errorf(`HasPendingBits() = %v, want %v`, actual, tt.want)
		}
	}
}

// TestWriteBitsErr checks that an error will be returned if it occurs when
// writing bits.
func TestWriteBitsErr(t *testing.T) {
	lim := limitWriter(1)
	w := NewWriter(&lim, 2)
	if _, err := w.WriteBits(0, 8); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if _, err := w.WriteBits(0, 8); err != io.EOF {
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
