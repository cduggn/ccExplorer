package aws

import "strconv"

func ConvertToFloat(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func isEmpty(s []string) string {
	if len(s) == 1 {
		return ""
	} else {
		return s[1]
	}

}
