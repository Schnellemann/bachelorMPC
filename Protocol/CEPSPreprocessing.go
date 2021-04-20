package protocol

import (
	netpack "MPC/Netpackage"
	"fmt"
	"math"
	"strconv"
)

type randoms struct {
	r1t int64
	r2t int64
}

func (prot *Ceps) createMatrix() [][]int64 {
	field := prot.shamir.field
	n := int(prot.config.ConstantConfig.NumberOfParties)
	m := int(float64(n) - (math.Ceil(float64(n)/float64(2)) - float64(1)))
	matrix := make([][]int64, m)
	for i := 1; i <= m; i++ {
		matrix[i-1] = make([]int64, n)
		for j := 1; j <= n; j++ {
			enumerator := int64(1)
			denominator := int64(1)
			for k := 1; k <= n; k++ {
				if k != j {
					beta_i := int64(n + i)
					enumerator = field.Multiply(enumerator, field.Minus(beta_i, int64(k)))
					denominator = field.Multiply(denominator, field.Minus(int64(j), int64(k)))
				}
				matrix[i-1][j-1] = field.Divide(enumerator, denominator)
			}
		}
	}
	return matrix
}

func (prot *Ceps) createRValues(fShares []*netpack.Share, gShares []*netpack.Share) {
	if len(fShares) != len(gShares) {
		fmt.Println("Prot/preprocess: the length of the two list are not equal")
		return
	}
	field := prot.shamir.field
	//rowwise in matrix compute pair of r-values
	for _, row := range prot.matrix {
		var resultF int64
		var resultG int64
		for i, share := range fShares {
			resultF = field.Add(resultF, field.Multiply(share.Value, row[i]))
		}
		for i, share := range gShares {
			resultG = field.Add(resultG, field.Multiply(share.Value, row[i]))
		}
		prot.listOfRandoms = append(prot.listOfRandoms, randoms{r1t: resultF, r2t: resultG})
	}

}

func (prot *Ceps) runPreprocess() {
	numberOfMults := prot.instructionTree.CountMults()
	rounds := int(math.Ceil(float64(numberOfMults) / (prot.config.ConstantConfig.NumberOfParties - float64(prot.degree))))
	for i := 1; i <= rounds; i++ {
		secret := prot.shamir.field.GetRandom()
		polyt := makeShamirSecretSharing(secret, prot.shamir.field, prot.degree)
		poly2t := makeShamirSecretSharing(secret, prot.shamir.field, 2*prot.degree)
		fIden := netpack.ShareIdentifier{Ins: "f" + strconv.Itoa(i), PartyNr: int(prot.config.VariableConfig.PartyNr)}
		gIden := netpack.ShareIdentifier{Ins: "g" + strconv.Itoa(i), PartyNr: int(prot.config.VariableConfig.PartyNr)}
		fShares := polyt.makeShares(int64(prot.config.ConstantConfig.NumberOfParties), fIden)
		gShares := poly2t.makeShares(int64(prot.config.ConstantConfig.NumberOfParties), gIden)
		prot.handleShare(fShares)
		prot.handleShare(gShares)
	}

	for i := 1; i <= rounds; i++ {
		//Wait for shares
		var randomfIdentifiers []netpack.ShareIdentifier
		var randomgIdentifiers []netpack.ShareIdentifier
		for j := 1; j <= int(prot.config.ConstantConfig.NumberOfParties); j++ {
			randomfIdentifiers = append(randomfIdentifiers, netpack.ShareIdentifier{Ins: "f" + strconv.Itoa(i), PartyNr: j})
			randomgIdentifiers = append(randomgIdentifiers, netpack.ShareIdentifier{Ins: "g" + strconv.Itoa(i), PartyNr: j})
		}
		prot.waitForShares(randomfIdentifiers)
		prot.waitForShares(randomgIdentifiers)
		//Handle the shares to create r-values
		var fShares []*netpack.Share
		var gShares []*netpack.Share
		prot.rShares.mu.Lock()
		for _, i := range randomfIdentifiers {
			fShares = append(fShares, prot.rShares.receivedShares[i])
		}
		for _, i := range randomgIdentifiers {
			gShares = append(gShares, prot.rShares.receivedShares[i])
		}
		prot.rShares.mu.Unlock()
		prot.createRValues(fShares, gShares)

	}

}
