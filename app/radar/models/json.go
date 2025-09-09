package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// 泛型 JSON 类型
type JSON[T any] struct {
	Data []T
}

// 实现 GormDataType 接口
func (j *JSON[T]) GormDataType() string {
	return "json"
}

// 实现 sql.Scanner 接口（读取数据库）
func (j *JSON[T]) Scan(value interface{}) error {
	if value == nil {
		j.Data = nil
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into JSON", value)
	}

	return json.Unmarshal(b, &j.Data)
}

// 实现 driver.Valuer 接口（写入数据库）
func (j JSON[T]) Value() (driver.Value, error) {
	if j.Data == nil {
		return []byte("[]"), nil // 空数组存储为 [] 而不是 NULL
	}
	return json.Marshal(j.Data)
}

// 实现 json.Marshaler 接口
func (j JSON[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}

// 实现 json.Unmarshaler 接口
func (j *JSON[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &j.Data)
}

// Getter / Setter
func (j *JSON[T]) Get() []T {
	return j.Data
}

func (j *JSON[T]) Set(data []T) {
	j.Data = data
}
