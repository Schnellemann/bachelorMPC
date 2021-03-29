package fields

import (
	"crypto/rand"
	"math"
	"math/big"
)

type ModPrime struct {
	p int64
}

func MakeModPrime(prime int64) *ModPrime {
	mk := new(ModPrime)
	mk.p = prime
	return mk
}

func (mp ModPrime) Multiply(a int64, b int64) int64 {
	result := (a * b) % mp.p
	return result
}

func (mp ModPrime) Add(a int64, b int64) int64 {
	result := (a + b) % mp.p
	return result
}

func (mp ModPrime) Minus(a int64, b int64) int64 {
	result := (a - b) % mp.p
	return result
}

func (mp ModPrime) Zero() int64 {
	return 0
}

func (mp ModPrime) Neg(a int64) int64 {
	return mp.Minus(mp.Zero(), a)
}

func (mp ModPrime) Pow(a int64, b int64) int64 {
	return int64(math.Pow(float64(a), float64(b))) % mp.p
}

func (mp ModPrime) GetRandom() int64 {
	bigP := big.NewInt(mp.p)
	randomNumber, _ := rand.Int(rand.Reader, bigP)
	return mp.Convert(randomNumber.Int64())
}

func (mp ModPrime) Convert(a int64) int64 {
	if a >= 0 {
		return a % mp.p
	}
	return a%mp.p + mp.p
}
