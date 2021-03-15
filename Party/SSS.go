package party

import fields "MPC/Fields"

type ShamirSecretSharing struct {
	field fields.Field
	poly  []int64
}

func makeShamirSecretSharing(secret int64, fieldImp fields.Field, degree int) *ShamirSecretSharing {
	SSS := new(ShamirSecretSharing)
	SSS.field = fieldImp
	SSS.poly = make([]int64, degree)
	SSS.poly[0] = secret
	//Fill the rest of the polynomial with random values from the field
	for i := 1; i < degree; i++ {
		SSS.poly[i] = SSS.field.GetRandom()
	}
	return SSS
}

func (s ShamirSecretSharing) lagrangeInterpolation() {

}

func (s ShamirSecretSharing) sendShares() {

}

func (s ShamirSecretSharing) makeShares() {

}
