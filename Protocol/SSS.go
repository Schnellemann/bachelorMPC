package protocol

import (
	fields "MPC/Fields"
	netpack "MPC/Netpackage"
	"errors"
	"math"
	"strconv"
)

type ShamirSecretSharing struct {
	field  fields.Field
	poly   []int64
	degree int
}

func makeShamirSecretSharing(secret int64, fieldImp fields.Field, degree int) *ShamirSecretSharing {
	SSS := new(ShamirSecretSharing)
	SSS.degree = degree
	SSS.field = fieldImp
	SSS.poly = make([]int64, degree+1)
	SSS.poly[0] = secret
	//Fill the rest of the polynomial with random values from the field
	for i := 1; i <= degree; i++ {
		SSS.poly[i] = SSS.field.GetRandom()
	}
	return SSS
}

func (s ShamirSecretSharing) lagrangeInterpolation(shares []netpack.Share) (secret int64, err error) {
	if !(len(shares) > s.degree) {
		return int64(0), errors.New("Lagrange: too few shares received")
	}
	result := float64(0)
	//Computation as showed in example in MPC-book page 42
	for i := 0; i <= s.degree; i++ {
		enumerator := float64(shares[i].Value)
		denominator := float64(1)
		for j := 0; j <= s.degree; j++ {
			if j != i {
				identifierI := shares[i].Identifier
				identifierJ := shares[j].Identifier
				// (┛ಠ_ಠ)┛彡┻━┻ (┛ಠ_ಠ)┛彡┻━┻ (┛ಠ_ಠ)┛彡┻━┻ (┛ಠ_ಠ)┛彡┻━┻ (┛ಠ_ಠ)┛彡┻━┻
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

func (s ShamirSecretSharing) makeShares(numberOfParties int64) (shares []netpack.Share) {
	for i := 1; i <= int(numberOfParties); i++ {
		share := new(netpack.Share)
		share.Identifier = "p" + strconv.Itoa(i)
		for j, v := range s.poly {
			share.Value = s.field.Add(share.Value, s.field.Multiply(v, s.field.Pow(int64(i), int64(j))))
		}
		shares = append(shares, *share)
	}
	return
}
