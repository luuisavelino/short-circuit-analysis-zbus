package zbus

import (
	"fmt"

	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
)

type Posicao_zbus struct {
	Posicao int
}

type Matrix [][]complex128

var posicao int
var zbus_positiva, zbus_zero Matrix
var barras_adicionadas map[string]Posicao_zbus
var elementosTipo3 []models.Element

func MontaZbus() (map[string]Posicao_zbus, error) {

	posicao = 0
	zbus_positiva, _ = Preenche_matriz_com_zeros(models.SystemSize["size"])
	zbus_zero, _ = Preenche_matriz_com_zeros(models.SystemSize["size"])
	barras_adicionadas = make(map[string]Posicao_zbus)

	fmt.Println("tipo1")
	AdicionaElementosTipo1()
	fmt.Println("tipo2")
	AdicionaElementosTipo2()
	fmt.Println("tipo3")
	AdicionaElementosTipo3()
	//fmt.Println("fim")

	models.Zbus = models.ZbusStr{
		Positiva: zbus_positiva.ArrayCmplxToArrayStr(),
		Negativa: zbus_positiva.ArrayCmplxToArrayStr(),
		Zero:     zbus_zero.ArrayCmplxToArrayStr(),
	}

	return barras_adicionadas, nil
}

func AdicionaElementosTipo1() {
	for _, dados_linha := range models.Elements["1"] {
		zbus_positiva.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_positiva)
		zbus_zero.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_zero)

		barras_adicionadas[dados_linha.De] = Posicao_zbus{
			Posicao: posicao,
		}

		posicao++
	}
}

func AdicionaElementosTipo2() {
	for len(models.Elements["2"]) != 0 {
		for nome_linha, linha := range models.Elements["2"] {
			_, existe_de := barras_adicionadas[linha.De]
			_, existe_para := barras_adicionadas[linha.Para]

			if existe_de && existe_para {
				elementosTipo3 = append(elementosTipo3, linha)
				delete(models.Elements["2"], nome_linha)

			} else if existe_de {
				zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_positiva)
				zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_zero)

				barras_adicionadas[linha.Para] = Posicao_zbus{
					Posicao: posicao,
				}

				delete(models.Elements["2"], nome_linha)
				posicao++

			} else if existe_para {
				zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_positiva)
				zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_zero)

				barras_adicionadas[linha.De] = Posicao_zbus{
					Posicao: posicao,
				}

				delete(models.Elements["2"], nome_linha)
				posicao++
			}
		}
	}
}

func AdicionaElementosTipo3() {
	for x := 0; x < len(elementosTipo3); x++ {
		linha := elementosTipo3[x]

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
}
