package protocol

import (
	netpack "MPC/Netpackage"
	"math"
	"strconv"
)

type randoms struct {
	r1t int64
	r2t int64
}

func (prot *Ceps) runPreprocess(numberOfMults int) {
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
		var randomIdentifiers []netpack.ShareIdentifier
		for j := 1; i <= int(prot.config.ConstantConfig.NumberOfParties); i++ {
			randomIdentifiers = append(randomIdentifiers, netpack.ShareIdentifier{Ins: "f" + strconv.Itoa(i), PartyNr: j})
			randomIdentifiers = append(randomIdentifiers, netpack.ShareIdentifier{Ins: "g" + strconv.Itoa(i), PartyNr: j})
		}
		prot.waitForShares(randomIdentifiers)

	}

}
