package hashmap

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"math/rand"
)

func hash[K comparable](key K) int {
	h := fnv.New32a()
	var buf [8]byte
	switch v := any(key).(type) {
	case string:
		h.Write([]byte(v))
	case int:
		binary.LittleEndian.PutUint64(buf[:], uint64(v))
		h.Write(buf[:8])
	case int8:
		buf[0] = byte(v)
		h.Write(buf[:1])
	case int16:
		binary.LittleEndian.PutUint16(buf[:], uint16(v))
		h.Write(buf[:2])
	case int32:
		binary.LittleEndian.PutUint32(buf[:], uint32(v))
		h.Write(buf[:4])
	case int64:
		binary.LittleEndian.PutUint64(buf[:], uint64(v))
		h.Write(buf[:8])
	case uint8:
		buf[0] = v
		h.Write(buf[:1])
	case uint16:
		binary.LittleEndian.PutUint16(buf[:], v)
		h.Write(buf[:2])
	case uint32:
		binary.LittleEndian.PutUint32(buf[:], v)
		h.Write(buf[:4])
	case uint64:
		binary.LittleEndian.PutUint64(buf[:], v)
		h.Write(buf[:8])
	case float32:
		binary.LittleEndian.PutUint32(buf[:], math.Float32bits(v))
		h.Write(buf[:4])
	case float64:
		binary.LittleEndian.PutUint64(buf[:], math.Float64bits(v))
		h.Write(buf[:8])
	case bool:
		if v {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
		h.Write(buf[:1])
	default:
		// For unsupported types, consider alternative handling
	}
	return int(h.Sum32())
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
