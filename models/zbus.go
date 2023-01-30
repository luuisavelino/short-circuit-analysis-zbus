package models

type Element struct {
	De         string `json:"de"`
	Para       string `json:"para"`
	Nome       string `json:"nome"`
	Z_positiva string `json:"z_positiva"`
	Z_zero     string `json:"z_zero"`
}

type Elements map[string]Element

type MatrixStr [][]string

type ZbusStr struct {
	Positiva MatrixStr `json:"positiva"`
	Negativa MatrixStr `json:"negativa"`
	Zero     MatrixStr `json:"zero"`
}

type Posicao_zbus struct {
	Posicao int
}

var ReadOps uint64
var WriteOps uint64
