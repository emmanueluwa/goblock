package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(test *testing.T) {
	stack := NewStack(101)
	stack.Push(1)
	stack.Push(2)

	val := stack.Pop()
	assert.Equal(test, val, 1)

	val2 := stack.Pop()
	assert.Equal(test, val2, 2)

	fmt.Println(stack.data)
}

func TestVM(test *testing.T) {
	//1 push to stack, 2 push to stack
	//int
	// data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}

	// string
	//0x61 = a
	//FOO
	// data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}

	//subtract
	// data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}

	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f} //, 0x03, 0x0a, 0x02, 0x0a, 0x0e}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(test, vm.Run())

	// 1 + 2 = 3
	// result := vm.stack.Pop()
	// assert.Equal(test, 3, result)

	// result := vm.stack.Pop().([]byte)
	// assert.Equal(test, "FOO", string(result))

	// result := vm.stack.Pop().([]byte)
	// fmt.Printf("%+v\n", vm.stack.data)

	valBytes, err := contractState.Get([]byte("FOO"))
	val := deserializeInt64(valBytes)
	assert.Nil(test, err)
	assert.Equal(test, val, int64(5))

	// assert.Equal(test, 1, result)
}
