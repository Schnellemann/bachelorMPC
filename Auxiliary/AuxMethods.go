package aux

import (
	"regexp"
	"strings"
)

func Remove(s []interface{}, i int) []interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func IsIpv4Regex(ipAddress string) bool {
	ipRegex, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	ipAddress = strings.Trim(ipAddress, " ")
	return ipRegex.MatchString(ipAddress)
}
