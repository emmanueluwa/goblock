package types

import (
	"fmt"
	"reflect"
)

type List[T any] struct {
	Data []T
}

func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

func (list *List[T]) Get(index int) T {
	if index > len(list.Data)-1 {
		err := fmt.Sprintf("the given index (%d) is higher than the length (%d)", index, len(list.Data))
		panic(err)
	}
	return list.Data[index]
}

func (list *List[T]) Insert(v T) {
	list.Data = append(list.Data, v)
}

func (list *List[T]) Clear() {
	list.Data = []T{}
}

// GetIndex, return the index v. If v deos not exist in list
// return -1.

func (list *List[T]) GetIndex(v T) int {
	for i := 0; i < list.Len(); i++ {
		if reflect.DeepEqual(v, list.Data[i]) {
			return i
		}
	}
	return -1
}

func (list *List[T]) Remove(v T) {
	index := list.GetIndex(v)
	if index == -1 {
		return
	}
	list.Pop(index)
}

func (list *List[T]) Pop(index int) {
	list.Data = append(list.Data[:index], list.Data[index+1:]...)
}

func (list *List[T]) Contains(v T) bool {
	for i := 0; i < len(list.Data); i++ {
		if reflect.DeepEqual(list.Data[i], v) {
			return true
		}
	}
	return false
}

func (list List[T]) Last() T {
	return list.Data[list.Len()-1]
}

func (list *List[T]) Len() int {
	return len(list.Data)
}
