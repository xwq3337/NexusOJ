package jsonx

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrInvalidInputType = errors.New("input must be a non-nil pointer")
	ErrNonPointerOutput = errors.New("output must be a non-nil pointer")
)

// Marshal 高效JSON序列化（无缩进）
func Marshal(v any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // 禁用HTML转义提高性能

	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	// 移除末尾换行符
	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	return data, nil
}

// MarshalString 序列化为JSON字符串
func MarshalString(v any) (string, error) {
	data, err := Marshal(v)
	return string(data), err
}

// MarshalIndent 带缩进的JSON序列化
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(data, []byte("\n")), nil
}

// Unmarshal 高性能JSON反序列化
func Unmarshal(data []byte, v any) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr || reflect.ValueOf(v).IsNil() {
		return ErrNonPointerOutput
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber() // 保留数字精度

	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

// UnmarshalString 从字符串反序列化
func UnmarshalString(s string, v any) error {
	return Unmarshal([]byte(s), v)
}

// ToMap 结构体转map[string]any
func ToMap(v any) (map[string]any, error) {
	data, err := Marshal(v)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	if err := Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// FromMap map转结构体
func FromMap(m map[string]any, v any) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr || reflect.ValueOf(v).IsNil() {
		return ErrNonPointerOutput
	}

	data, err := Marshal(m)
	if err != nil {
		return err
	}
	return Unmarshal(data, v)
}

// Transform 类型转换（通过JSON中转）
func Transform(input, output any) error {
	data, err := Marshal(input)
	if err != nil {
		return err
	}
	return Unmarshal(data, output)
}

/** Example usage:
// 结构体转JSON字符串
s, err := jsonx.MarshalString(user)
if err != nil {  ... handle error... }

// JSON字符串转结构体
err := jsonx.UnmarshalString(jsonStr, &user)

// 类型转换
var dto UserDTO
err := jsonx.Transform(rawData, &dto)

// Map转结构体
m := map[string]any{"name": "John", "age": 30}
err := jsonx.FromMap(m, &user)

// 结构体转Map
result, err := jsonx.ToMap(user)
*/
