package protocol

import (
	fields "MPC/Fields"
	netpack "MPC/Netpackage"
	"strconv"
	"testing"
)

//Testing the setup of SSS, with field Mod Prime
func TestModPrime(t *testing.T) {
	field := fields.MakeModPrime(43)
	degree := 10
	SSS := makeShamirSecretSharing(42, field, degree)
	if len(SSS.poly) != degree+1 {
		t.Errorf("The length of poly does not match degree, expected %v, got %v", degree+1, len(SSS.poly))
	}
	if SSS.poly[0] != 42 {
		t.Errorf("The polynomial does not have the secret as the constant")
	}
	for i := 1; i < degree; i++ {
		if SSS.poly[i] >= 43 || SSS.poly[i] < 0 {
			t.Errorf("The coeffecients in the polynomials are not in the range of the prime")
		}
	}
}

func TestShareMaking(t *testing.T) {
	SSS := makeShamirSecretSharing(42, fields.MakeModPrime(43), 1)
	SSS.poly = []int64{4, 7, 9, 1, 11}
	shares := SSS.makeShares(5, "p1")
	if shares[0].Value != (4+7+9+1+11)%43 {
		t.Errorf("First share not computed correctly got %v expected %v", shares[0].Value, (4+7+9+1+11)%43)
	}
	if shares[1].Value != (4+7*2+9*4+1*8+11*16)%43 {
		t.Errorf("First share not computed correctly got %v expected %v", shares[1].Value, (4+7*2+9*4+1*8+11*16)%43)
	}
	if shares[2].Value != (4+7*3+9*9+1*27+11*81)%43 {
		t.Errorf("First share not computed correctly got %v expected %v", shares[2].Value, (4+7*3+9*9+1*27+11*81)%43)
	}
	if shares[3].Value != (4+7*4+9*16+1*64+11*256)%43 {
		t.Errorf("First share not computed correctly got %v expected %v", shares[3].Value, (4+7*4+9*16+1*64+11*256)%43)
	}
	if shares[4].Value != (4+7*5+9*25+1*125+11*625)%43 {
		t.Errorf("First share not computed correctly got %v expected %v", shares[4].Value, (4+7*5+9*25+1*125+11*625)%43)
	}
}

func TestLagrangeInterpolation(t *testing.T) {
	secret := int64(28)
	SSS := makeShamirSecretSharing(secret, fields.MakeModPrime(43), 2)
	shares := SSS.makeShares(5, "p1")
	for i := 0; i < 5; i++ {
		shares[i].Identifier = "p" + strconv.Itoa(i+1)
	}
	secretLagrange, _ := SSS.lagrangeInterpolation(shares, 2)
	if secretLagrange != secret {
		t.Errorf("Wrong secret, got %v expected %v", secretLagrange, secret)
	}
}

func TestLagrangeInterpolation2(t *testing.T) {
	secret := int64(11)
	SSS := makeShamirSecretSharing(secret, fields.MakeModPrime(13), 1)
	SSS.poly = []int64{11, 8}
	shares := []netpack.Share{{Value: 9, Identifier: "o3"}, {Value: 6, Identifier: "o1"}}
	secretLagrange, _ := SSS.lagrangeInterpolation(shares, 1)
	if secretLagrange != secret {
		t.Errorf("Wrong secret, got %v expected %v", secretLagrange, secret)
	}
}
