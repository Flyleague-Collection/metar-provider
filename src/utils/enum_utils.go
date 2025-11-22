// Package utils
package utils

// Enum 枚举类型结构体，用于表示具有值和标签的枚举项
// T 是可比较的类型
type Enum[T comparable, V any] struct {
	Value T `json:"value"`
	Data  V `json:"data"`
}

// NewEnum 创建一个新的枚举实例
func NewEnum[T comparable, V any](value T, data V) *Enum[T, V] {
	return &Enum[T, V]{Value: value, Data: data}
}

// ManagerInterface 枚举管理器接口，定义了枚举管理的基本操作
// T 是可比较的类型
type ManagerInterface[T comparable, V any] interface {
	IsValidEnum(value T) bool
	GetEnum(value T) *Enum[T, V]
	GetEnums() map[T]*Enum[T, V]
}

// EnumManager 枚举管理器结构体，用于管理一组枚举项
// T 是可比较的类型
type EnumManager[T comparable, V any] struct {
	enums map[T]*Enum[T, V]
}

func NewEnums[T comparable, V any](enums ...*Enum[T, V]) *EnumManager[T, V] {
	manager := &EnumManager[T, V]{
		enums: make(map[T]*Enum[T, V]),
	}
	for _, e := range enums {
		manager.enums[e.Value] = e
	}
	return manager
}

func (manager *EnumManager[T, V]) IsValidEnum(value T) bool {
	return manager.enums[value] != nil
}

func (manager *EnumManager[T, V]) GetEnum(value T) *Enum[T, V] {
	return manager.enums[value]
}

func (manager *EnumManager[T, V]) GetEnums() map[T]*Enum[T, V] {
	return manager.enums
}
