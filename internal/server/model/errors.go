package models

import (
	"errors"
)

var (
	ErrorGetValue = errors.New("Ошибка получения значения")
	ErrorDiffType = errors.New("Ошибка: переменная в базе другого типа")
	ErrorNotDB    = errors.New("Ошибка: переменной нет в базе")
	ErrorUnType   = errors.New("Есть непредвиденный тип данных в базе")
)
