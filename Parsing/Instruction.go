package parsing

type Instruction interface {
}

type InstructionTree struct {
	Left        *InstructionTree
	Right       *InstructionTree
	Instruction Instruction
}

func createAddition(left string, right string, result string) *AddInstruction {
	return &AddInstruction{left, right, result}
}

func createMultiplication(left string, right string, result string, unique int) *MultInstruction {
	return &MultInstruction{left, right, result, unique}
}

func createScalar(scalar int, variable string, result string) *ScalarInstruction {
	return &ScalarInstruction{scalar, variable, result}
}
