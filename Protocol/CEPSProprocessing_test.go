package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	"testing"
)

func TestCreateMatrix(t *testing.T) {
	r1 := []int64{1, 8, 3}
	r2 := []int64{3, 3, 6}
	testMatrix := [][]int64{r1, r2}

	configs := config.MakeConfigs(ip, "p1+p2+p3", []int{4, 3, 3})
	peerlist := getXPeers(configs)
	//Make protocol
	prot := mkProtocol(configs[0], field.MakeModPrime(11), peerlist[0])
	n := int(prot.config.ConstantConfig.NumberOfParties)
	m := n - prot.degree

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if prot.matrix[i][j] != testMatrix[i][j] {
				t.Errorf("Matrix on entry %v,%v is expected to be %v but is %v", i+1, j+1, testMatrix[i][j], prot.matrix[i][j])
			}
		}
	}

}
