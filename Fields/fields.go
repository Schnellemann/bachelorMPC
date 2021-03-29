package fields

type Field interface {
	Multiply(int64, int64) int64
	Add(int64, int64) int64
	Minus(int64, int64) int64
	Zero() int64
	Pow(int64, int64) int64
	GetRandom() int64
	Neg(int64) int64
	Convert(int64) int64
	Divide(int64, int64) int64
}
