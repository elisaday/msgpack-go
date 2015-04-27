package msgpack

import (
	"io"
)

const (
	MP_INT8 = 0xd0
	MP_INT16 = 0xd1
	MP_INT32 = 0xd2
	MP_INT64 = 0xd3
	MP_UINT8 = 0xcc
	MP_UINT16 = 0xcd
	MP_UINT32 = 0xce
	MP_UINT64 = 0xcf
	MP_FIXNUM = 0x00
	MP_NEGATIVE_FIXNUM = 0xe0

	//! nil
	MP_NULL = 0xc0;

	//! boolean
	MP_FALSE = 0xc2;
	MP_TRUE = 0xc3;

	//! Floating point
	MP_FLOAT = 0xca;
	MP_DOUBLE = 0xcb;

	/*****************************************************
	* Variable length types
	*****************************************************/

	//! Raw bytes
	MP_RAW16 = 0xda;
	MP_RAW32 = 0xdb;
	MP_FIXRAW = 0xa0; //!< Last 5 bits is size

	/*****************************************************
	* Container types
	*****************************************************/

	//! Arrays
	MP_ARRAY16 = 0xdc
	MP_ARRAY32 = 0xdd
	MP_FIXARRAY = 0x90 //<! Lst 4 bits is size

	//! Maps
	MP_MAP16 = 0xde
	MP_MAP32 = 0xdf
	MP_FIXMAP = 0x80 //<! Last 4 bits is size

	//! Some helper bitmasks
	MAX_4BIT = 0xf
	MAX_5BIT = 0x1f
	MAX_7BIT = 0x7f
	MAX_8BIT = 0xff
	MAX_15BIT = 0x7fff
	MAX_16BIT = 0xffff
	MAX_31BIT = 0x7fffffff
	MAX_32BIT = 0xffffffff
)

type Bytes []byte

func PackUInt64(writer io.Writer, value uint64) (n int, err error) {
	return 0, nil
}
/*
func PackUint8(writer io.Writer, value uint8) (n int, err error) {
	if value >= REGULAR_UINT7_MAX {
		return writer.Write(Bytes{UINT8, value})
	}
	return writer.Write(Bytes{value})
}

func PackUint16(writer io.Writer, value uint16) (n int, err error) {
	if value >= REGULAR_UINT8_MAX {
		return writer.Write(Bytes{UINT16, byte(value >> 8), byte(value)})
	}
	return PackUint8(writer, uint8(value))
}

func PackUint32(writer io.Writer, value uint32) (n int, err error) {
	if value >= REGULAR_UINT16_MAX {
		return writer.Write(Bytes{UINT32, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
	}
	return PackUint16(writer, uint16(value))
}

func PackUint64(writer io.Writer, value uint64) (n int, err error) {
	if value >= REGULAR_UINT32_MAX {
		return writer.Write(Bytes{UINT64, byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
	}
	return PackUint32(writer, uint32(value))
}

func PackUint(writer io.Writer, value uint) (n int, err error) {
	switch unsafe.Sizeof(value) {
	case LEN_INT32:
		return PackUint32(writer, *(*uint32)(unsafe.Pointer(&value)))
	case LEN_INT64:
		return PackUint64(writer, *(*uint64)(unsafe.Pointer(&value)))
	}
	return 0, os.ErrNotExist // never get here
}

func PackInt8(writer io.Writer, value int8) (n int, err error) {
	if value < -SPECIAL_INT8 {
		return writer.Write(Bytes{INT8, byte(value)})
	}
	return writer.Write(Bytes{byte(value)})
}

func PackInt16(writer io.Writer, value int16) (n int, err error) {
	if value < -SPECIAL_INT16 || value >= SPECIAL_INT16 {
		return writer.Write(Bytes{INT16, byte(uint16(value) >> 8), byte(value)})
	}
	return PackInt8(writer, int8(value))
}

func PackInt32(writer io.Writer, value int32) (n int, err error) {
	if value < -SPECIAL_INT32 || value >= SPECIAL_INT32 {
		return writer.Write(Bytes{INT32, byte(uint32(value) >> 24), byte(uint32(value) >> 16), byte(uint32(value) >> 8), byte(value)})
	}
	return PackInt16(writer, int16(value))
}

func PackInt64(writer io.Writer, value int64) (n int, err error) {
	if value < -SPECIAL_INT64 || value >= SPECIAL_INT64 {
		return writer.Write(Bytes{INT64, byte(uint64(value) >> 56), byte(uint64(value) >> 48), byte(uint64(value) >> 40), byte(uint64(value) >> 32), byte(uint64(value) >> 24), byte(uint64(value) >> 16), byte(uint64(value) >> 8), byte(value)})
	}
	return PackInt32(writer, int32(value))
}

func PackInt(writer io.Writer, value int) (n int, err error) {
	switch unsafe.Sizeof(value) {
	case LEN_INT32:
		return PackInt32(writer, *(*int32)(unsafe.Pointer(&value)))
	case LEN_INT64:
		return PackInt64(writer, *(*int64)(unsafe.Pointer(&value)))
	}
	return 0, os.ErrNotExist // never get here
}

func PackNil(writer io.Writer) (n int, err error) {
	return writer.Write(Bytes{NIL})
}

func PackBool(writer io.Writer, value bool) (n int, err error) {
	if value {
		return writer.Write(Bytes{TRUE})
	}
	return writer.Write(Bytes{FALSE})
}

func PackFloat32(writer io.Writer, value float32) (n int, err error) {
	return PackUint32(writer, *(*uint32)(unsafe.Pointer(&value)))
}

func PackFloat64(writer io.Writer, value float64) (n int, err error) {
	return PackUint64(writer, *(*uint64)(unsafe.Pointer(&value)))
}

func PackBytes(writer io.Writer, value []byte) (n int, err error) {
	length := len(value)
	if length < MAXFIXRAW {
		n1, err := writer.Write(Bytes{FIXRAW | uint8(length)})
		if err != nil {
			return n1, err
		}
		n2, err := writer.Write(value)
		return n1 + n2, err
	} else if length < MAX16BIT {
		n1, err := writer.Write(Bytes{RAW16, byte(length >> 16), byte(length)})
		if err != nil {
			return n1, err
		}
		n2, err := writer.Write(value)
		return n1 + n2, err
	}
	n1, err := writer.Write(Bytes{RAW32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write(value)
	return n1 + n2, err
}

func PackUint16Array(writer io.Writer, value []uint16) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint16(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint16(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint16(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackUint32Array(writer io.Writer, value []uint32) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackUint64Array(writer io.Writer, value []uint64) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackUint64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackUintArray(writer io.Writer, value []uint) (n int, err error) {
	switch unsafe.Sizeof(0) {
	case 4:
		return PackUint32Array(writer, *(*[]uint32)(unsafe.Pointer(&value)))
	case 8:
		return PackUint64Array(writer, *(*[]uint64)(unsafe.Pointer(&value)))
	}
	return 0, os.ErrNotExist // never get here
}

func PackInt8Array(writer io.Writer, value []int8) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt8(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt8(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt8(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackInt16Array(writer io.Writer, value []int16) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt16(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt16(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt16(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackInt32Array(writer io.Writer, value []int32) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackInt64Array(writer io.Writer, value []int64) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackInt64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackIntArray(writer io.Writer, value []int) (n int, err error) {
	switch unsafe.Sizeof(0) {
	case LEN_INT32:
		return PackInt32Array(writer, *(*[]int32)(unsafe.Pointer(&value)))
	case LEN_INT64:
		return PackInt64Array(writer, *(*[]int64)(unsafe.Pointer(&value)))
	}
	return 0, os.ErrNotExist // never get here
}

func PackFloat32Array(writer io.Writer, value []float32) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackFloat32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackFloat32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackFloat32(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackFloat64Array(writer io.Writer, value []float64) (n int, err error) {
	length := len(value)
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackFloat64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackFloat64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, i := range value {
			_n, err := PackFloat64(writer, i)
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackArray(writer io.Writer, value reflect.Value) (n int, err error) {
	{
		elemType := value.Type().Elem()
		if (elemType.Kind() == reflect.Uint || elemType.Kind() == reflect.Uint8 || elemType.Kind() == reflect.Uint16 || elemType.Kind() == reflect.Uint32 || elemType.Kind() == reflect.Uint64 || elemType.Kind() == reflect.Uintptr) &&
			elemType.Kind() == reflect.Uint8 {
			return PackBytes(writer, value.Interface().([]byte))
		}
	}

	length := value.Len()
	if length < MAXFIXARRAY {
		n, err := writer.Write(Bytes{FIXARRAY | byte(length)})
		if err != nil {
			return n, err
		}
		for i := 0; i < length; i++ {
			_n, err := PackValue(writer, value.Index(i))
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for i := 0; i < length; i++ {
			_n, err := PackValue(writer, value.Index(i))
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{ARRAY32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for i := 0; i < length; i++ {
			_n, err := PackValue(writer, value.Index(i))
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackMap(writer io.Writer, value reflect.Value) (n int, err error) {
	keys := value.MapKeys()
	length := len(keys)
	if length < MAXFIXMAP {
		n, err := writer.Write(Bytes{FIXMAP | byte(length)})
		if err != nil {
			return n, err
		}
		for _, k := range keys {
			_n, err := PackValue(writer, k)
			if err != nil {
				return n, err
			}
			n += _n
			_n, err = PackValue(writer, value.MapIndex(k))
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else if length < MAX16BIT {
		n, err := writer.Write(Bytes{MAP16, byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, k := range keys {
			_n, err := PackValue(writer, k)
			if err != nil {
				return n, err
			}
			n += _n
			_n, err = PackValue(writer, value.MapIndex(k))
			if err != nil {
				return n, err
			}
			n += _n
		}
	} else {
		n, err := writer.Write(Bytes{MAP32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			return n, err
		}
		for _, k := range keys {
			_n, err := PackValue(writer, k)
			if err != nil {
				return n, err
			}
			n += _n
			_n, err = PackValue(writer, value.MapIndex(k))
			if err != nil {
				return n, err
			}
			n += _n
		}
	}
	return n, nil
}

func PackValue(writer io.Writer, value reflect.Value) (n int, err error) {
	if !value.IsValid() || value.Type() == nil {
		return PackNil(writer)
	}
	switch _value := value; _value.Kind() {
	case reflect.Bool:
		return PackBool(writer, _value.Bool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return PackUint64(writer, _value.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return PackInt64(writer, _value.Int())
	case reflect.Float32, reflect.Float64:
		return PackFloat64(writer, _value.Float())
	case reflect.Array:
		return PackArray(writer, _value)
	case reflect.Slice:
		return PackArray(writer, _value)
	case reflect.Map:
		return PackMap(writer, _value)
	case reflect.String:
		return PackBytes(writer, []byte(_value.String()))
	case reflect.Interface:
		__value := reflect.ValueOf(_value.Interface())

		if __value.Kind() != reflect.Interface {
			return PackValue(writer, __value)
		}
	}
	panic("unsupported type: " + value.Type().String())
}

func Pack(writer io.Writer, value interface{}) (n int, err error) {
	if value == nil {
		return PackNil(writer)
	}
	switch _value := value.(type) {
	case bool:
		return PackBool(writer, _value)
	case uint8:
		return PackUint8(writer, _value)
	case uint16:
		return PackUint16(writer, _value)
	case uint32:
		return PackUint32(writer, _value)
	case uint64:
		return PackUint64(writer, _value)
	case uint:
		return PackUint(writer, _value)
	case int8:
		return PackInt8(writer, _value)
	case int16:
		return PackInt16(writer, _value)
	case int32:
		return PackInt32(writer, _value)
	case int64:
		return PackInt64(writer, _value)
	case int:
		return PackInt(writer, _value)
	case float32:
		return PackFloat32(writer, _value)
	case float64:
		return PackFloat64(writer, _value)
	case []byte:
		return PackBytes(writer, _value)
	case []uint16:
		return PackUint16Array(writer, _value)
	case []uint32:
		return PackUint32Array(writer, _value)
	case []uint64:
		return PackUint64Array(writer, _value)
	case []uint:
		return PackUintArray(writer, _value)
	case []int8:
		return PackInt8Array(writer, _value)
	case []int16:
		return PackInt16Array(writer, _value)
	case []int32:
		return PackInt32Array(writer, _value)
	case []int64:
		return PackInt64Array(writer, _value)
	case []int:
		return PackIntArray(writer, _value)
	case []float32:
		return PackFloat32Array(writer, _value)
	case []float64:
		return PackFloat64Array(writer, _value)
	case string:
		return PackBytes(writer, Bytes(_value))
	default:
		return PackValue(writer, reflect.ValueOf(value))
	}
	return 0, nil // never get here
}
*/