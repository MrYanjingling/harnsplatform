package biz

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

type JSONMap map[string]interface{}

// Value 实现 driver.Valuer 接口，将JSONMap转换为JSON字节
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口，将JSON字节解析为JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

type StringSlice []string

// Value 实现 driver.Valuer 接口
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "", nil
	}
	return strings.Join(s, ","), nil
}

// Scan 实现 sql.Scanner 接口
func (s *StringSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	*s = strings.Split(string(b), ",")
	return nil
}

// JSON 自定义类型，用于存储JSON数据
type JSON json.RawMessage

// Value 实现 driver.Valuer 接口，将JSON转换为字节切片
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// Scan 实现 sql.Scanner 接口，将数据库中的值解析为JSON
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("JSON类型扫描错误：不支持的值类型")
	}

	if len(bytes) == 0 {
		*j = nil
		return nil
	}

	result := make(json.RawMessage, len(bytes))
	copy(result, bytes)
	*j = JSON(result)
	return nil
}
