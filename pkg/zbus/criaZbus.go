package zbus

import (
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
)

type Posicao_zbus struct {
	Posicao int
}

type Matrix [][]complex128

func MontaZbus() (models.ZbusStr, map[string]Posicao_zbus) {
	var zbus_positiva, _ = Preenche_matriz_com_zeros(models.SystemSize["size"])
	var zbus_zero, _ = Preenche_matriz_com_zeros(models.SystemSize["size"])

	var elementos_tipo_3 []models.Element
	var barras_adicionadas = make(map[string]Posicao_zbus)
	var posicao = 0

	for _, dados_linha := range models.Elements["1"] {
		zbus_positiva = zbus_positiva.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_positiva)
		zbus_zero = zbus_zero.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_zero)

		barras_adicionadas[dados_linha.De] = Posicao_zbus{
			Posicao: posicao,
		}

		posicao++
	}

	for len(models.Elements["2"]) != 0 {
		for nome_linha, linha := range models.Elements["2"] {
			_, existe_de := barras_adicionadas[linha.De]
			_, existe_para := barras_adicionadas[linha.Para]

			if existe_de && existe_para {
				elementos_tipo_3 = append(elementos_tipo_3, linha)
				delete(models.Elements["2"], nome_linha)

			} else if existe_de {
				zbus_positiva = zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_positiva)
				zbus_zero = zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_zero)

				barras_adicionadas[linha.Para] = Posicao_zbus{
					Posicao: posicao,
				}

				delete(models.Elements["2"], nome_linha)
				posicao++

			} else if existe_para {
				zbus_positiva = zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_positiva)
				zbus_zero = zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_zero)

				barras_adicionadas[linha.De] = Posicao_zbus{
					Posicao: posicao,
				}

				delete(models.Elements["2"], nome_linha)
				posicao++

			}
		}
	}

	for x := 0; x < len(elementos_tipo_3); x++ {
		linha := elementos_tipo_3[x]

		zbus_positiva = zbus_positiva.AdicionaElementoTipo3ComReducaoDeKron(
			barras_adicionadas[linha.De].Posicao,
			barras_adicionadas[linha.Para].Posicao,
			linha.Z_positiva,
			models.SystemSize["size"])
		zbus_zero = zbus_zero.AdicionaElementoTipo3ComReducaoDeKron(
			barras_adicionadas[linha.De].Posicao,
			barras_adicionadas[linha.Para].Posicao,
			linha.Z_zero,
			models.SystemSize["size"])
	}

	zbus := models.ZbusStr{
		Positiva: zbus_positiva.ArrayCmplxToArrayStr(),
		Negativa: zbus_positiva.ArrayCmplxToArrayStr(),
		Zero:     zbus_zero.ArrayCmplxToArrayStr(),
	}

	return zbus, barras_adicionadas
}
