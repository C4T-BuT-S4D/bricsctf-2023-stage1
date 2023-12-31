package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/klauspost/compress/flate"
)

func main() {
	compressed := []byte{ /* Packet 147 */
		0x74, 0x91, 0xcb, 0x6e, 0xe2, 0x48, 0x18, 0x85,
		0x33, 0xea, 0x9d, 0x25, 0xbf, 0x43, 0x6d, 0xa6,
		0xd5, 0xdd, 0x83, 0x6f, 0x24, 0x04, 0x30, 0x62,
		0xe1, 0x86, 0x84, 0x4b, 0xb0, 0x13, 0x88, 0x21,
		0xc1, 0x9b, 0xa8, 0xb0, 0x7f, 0xdb, 0x65, 0xd7,
		0xc5, 0x29, 0x17, 0x60, 0xf2, 0x8a, 0xb3, 0x99,
		0x47, 0x1a, 0x05, 0x25, 0x8b, 0x1e, 0x4d, 0xd7,
		0xea, 0xd3, 0x5f, 0x75, 0xea, 0x9c, 0x5f, 0xe7,
		0xe2, 0xcb, 0x3f, 0x17, 0xef, 0xe7, 0x8f, 0x8b,
		0x8b, 0x8b, 0x2f, 0x7f, 0x3f, 0xdc, 0x3f, 0x86,
		0xc8, 0xaa, 0x49, 0xc6, 0x09, 0x47, 0xd3, 0x30,
		0x7c, 0xb0, 0x1c, 0xd3, 0xd1, 0xb5, 0xa9, 0xa8,
		0x95, 0x8b, 0x9c, 0x7e, 0xdb, 0x74, 0xae, 0x7b,
		0xa6, 0x6d, 0x3a, 0x6d, 0xb7, 0x6f, 0xdb, 0xb6,
		0xae, 0x8d, 0x04, 0xe7, 0x10, 0x2b, 0x22, 0xb8,
		0x8b, 0x4a, 0x80, 0xca, 0xc0, 0x94, 0x1c, 0x40,
		0xd7, 0x46, 0x82, 0x2b, 0xe0, 0xca, 0x58, 0x00,
		0xcf, 0x54, 0xee, 0x22, 0xa7, 0xd7, 0xd6, 0xb5,
		0x11, 0x8e, 0x73, 0x30, 0x46, 0x82, 0x2b, 0x29,
		0xa8, 0x8b, 0x18, 0x6e, 0x0c, 0x9c, 0xc1, 0xd0,
		0xd6, 0xb5, 0x7b, 0x49, 0x32, 0xc2, 0x5d, 0x94,
		0x2b, 0x55, 0xb9, 0x96, 0xf5, 0x3f, 0x46, 0xe3,
		0x20, 0x74, 0x91, 0xa3, 0x6b, 0xeb, 0x2a, 0x93,
		0x38, 0x01, 0x63, 0xc6, 0x6b, 0x88, 0xf7, 0x12,
		0x8c, 0x15, 0xbc, 0xee, 0xa1, 0x56, 0xf5, 0xf9,
		0xf6, 0xd3, 0x36, 0x3c, 0x55, 0xe0, 0x22, 0x5c,
		0x55, 0x94, 0xc4, 0x58, 0x11, 0xc1, 0xad, 0xc6,
		0x38, 0x1e, 0x8f, 0x46, 0x2a, 0x24, 0x33, 0xf6,
		0x92, 0x02, 0x8f, 0x45, 0x02, 0x89, 0xae, 0xad,
		0x6b, 0x90, 0x86, 0x97}
	compressed = append(compressed, []byte{ /* Packet 147 */
		0x01, 0x57, 0x2e, 0xf2}...)
	compressed = append(compressed, []byte{ /* Packet 147 */
		0xc5, 0x1b, 0xa1, 0x14, 0x5b, 0x1d, 0xd3, 0x46,
		0xdf, 0x16, 0x84, 0xef, 0x9b, 0x01, 0xf2, 0x78,
		0x22, 0x05, 0x49, 0x90, 0x63, 0x0f, 0xd0, 0xdd,
		0x77, 0xe4, 0x55, 0x15, 0x85, 0x27, 0xd8, 0xdd,
		0x11, 0x65, 0x75, 0x2e, 0xbb, 0xe6, 0xe5, 0x35,
		0xfa, 0x76, 0x37, 0x0d, 0xfd, 0x45, 0x0b, 0x51,
		0x52, 0x02, 0x9a, 0x40, 0x5c, 0x8a, 0xef, 0x68,
		0x94, 0x4b, 0xc1, 0xc0, 0x72, 0x9c, 0xae, 0x69,
		0x9b, 0xb6, 0x69, 0x23, 0x5f, 0xec, 0x08, 0x05,
		0xf4, 0x88, 0x53, 0x2c, 0xc9, 0x87, 0x52, 0xd7,
		0xbc, 0x38, 0x86, 0x4a, 0xb9, 0x48, 0x41, 0xa3,
		0xac, 0x5c, 0x31, 0xda, 0xfa, 0x25, 0xf1, 0xfb,
		0xe4, 0xaf, 0xe6, 0xbf, 0x53, 0x46, 0x07, 0xaf,
		0x43, 0xdb, 0xec, 0xb7, 0x08, 0xc3, 0x19, 0x58,
		0xf8, 0x40, 0xd2, 0x0f, 0x3c, 0xc2, 0xae, 0xfa,
		0x40, 0x5c, 0xf1, 0xac, 0xf5, 0xc3, 0xfa, 0x71,
		0x7e, 0xda, 0xfb, 0xe5, 0xdb, 0xf7, 0x6e, 0x21,
		0x31, 0xa0, 0x89, 0x73, 0xcc, 0x33, 0x18, 0x1c,
		0x86, 0xbb, 0xcb, 0xc1, 0xeb, 0xd0, 0x36, 0xbb,
		0xba, 0xb6, 0x82, 0x14, 0x24, 0xc8, 0xdf, 0xf7,
		0x70, 0x56, 0x13, 0xfe, 0x19, 0xdd, 0xb8, 0xe1,
		0xb1, 0x48, 0x08, 0xcf, 0x5c, 0x94, 0xbd, 0x91,
		0xaa, 0x85, 0x12, 0x48, 0x29, 0x56, 0xf0, 0xb9,
		0x9a, 0xb1, 0xc0, 0x3c, 0xdb, 0xe3, 0x0c, 0x5c,
		0x04, 0xdc, 0x58, 0x3f, 0xb6, 0x80, 0x9f, 0xad,
		0xfa, 0x2d, 0xb9, 0x37, 0x56, 0xeb, 0x33, 0xf7,
		0x5a, 0x72, 0x7f, 0x86, 0xee, 0x7b, 0x7d, 0xa2,
		0x24, 0xe0, 0xa2, 0x1a, 0xea, 0x9a, 0x08, 0x3e,
		0x84, 0xd3, 0x3c, 0x8d, 0xf8, 0x9c, 0xc6, 0xed,
		0x8c, 0xdc, 0xb3, 0x28}...)
	compressed = append(compressed, []byte{ /* Packet 147 */
		0xdf, 0x4d, 0x03, 0xba}...)
	compressed = append(compressed, []byte{ /* Packet 147 */
		0x18, 0xcd, 0x8b, 0xf8, 0x72, 0xce, 0x9e, 0x2f,
		0x57, 0x07, 0xdc, 0xde, 0xec, 0x67, 0x85, 0x20,
		0x41, 0x7b, 0x45, 0xfc, 0xf1, 0x86, 0x06, 0xe3,
		0xd9, 0x29, 0x08, 0x7f, 0xe6, 0xc1, 0x78, 0xdb,
		0x09, 0x8a, 0x80, 0xfa, 0xcc, 0x3f, 0x6e, 0x8b,
		0xf2, 0x14, 0x3d, 0xdd, 0x38, 0x41, 0xe1, 0xb7,
		0xb7, 0x85, 0xd7, 0x44, 0x4f, 0xdb, 0xc6, 0x67,
		0xb7, 0x2c, 0x22, 0xf3, 0xbe, 0x19, 0x2d, 0xd5,
		0x82, 0x79, 0x66, 0x06, 0x57, 0x93, 0x29, 0x09,
		0xc2, 0xdd, 0xd5, 0x0b, 0x7b, 0xd8, 0x85, 0x55,
		0x38, 0xeb, 0x76, 0x3d, 0x3e, 0x6e, 0x8e, 0x9e,
		0xa7, 0x6b, 0xba, 0x06, 0x0c, 0x13, 0x3a, 0xc4,
		0x87, 0x8c, 0xc2, 0x55, 0x09, 0xbc, 0xb0, 0x4f,
		0x20, 0xcd, 0x72, 0x35, 0xf1, 0xca, 0x5d, 0x1a,
		0xab, 0xf1, 0x2d, 0x9d, 0x79, 0x7f, 0x5e, 0xd9,
		0x95, 0x14, 0x4a, 0x70, 0x93, 0xc1, 0xd7, 0x0a,
		0xd7, 0xf5, 0x51, 0xc8, 0x64, 0x78, 0xaa, 0x57,
		0xd7, 0x91, 0x13, 0xf2, 0x9f, 0xf5, 0xed, 0x63,
		0x67, 0xf9, 0x2a, 0x9e, 0xfd, 0xcd, 0x55, 0xb4,
		0xfc, 0x1a, 0xd7, 0x32, 0x7d, 0x51, 0xa2, 0x04,
		0x3e, 0x9c, 0x15, 0x49, 0xb9, 0x2d, 0x3c, 0x27,
		0x0a, 0x97, 0x27, 0xbf, 0x58, 0x1f, 0xb7, 0xe1,
		0xb2, 0x7d, 0x1f, 0x6e, 0xdf, 0xa2, 0x70, 0x5e,
		0xf8, 0x93, 0x59, 0xc7, 0x67, 0x9b, 0x3c, 0x08,
		0xb7, 0x6f, 0x01, 0x9b, 0x1d, 0xfd, 0xa7, 0x0d,
		0xf3, 0xc3, 0x79, 0x1e, 0xb1, 0x2d, 0x31, 0xa3,
		0xa5, 0x4a, 0xe7, 0x99, 0x39, 0x51, 0xf1, 0x2e,
		0xba, 0x93, 0xcd, 0xf3, 0xcb, 0x24, 0x0c, 0x0e,
		0xde, 0x8d, 0x1d, 0xb5, 0x67, 0x54, 0x1e, 0x5e,
		0xd2}...)
	compressed = append(compressed, []byte{ /* Packet 147 */
		0x5e, 0xf9, 0x2f, 0x00}...)

	// by rfc
	compressed = append(compressed, []byte{
		0x00, 0x00, 0xff, 0xff}...)

	reader := bytes.NewReader(compressed)
	dict := make([]byte, 32768)
	fr := flate.NewReaderDict(reader, dict)

	buf := make([]byte, 1)
	for {
		n, err := fr.Read(buf)
		if n > 0 {
			fmt.Print(hex.EncodeToString(buf))
		}

		if err != nil {
			fmt.Println()
			fmt.Println("ended on error", err)
			break
		}
	}
}
