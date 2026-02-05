package service

import "strconv"

func getValueFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func getValueInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func getStringFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func getStringInt(i int64) string {
	return strconv.FormatInt(i, 10)
}
