package models

type Element struct {
	De         string `json:"de"`
	Para       string `json:"para"`
	Nome       string `json:"nome"`
	Z_positiva string `json:"z_positiva"`
	Z_zero     string `json:"z_zero"`
}

var Elements = make(map[string]map[string]Element)
var SystemSize = make(map[string]int)

type MatrixStr [][]string

type ZbusStr struct {
	Positiva MatrixStr `json:"positiva"`
	Negativa MatrixStr `json:"negativa"`
	Zero     MatrixStr `json:"zero"`
}

var Zbus ZbusStr




