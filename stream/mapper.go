package stream

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
