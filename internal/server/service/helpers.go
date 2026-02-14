package service

import "strconv"

func getValueFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func getValueInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
