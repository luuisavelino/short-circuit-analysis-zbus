package zbus

import (
	"fmt"

	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
)

type Posicao_zbus struct {
	Posicao int
}

type Matrix [][]complex128

func MontaZbus(elementos_tipo_1 map[string]models.Element, elementos_tipo_2_3 map[string]models.Element, tamanho_do_sistema int) (models.ZbusStr, map[string]Posicao_zbus) {
	var zbus_positiva, _ = Preenche_matriz_com_zeros(tamanho_do_sistema)
	var zbus_zero, _ = Preenche_matriz_com_zeros(tamanho_do_sistema)

	var elementos_tipo_3 []models.Element
	var barras_adicionadas = make(map[string]Posicao_zbus)
	var posicao = 0

	// Adiciona os elementos do tipo 1
	// Loop passando por todos os elementos do tipo 1, e adicionando cada um na matriz Zbus
	for _, dados_linha := range elementos_tipo_1 {

		zbus_positiva = zbus_positiva.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_positiva)
		zbus_zero = zbus_zero.AdicionaElementoTipo1NaZbus(posicao, dados_linha.Z_zero)

		fmt.Println("Adicionado elemento tipo 1 -> Barra: "+dados_linha.De+"\t\tImpedancia:", dados_linha.Z_zero)

		barras_adicionadas[dados_linha.De] = Posicao_zbus{
			Posicao: posicao,
		}

		posicao++
	}

	// Adiciona os elementos do tipo 2
	// Valida se o elemento é do tipo 2, caso seja, adiciona na Zbus
	// Caso o elemento seja do tipo 3, ele adiciona em uma lista que será utilizada futuramente para adicionar os elementos tipo 3
	for len(elementos_tipo_2_3) != 0 {
		for nome_linha, linha := range elementos_tipo_2_3 {

			_, existe_de := barras_adicionadas[linha.De]
			_, existe_para := barras_adicionadas[linha.Para]

			if existe_de && existe_para {
				elementos_tipo_3 = append(elementos_tipo_3, linha)
				delete(elementos_tipo_2_3, nome_linha)

			} else if existe_de {
				zbus_positiva = zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_positiva)
				zbus_zero = zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.De].Posicao, posicao, linha.Z_zero)

				fmt.Println("Adicionado elemento tipo 2 -> Linha: "+linha.De+"-"+linha.Para+"\tImpedancia:", linha.Z_zero)

				barras_adicionadas[linha.Para] = Posicao_zbus{
					Posicao: posicao,
				}

				delete(elementos_tipo_2_3, nome_linha)
				posicao++

			} else if existe_para {
				zbus_positiva = zbus_positiva.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_positiva)
				zbus_zero = zbus_zero.AdicionaElementoTipo2NaZbus(barras_adicionadas[linha.Para].Posicao, posicao, linha.Z_zero)

				fmt.Println("Adicionado elemento tipo 2 -> Linha: "+linha.De+"-"+linha.Para+"\tImpedancia:", linha.Z_zero)

				barras_adicionadas[linha.De] = Posicao_zbus{
					Posicao: posicao,
				}

				delete(elementos_tipo_2_3, nome_linha)
				posicao++

			}
		}
	}

	// Com a lista criada de elementos do tipo 3, adicionamos na Zbus
	for x := 0; x < len(elementos_tipo_3); x++ {
		linha := elementos_tipo_3[x]

		fmt.Println("Adicionado elemento tipo 3 -> Linha: "+linha.De+"-"+linha.Para+" \tImpedancia:", linha.Z_zero, " \tRealizando redução de Kron")

		zbus_positiva = zbus_positiva.AdicionaElementoTipo3ComReducaoDeKron(
			barras_adicionadas[linha.De].Posicao,
			barras_adicionadas[linha.Para].Posicao,
			linha.Z_positiva,
			tamanho_do_sistema)
		zbus_zero = zbus_zero.AdicionaElementoTipo3ComReducaoDeKron(
			barras_adicionadas[linha.De].Posicao,
			barras_adicionadas[linha.Para].Posicao,
			linha.Z_zero,
			tamanho_do_sistema)
	}

	zbus := models.ZbusStr{
		Positiva: zbus_positiva.ArrayCmplxToArrayStr(),
		Negativa: zbus_positiva.ArrayCmplxToArrayStr(),
		Zero:     zbus_zero.ArrayCmplxToArrayStr(),
	}

	return zbus, barras_adicionadas
}
