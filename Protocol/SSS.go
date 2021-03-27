package protocol

import (
	fields "MPC/Fields"
	netpack "MPC/Netpackage"
	"errors"
	"math"
	"strconv"
)

type ShamirSecretSharing struct {
	field fields.Field
	poly  []int64
}

func makeShamirSecretSharing(secret int64, fieldImp fields.Field, degree int) *ShamirSecretSharing {
	SSS := new(ShamirSecretSharing)
	SSS.field = fieldImp
	SSS.poly = make([]int64, degree+1)
	SSS.poly[0] = secret
	//Fill the rest of the polynomial with random values from the field
	for i := 1; i <= degree; i++ {
		SSS.poly[i] = SSS.field.GetRandom()
	}
	return SSS
}

func (s ShamirSecretSharing) lagrangeInterpolation(shares []netpack.Share, degree int) (secret int64, err error) {
	if !(len(shares) > degree) {
		return int64(0), errors.New("Lagrange: too few shares received")
	}
	result := float64(0)
	//Computation as showed in example in MPC-book page 42
	for i := 0; i <= degree; i++ {
		enumerator := float64(shares[i].Value)
		denominator := float64(1)
		for j := 0; j <= degree; j++ {
			if j != i {
				identifierI := shares[i].Identifier
				identifierJ := shares[j].Identifier
				numberI, _ := strconv.Atoi(string(identifierI[1]))
				numberJ, _ := strconv.Atoi(string(identifierJ[1]))

				enumerator = enumerator * float64(-numberJ)
				denominator = denominator * float64(numberI-numberJ)
			}
		}
		result += enumerator / denominator
	}
	return s.field.Convert((int64(math.Round(result)))), nil
}

func (s ShamirSecretSharing) makeShares(numberOfParties int64, identifier string) (shares []netpack.Share) {
	for i := 1; i <= int(numberOfParties); i++ {
		share := new(netpack.Share)
		//TODO this is wrong should be my party num always
		share.Identifier = identifier
		for j, v := range s.poly {
			share.Value = s.field.Add(share.Value, s.field.Multiply(v, s.field.Pow(int64(i), int64(j))))
		}
		shares = append(shares, *share)
	}
	return
}
