package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(test *testing.T) {
	//1 push to stack, 2 push to stack
	data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVM(data)
	assert.Nil(test, vm.Run())

	// 1 + 2 = 3
	assert.Equal(test, byte(3), vm.stack[vm.sp])
}
