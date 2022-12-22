package core

import (
	"encoding/binary"
)

type Instruction byte

const (
	InstructionPushInt  Instruction = 0x0a //10
	InstructionAdd      Instruction = 0x0b //11
	InstructionPushByte Instruction = 0x0c
	InstructionPack     Instruction = 0x0d
	InstructionSub      Instruction = 0x0e
	InstructionStore    Instruction = 0x0f
	InstructionGet      Instruction = 0xae
	InstructionMul      Instruction = 0xea
	InstructionDiv      Instruction = 0xfd
)

type Stack struct {
	data []any
	sp   int
}

// preallocate, no need to grow as we go, takes little bit of memory upfront but saves us some space
func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

func (s *Stack) Push(val any) {
	s.data = append([]any{val}, s.data...)

	// s.data[s.sp] = val
	s.sp++
}

func (s *Stack) Pop() any {
	val := s.data[0]
	//shift array back one slot
	s.data = append(s.data[:0], s.data[1:]...)
	s.sp--

	return val
}

type VM struct {
	data          []byte
	ip            int //instruction pointer
	stack         *Stack
	contractState *State
}

func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		data:          data,
		ip:            0,
		stack:         NewStack(128),
		contractState: contractState,
	}
}

func (vm *VM) Run() error {
	for {
		instructionP := Instruction(vm.data[vm.ip])

		if err := vm.Exec(instructionP); err != nil {
			return err
		}

		vm.ip++

		if vm.ip > len(vm.data)-1 {
			break
		}
	}
	return nil
}

func (vm *VM) Exec(instruction Instruction) error {
	switch instruction {
	case InstructionGet:
		key := vm.stack.Pop().([]byte)

		val, err := vm.contractState.Get(key)
		if err != nil {
			return err
		}

		vm.stack.Push(val)

	case InstructionStore:
		var (
			key             = vm.stack.Pop().([]byte)
			value           = vm.stack.Pop()
			serializedValue []byte
		)

		switch val := value.(type) {
		case int:
			serializedValue = serializeInt64(int64(val))
		default:
			panic("TODO: Unknown type")
		}

		vm.contractState.Put(key, serializedValue)

	case InstructionPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1]))

	case InstructionPushByte:
		vm.stack.Push(byte(vm.data[vm.ip-1]))

	case InstructionPack:
		n := vm.stack.Pop().(int)

		b := make([]byte, n)

		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte)
		}

		vm.stack.Push(b)

	case InstructionSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a - b
		vm.stack.Push(c)

	case InstructionAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a + b
		vm.stack.Push(c)

	case InstructionMul:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a * b
		vm.stack.Push(c)

	case InstructionDiv:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a / b
		vm.stack.Push(c)
	}

	return nil
}

func serializeInt64(value int64) []byte {
	buffer := make([]byte, 8)

	binary.LittleEndian.PutUint64(buffer, uint64(value))

	return buffer
}

func deserializeInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}
