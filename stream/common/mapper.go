package common

import (
	"fmt"
	"strconv"
)

func StringToInt(elem string) (int, error) {
	return strconv.Atoi(elem)
}

func IntToString(elem int) (string, error) {
	return strconv.Itoa(elem), nil
}

func StringToFloat(elem string) (float64, error) {
	return strconv.ParseFloat(elem, 64)
}

func FloatToString(elem float64) (string, error) {
	return fmt.Sprintf("%g", elem), nil
}

func IntToFloat64(elem int) (float64, error) {
	return float64(elem), nil
}
