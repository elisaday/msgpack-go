package msgpack

import (
	"errors"
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

func PackUInt64(writer io.Writer, value uint64) (count int, err error) {
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

func PackUInt32(writer io.Writer, value uint32) (count int, err error) {
	return PackUInt64(writer, uint64(value))
}

func PackInt64(writer io.Writer, value int64) (count int, err error) {
	var n uint64
	n = uint64(value)
	if value >= 0 {
		switch {
		case value <= MAX_7BIT:
			return writer.Write(Bytes{MP_FIXNUM | uint8(n)})
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

func PackInt32(writer io.Writer, value int32) (count int, err error) {
	return PackInt64(writer, int64(value))
}

func PackBool(writer io.Writer, value bool) (count int, err error) {
	if value {
		return writer.Write(Bytes{MP_TRUE})
	} else {
		return writer.Write(Bytes{MP_FALSE})
	}
}

func PackFloat(writer io.Writer, value float32) (count int, err error) {
	var n uint32
	n = math.Float32bits(value)
	return writer.Write(Bytes{MP_FLOAT,
		uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)})
}

func PackDouble(writer io.Writer, value float64) (count int, err error) {
	var n uint64
	n = math.Float64bits(value)
	return writer.Write(Bytes{MP_DOUBLE,
		uint8(n >> 56), uint8(n >> 48), uint8(n >> 40), uint8(n >> 32),
		uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)})
}

func PackRawBuffer(writer io.Writer, value []uint8) (count int, err error) {
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

var ErrUnpackOverflow = errors.New("unpack overflow")

func unpackHeader(buf []byte, offset *uint32) (header uint8, err error) {
	off := *offset
	if int(off) >= len(buf) {
		return 0, ErrUnpackOverflow
	}

	(*offset)++
	return buf[off], nil
}

func UnpackUInt64(buf []byte, offset *uint32) (val uint64, err error) {
	header, err := unpackHeader(buf, offset)
	if err != nil {
		return 0, err
	}

	if header <= MAX_7BIT {
		return uint64(header), nil
	}

	off := *offset

	switch header {
	case MP_UINT8:
		(*offset)++
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return uint64(buf[off]), nil
	case MP_UINT16:
		(*offset) += 2
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return (uint64(buf[off]) << 8) | uint64(buf[off+1]), nil
	case MP_UINT32:
		(*offset) += 4
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return (uint64(buf[off]) << 24) | (uint64(buf[off+1]) << 16) |
			(uint64(buf[off+2]) << 8) | uint64(buf[off+3]), nil
	case MP_UINT64:
		(*offset) += 8
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return (uint64(buf[off]) << 56) | (uint64(buf[off+1]) << 48) |
			(uint64(buf[off+2]) << 40) | (uint64(buf[off+3]) << 32) |
			(uint64(buf[off+4]) << 24) | (uint64(buf[off+5]) << 16) |
			(uint64(buf[off+6]) << 8) | uint64(buf[off+7]), nil
	default:
		return 0, errors.New("invalid type header" + string(header))
	}
}

func UnpackInt64(buf []byte, offset *uint32) (val int64, err error) {
	header, err := unpackHeader(buf, offset)
	if err != nil {
		return 0, err
	}

	if header <= MAX_7BIT {
		return int64(header), nil
	}

	if (header & MP_NEGATIVE_FIXNUM) == MP_NEGATIVE_FIXNUM {
		return int64(header&0x1f) - 32, nil
	}

	off := *offset

	switch header {
	case MP_INT8:
		(*offset)++
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return int64(int8(buf[off])), nil
	case MP_INT16:
		(*offset) += 2
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return int64((int16(buf[off]) << 8) | int16(buf[off+1])), nil
	case MP_INT32:
		(*offset) += 4
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return int64((int32(buf[off]) << 24) | (int32(buf[off+1]) << 16) |
			(int32(buf[off+2]) << 8) | int32(buf[off+3])), nil
	case MP_INT64:
		(*offset) += 8
		if int(*offset) > len(buf) {
			return 0, ErrUnpackOverflow
		}
		return (int64(buf[off]) << 56) | (int64(buf[off+1]) << 48) |
			(int64(buf[off+2]) << 40) | (int64(buf[off+3]) << 32) |
			(int64(buf[off+4]) << 24) | (int64(buf[off+5]) << 16) |
			(int64(buf[off+6]) << 8) | int64(buf[off+7]), nil
	default:
		return 0, errors.New("invalid type header" + string(header))
	}
}

func UnpackUInt32(buf []byte, offset *uint32) (val uint32, err error) {
	return 0, nil
}

func UnpackInt32(buf []byte, offset *uint32) (val int32, err error) {
	return 0, nil
}

func UnpackBool(buf []byte, offset *uint32) (val bool, err error) {
	return true, nil
}

func UnpackFloat(buf []byte, offset *uint32) (val float32, err error) {
	return 0, nil
}

func UnpackDouble(buf []byte, offset *uint32) (val float64, err error) {
	return 0, nil
}

func UnpackRawBuffer(buf []byte, offset *uint32) (val []byte, err error) {
	return nil, nil
}
