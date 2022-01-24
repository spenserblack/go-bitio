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
		fmt.Printf("0x%3X\n", bits)
	}
	// Output:
	// 0x123
	// 0x456
}
