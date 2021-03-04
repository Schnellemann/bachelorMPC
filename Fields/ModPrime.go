package fields

import (
	"crypto/rand"
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

func (mp ModPrime) GetRandom() int64 {
	bigP := big.NewInt(mp.p)
	randomNumber, _ := rand.Int(rand.Reader, bigP)
	return randomNumber.Int64()
}
