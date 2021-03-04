package fields

type field interface {
	multiply(int64, int64) int64
	add(int64, int64) int64
	getRandom() int64
}
