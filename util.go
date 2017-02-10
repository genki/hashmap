package hashmap

import (
	"strconv"
	"encoding/gob"
	"bytes"
	"github.com/spaolacci/murmur3"
)

// intSizeBytes is the size in byte of an int or uint value.
const intSizeBytes = strconv.IntSize >> 3

// roundUpPower2 rounds a number to the next power of 2.
func roundUpPower2(i uint64) uint64 {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}

// log2 computes the binary logarithm of x, rounded up to the next integer.
func log2(i uint64) uint64 {
	var n, p uint64
	for p = 1; p < i; p += p {
		n++
	}
	return n
}

func Hash(key interface{}) uint64 {
	switch key.(type) {
	case []byte: return murmur3.Sum64(key.([]byte))
	default:
		if hashable, ok := key.(interface{Hash() uint64}); ok {
			return hashable.Hash()
		} else {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(key); err != nil {
				panic("Unsupported key type.")
			}
			return murmur3.Sum64(buf.Bytes())
		}
	}
}
