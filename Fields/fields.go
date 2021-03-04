package fields

type Field interface {
	Multiply(int64, int64) int64
	Add(int64, int64) int64
	GetRandom() int64
}
