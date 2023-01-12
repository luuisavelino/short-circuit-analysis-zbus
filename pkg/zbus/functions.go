package zbus

import (
	"strconv"

	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
)

func (m Matrix) ArrayCmplxToArrayStr() models.MatrixStr {
	size := len(m)
	_, matrixStr := Preenche_matriz_com_zeros(size)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			matrixStr[x] = append(matrixStr[y], strconv.FormatComplex(m[x][y], 'g', 'g', 64))
		}
		matrixStr = append(matrixStr, matrixStr[x])
	}

	return matrixStr
}

func Preenche_matriz_com_zeros(tamanho int) (Matrix, models.MatrixStr) {
	var matrix = make(Matrix, 0)
	var matrix_str = make(models.MatrixStr, 0)

	// Adiciona elementos 0 na matriz zbus
	for i := 0; i < tamanho; i++ {
		temp1 := make([]complex128, 0)
		temp2 := make([]string, 0)
		for j := 0; j < tamanho; j++ {
			temp1 = append(temp1, 0.0)
			temp2 = append(temp2, "0.0")
		}

		matrix = append(matrix, Matrix{temp1}...)
		matrix_str = append(matrix_str, models.MatrixStr{temp2}...)
	}

	return matrix, matrix_str
}

func Aumenta_tamanho_da_matriz(matrix Matrix, tamanho int) Matrix {

	for x := 0; x < tamanho; x++ {
		temp := make([]complex128, 0)
		for j := 0; j <= len(matrix); j++ {
			temp = append(temp, 0)
		}

		matrix = append(matrix, Matrix{temp}...)

		for i := 0; i < len(matrix); i++ {
			matrix[i] = append(matrix[i], 0)
		}
	}

	return matrix
}
