package aux

func Remove(s []interface{}, i int) []interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
