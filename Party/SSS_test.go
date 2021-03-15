package party

import (
	fields "MPC/Fields"
	"fmt"
	"testing"
)

//Testing the setup of SSS, with field Mod Prime
func TestSSS(t *testing.T) {
	field := fields.MakeModPrime(43)
	degree := 10
	SSS := makeShamirSecretSharing(42, field, degree)
	fmt.Println(SSS.poly)
	if len(SSS.poly) != degree {
		t.Errorf("The length of poly does not match degree")
	}
	if SSS.poly[0] != 42 {
		t.Errorf("The polynomial does not have the secret as the constant")
	}
	for i := 1; i < degree; i++ {
		if SSS.poly[i] > 43 || SSS.poly[i] <= 0 {
			t.Errorf("The coeffecients in the polynomials are not in the range of the prime")
		}

	}
}
