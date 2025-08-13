package client

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"hash/crc32"
// 	"reflect"
// )

// func MergeAddCRC32Sum(args ...interface{}) ([]byte, uint32) {
// 	var buffer bytes.Buffer

// 	for _, arg := range args {
// 		// 使用反射来处理不同类型的参数
// 		if arg == nil {
// 			continue
// 		}
// 		v := reflect.ValueOf(arg)
// 		switch v.Kind() {
// 		case reflect.Struct:
// 			// 如果是结构体，遍历字段并写入缓冲区
// 			for i := 0; i < v.NumField(); i++ {
// 				field := v.Field(i)
// 				binary.Write(&buffer, binary.LittleEndian, field.Interface())
// 			}
// 		default:
// 			// 其他类型直接写入缓冲区
// 			binary.Write(&buffer, binary.LittleEndian, arg)
// 		}
// 	}

// 	// 计算 CRC32
// 	crc := crc32.ChecksumIEEE(buffer.Bytes())
// 	binary.Write(&buffer, binary.LittleEndian, crc)
// 	return buffer.Bytes(), crc
// }
