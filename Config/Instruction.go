package config

type Instruction struct {
	Op     Operator
	Left   string
	Right  string
	Result string
}

type InstructionTree struct {
	Left        *InstructionTree
	Right       *InstructionTree
	Instruction *Instruction
}

type Operator int

const (
	Add = iota
	Multiply
	Scalar
)

func (op Operator) String() string {
	return [...]string{"add", "mul", "scalar"}[op]
}

func createAddition(left string, right string, result string) Instruction {
	return Instruction{Add, left, right, result}
}

func createMultiplication(left string, right string, result string) Instruction {
	return Instruction{Multiply, left, right, result}
}

func createScalar(left string, right string, result string) Instruction {
	return Instruction{Scalar, left, right, result}
}
