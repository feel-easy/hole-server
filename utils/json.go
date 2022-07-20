package utils

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	jsoniter.RegisterTypeEncoderFunc("[]uint8", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		t := *((*[]byte)(ptr))
		stream.WriteString(string(t))
	}, nil)
	jsoniter.RegisterTypeDecoderFunc("[]uint8", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		str := iter.ReadString()
		*((*[]byte)(ptr)) = []byte(str)
	})
}

func MarshalJson(v interface{}) []byte {
	data, _ := jsoniter.Marshal(v)
	return data
}

func UnmarshalJson(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
