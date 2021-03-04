package fields

import (
	"crypto/rand"
	"math/big"
)

type modPrime struct {
	p int64
}

func makeModPrime(prime int64) *modPrime {
	mk := new(modPrime)
	mk.p = prime
	return mk
}

func (mp modPrime) multiply(a int64, b int64) int64 {
	result := (a * b) % mp.p
	return result
}

func (mp modPrime) add(a int64, b int64) int64 {
	result := (a + b) % mp.p
	return result
}

func (mp modPrime) getRandom() int64 {
	bigP := big.NewInt(mp.p)
	randomNumber, _ := rand.Int(rand.Reader, bigP)
	return randomNumber.Int64()
}
