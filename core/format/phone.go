package format

import (
	"strings"
)

// Will format phone into "0" prefix phone number
func Phone(phone string) string {
	if strings.HasPrefix(phone, "62") {
		return phone
	}
	if strings.HasPrefix(phone, "0") {
		return strings.Replace(phone, "0", "62", 1)
	}
	if strings.HasPrefix(phone, "+62") {
		return strings.Replace(phone, "+62", "62", 1)
	}
	return phone
}
