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

	// push FOO to stack, pack, push 5 to stack, store
	// data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x03, 0x0a, 0x0d, 0x0f}
	dataOther := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4d, 0x0c, 0x4f, 0x0c, 0x46, 0x03, 0x0a, 0x0d, 0x0f}

	data = append(data, dataOther...)

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

	fmt.Printf("%+v\n", contractState)

	valBytes, err := contractState.Get([]byte("FOO"))
	assert.Nil(test, err)
	val := deserializeInt64(valBytes)
	assert.Equal(test, val, int64(5))
}

func TestVMGet(test *testing.T) {
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x03, 0x0a, 0x0d, 0xae}

	data = append(data, pushFoo...)

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(test, vm.Run())

	fmt.Printf("%+v", vm.stack.data)
	val := vm.stack.Pop().([]byte)
	valSerialized := deserializeInt64(val)

	assert.Equal(test, valSerialized, int64(5))

	// valBytes, err := contractState.Get([]byte("FOO"))
	// assert.Nil(test, err)
	// val := deserializeInt64(valBytes)
	// assert.Equal(test, val, int64(5))
}
