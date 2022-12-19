package models

type Element struct {
	Id         int    `json:"id"`
	De         string `json:"de"`
	Para       string `json:"para"`
	Nome       string `json:"nome"`
	Z_positiva string `json:"z_positiva"`
	Z_zero     string `json:"z_zero"`
}

type MatrixStr [][]string

type ZbusStr struct {
	Positiva MatrixStr `json:"positiva"`
	Negativa MatrixStr `json:"negativa"`
	Zero     MatrixStr `json:"zero"`
}

var Zbus ZbusStr
