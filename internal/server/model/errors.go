package models

import (
	"errors"
)

var (
	ErrorGetValue = errors.New("ошибка получения значения")
	ErrorDiffType = errors.New("ошибка: переменная в базе другого типа")
	ErrorNotDB    = errors.New("ошибка: переменной нет в базе")
	ErrorUnType   = errors.New("есть непредвиденный тип данных в базе")
)
