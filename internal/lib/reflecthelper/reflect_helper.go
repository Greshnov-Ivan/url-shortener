package reflecthelper

import "reflect"

func ComparePointers[T any](a, b *T) bool {
	// Если оба указателя равны nil, считаем их равными
	if a == nil && b == nil {
		return true
	}

	// Если один из указателей nil, считаем их не равными
	if a == nil || b == nil {
		return false
	}

	// Если оба указателя не nil, сравниваем их значения
	return reflect.DeepEqual(*a, *b)
}
