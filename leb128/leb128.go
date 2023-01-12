package leb128

import "fmt"

func Int32ToULEB128(i_ int32) []byte {
	i := uint32(i_)
	var result []byte
	c := 0
	for {
		currentByte := byte(i & 0b01111111)
		i >>= 7
		if i != 0 {
			currentByte |= 0b10000000
		}
		result = append(result, currentByte)
		c++

		if c >= 5 || i == 0 {
			return result
		}
	}
}

// From: https://github.com/aviate-labs/leb128/blob/v0.1.0/leb.go
func Int32ToLEB128(n int32) []byte {
	leb := make([]byte, 0)
	for {
		var (
			b    = byte(n & 0x7F)
			sign = byte(n & 0x40)
		)
		if n >>= 7; sign == 0 && n != 0 || n != -1 && (n != 0 || sign != 0) {
			b |= 0x80
		}
		leb = append(leb, b)
		if b&0x80 == 0 {
			break
		}
	}
	return leb
}

func LEB128ToInt32(bytes []byte) (int, error) {
	result := 0
	shift := 0

	for _, b := range bytes {
		result |= int(b&0x7f) << shift
		shift += 7
		if b&0x80 == 0 {
			return result, nil
		}
	}

	return 0, fmt.Errorf("invalid LEB128 encoding")
}

func NumBytesInLEB128(bytes []byte) (int, error) {
	result := 0
	shift := 0

	for i := 0; i < len(bytes); i++ {
		result |= int(bytes[i]&0x7f) << shift
		shift += 7
		if bytes[i]&0x80 == 0 {
			return i + 1, nil
		}
	}

	return 0, fmt.Errorf("invalid LEB128 encoding")
}
