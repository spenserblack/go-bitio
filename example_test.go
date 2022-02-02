package bitio_test

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spenserblack/go-bitio"
)

func ExampleReader_ReadBits() {
	buff := bytes.NewBuffer([]byte{0x12, 0x34, 0x56})
	r := bitio.NewReader(buff, 3)

	for i := 0; i < 2; i++ {
		bits, _, err := r.ReadBits(12)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("0x%03X\n", bits)
	}
	// Output:
	// 0x123
	// 0x456
}

func ExampleWriter_WriteBits() {
	var buff bytes.Buffer
	w := bitio.NewWriter(&buff, 3)

	if _, err := w.WriteBits(0x123456, 24); err != nil {
		log.Fatal(err)
	}

	for _, b := range buff.Bytes() {
		fmt.Printf("0x%02X\n", b)
	}
	// Output:
	// 0x12
	// 0x34
	// 0x56
}

// A Writer writes bits to the underlying writer, but only when a "chunk" is
// filled. To write a partial chunk (for example, half of a bit), the write
// must be committed.
func ExampleWriter() {
	var buff bytes.Buffer

	// A writer with a chunk size of 2 bytes.
	w := bitio.NewWriter(&buff, 2)

	// 12 1 bits written. A chunk requires 16 bits (2 bytes), so the write is currently
	// pending.
	written, _ := w.WriteBits(0xFFF, 12)
	fmt.Printf("%d bits written\n", written)

	// Now we commit the pending bits.
	written, _ = w.CommitPending()
	fmt.Printf("%d bits written\n", written)

	for i, b := range buff.Bytes() {
		fmt.Printf("byte %d = 0x%02X\n", i, b)
	}

	// Output:
	// 0 bits written
	// 16 bits written
	// byte 0 = 0xFF
	// byte 1 = 0xF0
}
