package config

type instruction struct {
	op     operator
	left   string
	right  string
	result string
}

type operator int

const (
	add = iota
	multiply
	scalar
)

func (op operator) String() string {
	return [...]string{"add", "mul", "scalar"}[op]
}

func createAddition(left string, right string, result string) instruction {
	return instruction{add, left, right, result}
}

func createMultiplication(left string, right string, result string) instruction {
	return instruction{multiply, left, right, result}
}

func createScalar(left string, right string, result string) instruction {
	return instruction{scalar, left, right, result}
}
