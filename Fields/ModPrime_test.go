package fields

import (
	"testing"
)

func TestDivide(t *testing.T) {
	mp := MakeModPrime(7)
	result := mp.Divide(3, 2)
	if result != 5 {
		t.Errorf("Wrong division in modPrime, got %v , should have gotten 5", result)
	}
}
