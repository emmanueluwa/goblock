package core

type Instruction byte

const (
	InstructionPush Instruction = 0x0a //10
	InstructionAdd  Instruction = 0x0b //11
)

type VM struct {
	data  []byte
	ip    int //instruction pointer
	stack []byte
	sp    int //stacl pointer
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: make([]byte, 1024),
		sp:    -1,
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
	case InstructionPush:
		vm.pushStack(vm.data[vm.ip-1])
	case InstructionAdd:
		a := vm.stack[0]
		b := vm.stack[1]
		c := a + b
		vm.pushStack(c)
	}

	return nil
}

func (vm *VM) pushStack(b byte) {
	vm.sp++
	vm.stack[vm.sp] = b
}
