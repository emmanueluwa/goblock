package core

type Instruction byte

const (
	InstructionPushInt  Instruction = 0x0a //10
	InstructionAdd      Instruction = 0x0b //11
	InstructionPushByte Instruction = 0x0c
	InstructionPack     Instruction = 0x0d
	InstructionSub      Instruction = 0x0e
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
	s.data[s.sp] = val
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
	data  []byte
	ip    int //instruction pointer
	stack *Stack
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: NewStack(128),
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
	}

	return nil
}
