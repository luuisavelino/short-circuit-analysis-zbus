package zbus

import (
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
)

type Posicao_zbus struct {
	Posicao int
	Tipo    string
}

type Matrix [][]complex128

func MontaZbus(ElementsSequencia map[string]models.Elements, SystemSize map[string]int) (models.ZbusStr, map[string]Posicao_zbus, error) {
	zbus_positiva, zbus_zero, barras_adicionadas, posicao := AdicionaElementosTipo1(ElementsSequencia["1"], SystemSize["size"])
	zbus_positiva, zbus_zero, barras_adicionadas, elementosTipo3 := AdicionaElementosTipo2(ElementsSequencia["2"], zbus_positiva, zbus_zero, barras_adicionadas, posicao)
	zbus_positiva, zbus_zero, barras_adicionadas = AdicionaElementosTipo3(elementosTipo3, zbus_positiva, zbus_zero, barras_adicionadas, SystemSize["size"])

	zbus := models.ZbusStr{
		Positiva: zbus_positiva.ArrayCmplxToArrayStr(),
		Negativa: zbus_positiva.ArrayCmplxToArrayStr(),
		Zero:     zbus_zero.ArrayCmplxToArrayStr(),
	}

	return zbus, barras_adicionadas, nil
}

func AdicionaElementosTipo1(Elements models.Elements, SystemSize int) (Matrix, Matrix, map[string]Posicao_zbus, int) {
	var posicao int = 0
	var barras_adicionadas = make(map[string]Posicao_zbus)
	zbus_positiva, _ := Preenche_matriz_com_zeros(SystemSize)
	zbus_zero, _ := Preenche_matriz_com_zeros(SystemSize)

	for _, dados_linha := range Elements {
		zbus_positiva.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_positiva)
		zbus_zero.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_zero)

		barras_adicionadas[dados_linha.De] = Posicao_zbus{
			Posicao: posicao,
			Tipo:    "1",
		}

		posicao++
	}

	return zbus_positiva, zbus_zero, barras_adicionadas, posicao
}

func AdicionaElementosTipo2(Elements models.Elements, zbus_positiva, zbus_zero Matrix, barras_adicionadas map[string]Posicao_zbus, posicao int) (Matrix, Matrix, map[string]Posicao_zbus, []models.Element) {
	var elementosTipo3 []models.Element

	for len(Elements) != 0 {
		for nome_linha, linha := range Elements {
			_, existe_de := barras_adicionadas[linha.De]
			_, existe_para := barras_adicionadas[linha.Para]

			if existe_de && existe_para {
				elementosTipo3 = append(elementosTipo3, linha)
				delete(Elements, nome_linha)

			} else if existe_de {
				zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_positiva)
				zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_zero)

				barras_adicionadas[linha.Para] = Posicao_zbus{
					Posicao: posicao,
					Tipo:    "2",
				}

				delete(Elements, nome_linha)
				posicao++

			} else if existe_para {
				zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_positiva)
				zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_zero)

				barras_adicionadas[linha.De] = Posicao_zbus{
					Posicao: posicao,
					Tipo:    "2",
				}

				delete(Elements, nome_linha)
				posicao++
			}
		}
	}

	return zbus_positiva, zbus_zero, barras_adicionadas, elementosTipo3
}

func AdicionaElementosTipo3(elementosTipo3 []models.Element, zbus_positiva, zbus_zero Matrix, barras_adicionadas map[string]Posicao_zbus, SystemSize int) (Matrix, Matrix, map[string]Posicao_zbus) {
	for x := 0; x < len(elementosTipo3); x++ {
		linha := elementosTipo3[x]

		zbus_positiva = zbus_positiva.AdicionaElementoTipo3ComReducaoDeKron(
			barras_adicionadas[linha.De].Posicao,
			barras_adicionadas[linha.Para].Posicao,
			linha.Z_positiva,
			SystemSize)
		zbus_zero = zbus_zero.AdicionaElementoTipo3ComReducaoDeKron(
			barras_adicionadas[linha.De].Posicao,
			barras_adicionadas[linha.Para].Posicao,
			linha.Z_zero,
			SystemSize)
	}

	return zbus_positiva, zbus_zero, barras_adicionadas
}
