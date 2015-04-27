package msgpack

import (
	"io"
	"math"
)

const (
	MP_INT8            = 0xd0
	MP_INT16           = 0xd1
	MP_INT32           = 0xd2
	MP_INT64           = 0xd3
	MP_UINT8           = 0xcc
	MP_UINT16          = 0xcd
	MP_UINT32          = 0xce
	MP_UINT64          = 0xcf
	MP_FIXNUM          = 0x00
	MP_NEGATIVE_FIXNUM = 0xe0

	//! nil
	MP_NULL = 0xc0

	//! boolean
	MP_FALSE = 0xc2
	MP_TRUE  = 0xc3

	//! Floating point
	MP_FLOAT  = 0xca
	MP_DOUBLE = 0xcb

	/*****************************************************
	* Variable length types
	*****************************************************/

	//! Raw bytes
	MP_RAW16  = 0xda
	MP_RAW32  = 0xdb
	MP_FIXRAW = 0xa0 //!< Last 5 bits is size

	/*****************************************************
	* Container types
	*****************************************************/

	//! Arrays
	MP_ARRAY16  = 0xdc
	MP_ARRAY32  = 0xdd
	MP_FIXARRAY = 0x90 //<! Lst 4 bits is size

	//! Maps
	MP_MAP16  = 0xde
	MP_MAP32  = 0xdf
	MP_FIXMAP = 0x80 //<! Last 4 bits is size

	//! Some helper bitmasks
	MAX_4BIT  = 0xf
	MAX_5BIT  = 0x1f
	MAX_7BIT  = 0x7f
	MAX_8BIT  = 0xff
	MAX_15BIT = 0x7fff
	MAX_16BIT = 0xffff
	MAX_31BIT = 0x7fffffff
	MAX_32BIT = 0xffffffff
)

type Bytes []uint8

func PackUInt64(writer io.Writer, value uint64) (ret int, err error) {
	switch {
	case value <= MAX_7BIT:
		return writer.Write(Bytes{MP_FIXNUM | uint8(value)})
	case value <= MAX_8BIT:
		return writer.Write(Bytes{MP_UINT8, uint8(value)})
	case value <= MAX_16BIT:
		return writer.Write(Bytes{MP_UINT16, uint8(value >> 8), uint8(value)})
	case value <= MAX_32BIT:
		return writer.Write(Bytes{MP_UINT32,
			uint8(value >> 24), uint8(value >> 16), uint8(value >> 8), uint8(value)})
	default:
		return writer.Write(Bytes{MP_UINT64,
			uint8(value >> 56), uint8(value >> 48), uint8(value >> 40), uint8(value >> 32),
			uint8(value >> 24), uint8(value >> 16), uint8(value >> 8), uint8(value)})
	}
}

func PackInt64(writer io.Writer, value int64) (ret int, err error) {
	var n uint64
	n = n
	if value >= 0 {
		switch {
		case value <= MAX_7BIT:
			return writer.Write(Bytes{uint8(n)})
		case value <= MAX_15BIT:
			return writer.Write(Bytes{MP_INT16, uint8(uint64(n) >> 8), uint8(n)})
		case value <= MAX_31BIT:
			return writer.Write(Bytes{MP_INT32,
				uint8(n >> 24), uint8(n >> 16),
				uint8(n >> 8), uint8(n)})
		default:
			return writer.Write(Bytes{MP_INT64,
				uint8(n >> 56), uint8(n >> 48), uint8(n >> 40), uint8(n >> 32),
				uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)})
		}
	} else {
		switch {
		case value >= -(MAX_5BIT + 1):
			return writer.Write(Bytes{MP_NEGATIVE_FIXNUM | uint8(value)})
		case value >= -(int64(MAX_7BIT) + 1):
			return writer.Write(Bytes{MP_INT8, uint8(value)})
		case value >= -(int64(MAX_15BIT) + 1):
			return writer.Write(Bytes{MP_INT16, uint8(n >> 8), uint8(n)})
		case value >= -(int64(MAX_31BIT) + 1):
			return writer.Write(Bytes{MP_INT32,
				uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(value)})
		default:
			return writer.Write(Bytes{MP_INT64,
				uint8(n >> 56), uint8(n >> 48), uint8(n >> 40), uint8(n >> 32),
				uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)})
		}
	}
}

func PackBool(writer io.Writer, value bool) (ret int, err error) {
	if value {
		return writer.Write(Bytes{MP_TRUE})
	} else {
		return writer.Write(Bytes{MP_FALSE})
	}
}

func PackFloat(writer io.Writer, value float32) (ret int, err error) {
	var n uint32
	n = math.Float32bits(value)
	return writer.Write(Bytes{MP_FLOAT,
		uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)})
}

func PackDouble(writer io.Writer, value float64) (ret int, err error) {
	var n uint64
	n = math.Float64bits(value)
	return writer.Write(Bytes{MP_DOUBLE,
		uint8(n >> 56), uint8(n >> 48), uint8(n >> 40), uint8(n >> 32),
		uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)})
}

func PackRawBuffer(writer io.Writer, value []uint8) (ret int, err error) {
	var length uint64
	length = uint64(len(value))
	var n int
	var e error

	if length <= MAX_5BIT {
		n, e = writer.Write(Bytes{MP_FIXRAW | uint8(length)})
	} else if length <= MAX_16BIT {
		n, e = writer.Write(Bytes{MP_RAW16, uint8(uint16(length) >> 8), uint8(length)})
	} else {
		n, e = writer.Write(Bytes{MP_RAW32,
			uint8(length >> 24), uint8(length >> 16), uint8(length >> 8), uint8(length)})
	}

	if e != nil {
		return n, e
	}

	return writer.Write(value)
}
