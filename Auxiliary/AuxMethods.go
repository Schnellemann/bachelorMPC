package aux

import (
	netpack "MPC/Netpackage"
)

func Remove(s []interface{}, i int) []interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func RemoveDuplicateValues(slice []netpack.ShareIdentifier) []netpack.ShareIdentifier {
	keys := make(map[netpack.ShareIdentifier]bool)
	list := []netpack.ShareIdentifier{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
