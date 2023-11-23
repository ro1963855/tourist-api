package utils

import "strconv"

func StringToUnit(input string) (uint, error) {
	transformedUint, err := strconv.ParseUint(input, 10, 64)
	return uint(transformedUint), err
}
