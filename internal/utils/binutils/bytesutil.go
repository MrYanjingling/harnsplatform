package binutils

import "math"

func ParseUint64(b []byte) uint64 {
	return (uint64(b[0]) << 56) |
		(uint64(b[1]) << 48) |
		(uint64(b[2]) << 40) |
		(uint64(b[3]) << 32) |
		(uint64(b[4]) << 24) |
		(uint64(b[5]) << 16) |
		(uint64(b[6]) << 8) |
		uint64(b[7])
}

// ABCD EFGH
func ParseUint64BigEndian(b []byte) uint64 {
	return (uint64(b[0]) << 56) |
		(uint64(b[1]) << 48) |
		(uint64(b[2]) << 40) |
		(uint64(b[3]) << 32) |
		(uint64(b[4]) << 24) |
		(uint64(b[5]) << 16) |
		(uint64(b[6]) << 8) |
		uint64(b[7])
}

// BADC FEHG
func ParseUint64BigEndianByteSwap(b []byte) uint64 {
	return (uint64(b[1]) << 56) |
		(uint64(b[0]) << 48) |
		(uint64(b[3]) << 40) |
		(uint64(b[2]) << 32) |
		(uint64(b[5]) << 24) |
		(uint64(b[4]) << 16) |
		(uint64(b[7]) << 8) |
		uint64(b[6])
}

// ParseUint64LittleEndian 解析
// HGFE DCBA
func ParseUint64LittleEndian(b []byte) uint64 {
	return (uint64(b[7]) << 56) |
		(uint64(b[6]) << 48) |
		(uint64(b[5]) << 40) |
		(uint64(b[4]) << 32) |
		(uint64(b[3]) << 24) |
		(uint64(b[2]) << 16) |
		(uint64(b[1]) << 8) |
		uint64(b[0])
}

// GHEF CDAB
func ParseUint64LittleEndianByteSwap(b []byte) uint64 {
	return (uint64(b[6]) << 56) |
		(uint64(b[7]) << 48) |
		(uint64(b[4]) << 40) |
		(uint64(b[5]) << 32) |
		(uint64(b[2]) << 24) |
		(uint64(b[3]) << 16) |
		(uint64(b[0]) << 8) |
		uint64(b[1])
}

// ParseUint32 解析
func ParseUint32(buf []byte) uint32 {
	return uint32(buf[0])<<24 +
		uint32(buf[1])<<16 +
		uint32(buf[2])<<8 +
		uint32(buf[3])
}

// ABCD
func ParseUint32BigEndian(buf []byte) uint32 {
	return uint32(buf[0])<<24 +
		uint32(buf[1])<<16 +
		uint32(buf[2])<<8 +
		uint32(buf[3])
}

// BADC
func ParseUint32BigEndianByteSwap(buf []byte) uint32 {
	return uint32(buf[1])<<24 +
		uint32(buf[0])<<16 +
		uint32(buf[3])<<8 +
		uint32(buf[2])
}

// ParseUint32LittleEndian 解析
// DCBA
func ParseUint32LittleEndian(buf []byte) uint32 {
	return uint32(buf[3])<<24 +
		uint32(buf[2])<<16 +
		uint32(buf[1])<<8 +
		uint32(buf[0])
}

// CDAB
func ParseUint32LittleEndianByteSwap(buf []byte) uint32 {
	return uint32(buf[2])<<24 +
		uint32(buf[3])<<16 +
		uint32(buf[0])<<8 +
		uint32(buf[1])
}

// ParseUint16 解析
func ParseUint16(buf []byte) uint16 {
	return uint16(buf[0])<<8 + uint16(buf[1])
}

func ParseUint16BigEndian(buf []byte) uint16 {
	return uint16(buf[0])<<8 + uint16(buf[1])
}

// ParseUint16LittleEndian 解析
func ParseUint16LittleEndian(buf []byte) uint16 {
	return uint16(buf[1])<<8 + uint16(buf[0])
}

// ParseFloat32 解析
func ParseFloat32(buf []byte) float32 {
	val := ParseUint32(buf)
	return math.Float32frombits(val)
}

func ParseFloat32BigEndian(buf []byte) float32 {
	val := ParseUint32BigEndian(buf)
	return math.Float32frombits(val)
}

func ParseFloat32BigEndianByteSwap(buf []byte) float32 {
	val := ParseUint32BigEndianByteSwap(buf)
	return math.Float32frombits(val)
}

// ParseFloat32LittleEndian 解析
func ParseFloat32LittleEndian(buf []byte) float32 {
	val := ParseUint32LittleEndian(buf)
	return math.Float32frombits(val)
}

func ParseFloat32LittleEndianByteSwap(buf []byte) float32 {
	val := ParseUint32LittleEndianByteSwap(buf)
	return math.Float32frombits(val)
}

func ParseFloat64BigEndian(buf []byte) float64 {
	val := ParseUint64BigEndian(buf)
	return math.Float64frombits(val)
}

func ParseFloat64BigEndianByteSwap(buf []byte) float64 {
	val := ParseUint64BigEndianByteSwap(buf)
	return math.Float64frombits(val)
}

// ParseFloat64LittleEndian 解析
func ParseFloat64LittleEndian(buf []byte) float64 {
	val := ParseUint64LittleEndian(buf)
	return math.Float64frombits(val)
}

func ParseFloat64LittleEndianByteSwap(buf []byte) float64 {
	val := ParseUint64LittleEndianByteSwap(buf)
	return math.Float64frombits(val)
}

// Uint32ToBytesBigEndian 编码
func Uint32ToBytesBigEndian(value uint32) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(value >> 24)
	buf[1] = byte(value >> 16)
	buf[2] = byte(value >> 8)
	buf[3] = byte(value)
	return buf
}

func Uint32ToBytesBigEndianByteSwap(value uint32) []byte {
	buf := make([]byte, 4)
	buf[1] = byte(value >> 24)
	buf[0] = byte(value >> 16)
	buf[3] = byte(value >> 8)
	buf[2] = byte(value)
	return buf
}

// Uint32ToBytesLittleEndian 编码
func Uint32ToBytesLittleEndian(value uint32) []byte {
	buf := make([]byte, 4)
	buf[3] = byte(value >> 24)
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
	return buf
}

func Uint32ToBytesLittleEndianByteSwap(value uint32) []byte {
	buf := make([]byte, 4)
	buf[2] = byte(value >> 24)
	buf[3] = byte(value >> 16)
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
	return buf
}

func Uint64ToBytesBigEndian(value uint64) []byte {
	buf := make([]byte, 8)
	buf[0] = byte(value >> 56)
	buf[1] = byte(value >> 48)
	buf[2] = byte(value >> 40)
	buf[3] = byte(value >> 32)
	buf[4] = byte(value >> 24)
	buf[5] = byte(value >> 16)
	buf[6] = byte(value >> 8)
	buf[7] = byte(value)
	return buf
}

func Uint64ToBytesBigEndianByteSwap(value uint64) []byte {
	buf := make([]byte, 8)
	buf[1] = byte(value >> 56)
	buf[0] = byte(value >> 48)
	buf[3] = byte(value >> 40)
	buf[2] = byte(value >> 32)
	buf[5] = byte(value >> 24)
	buf[4] = byte(value >> 16)
	buf[7] = byte(value >> 8)
	buf[6] = byte(value)
	return buf
}

// Uint64ToBytesLittleEndian 编码
func Uint64ToBytesLittleEndian(value uint64) []byte {
	buf := make([]byte, 8)
	buf[7] = byte(value >> 56)
	buf[6] = byte(value >> 48)
	buf[5] = byte(value >> 40)
	buf[4] = byte(value >> 32)
	buf[3] = byte(value >> 24)
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
	return buf
}

func Uint64ToBytesLittleEndianByteSwap(value uint64) []byte {
	buf := make([]byte, 8)
	buf[6] = byte(value >> 56)
	buf[7] = byte(value >> 48)
	buf[4] = byte(value >> 40)
	buf[5] = byte(value >> 32)
	buf[2] = byte(value >> 24)
	buf[3] = byte(value >> 16)
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
	return buf
}

func Uint16ToBytesBigEndian(value uint16) []byte {
	buf := make([]byte, 2)
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
	return buf
}

// Uint16ToBytesLittleEndian 编码
func Uint16ToBytesLittleEndian(value uint16) []byte {
	buf := make([]byte, 2)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
	return buf
}

func Float32ToBytesBigEndian(value float32) []byte {
	val := math.Float32bits(value)
	return Uint32ToBytesBigEndian(val)
}

func Float32ToBytesBigEndianByteSwap(value float32) []byte {
	val := math.Float32bits(value)
	return Uint32ToBytesBigEndianByteSwap(val)
}

// Float32ToBytesLittleEndian 编码
func Float32ToBytesLittleEndian(value float32) []byte {
	val := math.Float32bits(value)
	return Uint32ToBytesLittleEndian(val)
}

func Float32ToBytesLittleEndianByteSwap(value float32) []byte {
	val := math.Float32bits(value)
	return Uint32ToBytesLittleEndianByteSwap(val)
}

func Float64ToBytesBigEndian(value float64) []byte {
	val := math.Float64bits(value)
	return Uint64ToBytesBigEndian(val)
}

func Float64ToBytesBigEndianByteSwap(value float64) []byte {
	val := math.Float64bits(value)
	return Uint64ToBytesBigEndianByteSwap(val)
}

// Float64ToBytesLittleEndian 编码
func Float64ToBytesLittleEndian(value float64) []byte {
	val := math.Float64bits(value)
	return Uint64ToBytesLittleEndian(val)
}

func Float64ToBytesLittleEndianByteSwap(value float64) []byte {
	val := math.Float64bits(value)
	return Uint64ToBytesLittleEndianByteSwap(val)
}

// WriteUint64 编码
// ABCD EFGH
func WriteUint64BigEndian(buf []byte, value uint64) {
	buf[0] = byte(value >> 56)
	buf[1] = byte(value >> 48)
	buf[2] = byte(value >> 40)
	buf[3] = byte(value >> 32)
	buf[4] = byte(value >> 24)
	buf[5] = byte(value >> 16)
	buf[6] = byte(value >> 8)
	buf[7] = byte(value)
}

// BADC FEHG
func WriteUint64BigEndianByteSwap(buf []byte, value uint64) {
	buf[1] = byte(value >> 56)
	buf[0] = byte(value >> 48)
	buf[3] = byte(value >> 40)
	buf[2] = byte(value >> 32)
	buf[5] = byte(value >> 24)
	buf[4] = byte(value >> 16)
	buf[7] = byte(value >> 8)
	buf[6] = byte(value)
}

// WriteUint64LittleEndian 编码
// HGFE DCBA
func WriteUint64LittleEndian(buf []byte, value uint64) {
	buf[7] = byte(value >> 56)
	buf[6] = byte(value >> 48)
	buf[5] = byte(value >> 40)
	buf[4] = byte(value >> 32)
	buf[3] = byte(value >> 24)
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
}

// GHEF CDAB
func WriteUint64LittleEndianByteSwap(buf []byte, value uint64) {
	buf[6] = byte(value >> 56)
	buf[7] = byte(value >> 48)
	buf[4] = byte(value >> 40)
	buf[5] = byte(value >> 32)
	buf[2] = byte(value >> 24)
	buf[3] = byte(value >> 16)
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
}

// ABCD
func WriteUint32BigEndian(buf []byte, value uint32) {
	buf[0] = byte(value >> 24)
	buf[1] = byte(value >> 16)
	buf[2] = byte(value >> 8)
	buf[3] = byte(value)
}

// BADC
func WriteUint32BigEndianByteSwap(buf []byte, value uint32) {
	buf[1] = byte(value >> 24)
	buf[0] = byte(value >> 16)
	buf[3] = byte(value >> 8)
	buf[2] = byte(value)
}

// WriteUint32LittleEndian 编码
// DCBA
func WriteUint32LittleEndian(buf []byte, value uint32) {
	buf[3] = byte(value >> 24)
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
}

// CDAB
func WriteUint32LittleEndianByteSwap(buf []byte, value uint32) {
	buf[2] = byte(value >> 24)
	buf[3] = byte(value >> 16)
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
}

// WriteUint24 编码
func WriteUint24BigEndian(buf []byte, value uint32) {
	buf[0] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[2] = byte(value)
}

// WriteUint24LittleEndian 编码
func WriteUint24LittleEndian(buf []byte, value uint32) {
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
}

// WriteUint16 编码
// func WriteUint16(buf []byte, value uint16) {
// 	buf[0] = byte(value >> 8)
// 	buf[1] = byte(value)
// }

func WriteUint16BigEndian(buf []byte, value uint16) {
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
}

// WriteUint16LittleEndian 编码
func WriteUint16LittleEndian(buf []byte, value uint16) {
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
}

// WriteFloat32BigEndian 编码
func WriteFloat32BigEndian(buf []byte, value float32) {
	val := math.Float32bits(value)
	WriteUint32BigEndian(buf, val)
}

func WriteFloat32BigEndianByteSwap(buf []byte, value float32) {
	val := math.Float32bits(value)
	WriteUint32BigEndianByteSwap(buf, val)
}

// WriteFloat32LittleEndian 编码
func WriteFloat32LittleEndian(buf []byte, value float32) {
	val := math.Float32bits(value)
	WriteUint32LittleEndian(buf, val)
}

func WriteFloat32LittleEndianByteSwap(buf []byte, value float32) {
	val := math.Float32bits(value)
	WriteUint32LittleEndianByteSwap(buf, val)
}

// WriteFloat64BigEndian 编码
func WriteFloat64BigEndian(buf []byte, value float64) {
	val := math.Float64bits(value)
	WriteUint64BigEndian(buf, val)
}

func WriteFloat64BigEndianByteSwap(buf []byte, value float64) {
	val := math.Float64bits(value)
	WriteUint64BigEndianByteSwap(buf, val)
}

// WriteFloat64LittleEndian 编码
func WriteFloat64LittleEndian(buf []byte, value float64) {
	val := math.Float64bits(value)
	WriteUint64LittleEndian(buf, val)
}

func WriteFloat64LittleEndianByteSwap(buf []byte, value float64) {
	val := math.Float64bits(value)
	WriteUint64LittleEndianByteSwap(buf, val)
}

// BoolToAscii 编码
func BoolToAscii(buf []byte) []byte {
	length := len(buf)
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		if buf[i] == 0 {
			ret[i] = '0'
		} else {
			ret[i] = '1'
		}
	}
	return ret
}

// AsciiToBool 编码
func AsciiToBool(buf []byte) []byte {
	length := len(buf)
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		if buf[i] == '0' {
			ret[i] = 0
		} else {
			ret[i] = 1
		}
	}
	return ret
}

// Dup 复制
func Dup(buf []byte) []byte {
	b := make([]byte, len(buf))
	copy(b, buf)
	return b
}

// BoolToByte 编码
func BoolToByte(buf []bool) []byte {
	r := make([]byte, len(buf))
	for i, v := range buf {
		if v {
			r[i] = 1
		}
	}
	return r
}

// ByteToBool 编码
func ByteToBool(buf []byte) []bool {
	r := make([]bool, len(buf))
	for i, v := range buf {
		if v > 0 {
			r[i] = true
		}
	}
	return r
}

// ShrinkBool 压缩布尔类型
func ShrinkBool(buf []byte) []byte {
	length := len(buf)
	// length = length % 8 == 0 ? length / 8 : length / 8 + 1;
	ln := length >> 3    // length/8
	if length&0x07 > 0 { // length%8
		ln++
	}

	b := make([]byte, ln)

	for i := 0; i < length; i++ {
		if buf[i] > 0 {
			// b[i/8] += 1 << (i % 8)
			b[i>>3] += 1 << (i & 0x07)
		}
	}

	return b
}

// ExpandBool 展开布尔类型
func ExpandBool(buf []byte, count int) []byte {
	if count > len(buf) {
		count = len(buf)
	}
	expandLength := count << 3
	b := make([]byte, expandLength)
	for i := 0; i < expandLength; i++ {
		if buf[i>>3]&(1<<(i&0x07)) > 0 {
			b[i] = 1
		}
	}
	return b
}
