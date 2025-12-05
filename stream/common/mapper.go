package common

import (
	"fmt"
	"strconv"
)

// StringToInt converts a string to an int.
// It returns an error if the string cannot be parsed as an integer.
func StringToInt(elem string) (int, error) {
	return strconv.Atoi(elem)
}

// IntToString converts an int to its string representation.
func IntToString(elem int) (string, error) {
	return strconv.Itoa(elem), nil
}

// StringToFloat converts a string to a float64.
// It returns an error if the string cannot be parsed as a floating-point number.
func StringToFloat(elem string) (float64, error) {
	return strconv.ParseFloat(elem, 64)
}

// FloatToString converts a float64 to its string representation.
// It uses the %g format, which removes trailing zeros.
func FloatToString(elem float64) (string, error) {
	return fmt.Sprintf("%g", elem), nil
}

// IntToFloat64 converts an int to a float64.
func IntToFloat64(elem int) (float64, error) {
	return float64(elem), nil
}
