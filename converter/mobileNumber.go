package converter

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// austrianMobileProviders see https://www.telefonabc.at/vorwahl_uebersicht.aspx
var austrianMobileProviders = []string{
	"0664", // A1
	"0676", // Magenta
	"0678", // Magenta (was UPC)
	"0660", // 3 Drei.
	"0699", // 3 Drei.
	"0650", // tele.ring
	"0680", // bob
	"0699", // yesss
	"0681", // yesss
	"0667", // m:tel
}

// GetMobileNumber extracts the mobile number without blanks.
// The + sign is allowed and returned, too.
func GetMobileNumber(input string) (string, error) {
	var b bytes.Buffer
	for i := 0; i < len(input); i++ {
		c := input[i]
		// valid chars are passed through
		if '0' <= c && c <= '9' || c == '+' {
			b.WriteByte(c)
			continue
		}
	}
	return b.String(), nil
}

// PrepareMobileNumber returns 0043664.... and nil when
// input could be treated as mobile number. The input parameter
// can either be a string, int or float64 value
func PrepareMobileNumber(input interface{}) (string, error) {
	var mobileNumber string
	var isString bool
	switch val := input.(type) {
	case string:
		isString = true
		var err error
		mobileNumber, err = GetMobileNumber(val)
		if err != nil {
			return "", err
		}
	case int:
		mobileNumber = "0" + fmt.Sprintf("%d", val)
	case float64:
		mobileNumber = "0" + fmt.Sprintf("%.0f", val)
	default:
		return "", errors.New("unsupported input")
	}

	// check int and float
	if !isString {
		for _, provider := range austrianMobileProviders {
			if strings.HasPrefix(mobileNumber, provider) {
				mobileNumber = strings.Replace(mobileNumber, provider, "0043"+provider[1:], 1)
				break
			}
		}
		if !strings.HasPrefix(mobileNumber, "00") {
			mobileNumber = "0" + mobileNumber
		}
	}

	if len(mobileNumber) < 8 {
		return "", fmt.Errorf("mobile number %q is too short", mobileNumber)
	}

	// check +43664... and convert it to 0043664...
	if strings.HasPrefix(mobileNumber, "+") {
		mobileNumber = strings.Replace(mobileNumber, "+", "00", 1)
	}

	// check austrian providers only
	for _, provider := range austrianMobileProviders {
		p := provider[1:]
		if strings.HasPrefix(mobileNumber, provider) {
			mobileNumber = strings.Replace(mobileNumber, provider, "0043"+p, 1)
			return mobileNumber, nil
		}
		if strings.HasPrefix(mobileNumber, p) {
			mobileNumber = strings.Replace(mobileNumber, p, "0043"+p, 1)
			return mobileNumber, nil
		}
	}

	if !strings.HasPrefix(mobileNumber, "00") {
		mobileNumber = "00" + mobileNumber
	}

	return mobileNumber, nil
}
