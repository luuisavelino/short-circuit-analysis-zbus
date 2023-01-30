package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
	"github.com/luuisavelino/short-circuit-analysis-zbus/pkg/zbus"
)

func adicionaBarraFicticia(line string, point int, elements models.Elements) (models.Elements, error) {

	nums := strings.Split(line, "-")
	de, para := nums[0], nums[1]

	z_pos, err := strconv.ParseComplex(elements[line].Z_positiva, 64)
	if err != nil {
		return nil, err
	}

	z_zero, err := strconv.ParseComplex(elements[line].Z_zero, 64)
	if err != nil {
		return nil, err
	}

	elements[de+"-ficticia"] = models.Element{
		De:         de,
		Para:       "ficticia",
		Nome:       "Bara ficticia",
		Z_positiva: fmt.Sprint(z_pos * complex(float64(point), 0) / 100),
		Z_zero:     fmt.Sprint(z_zero * complex(float64(point), 0) / 100),
	}

	elements["ficticia-"+para] = models.Element{
		De:         "ficticia",
		Para:       para,
		Nome:       "Barra ficticia",
		Z_positiva: fmt.Sprint(z_pos * complex(float64(100-point), 0) / 100),
		Z_zero:     fmt.Sprint(z_zero * complex(float64(100-point), 0) / 100),
	}

	atomic.AddUint64(&models.WriteOps, 1)

	return elements, nil
}

func deleteLine(line string, elements models.Elements) models.Elements {
	delete(elements, line)
	return elements
}

func GetData(c *gin.Context) (map[string]models.Elements, map[string]int, error) {
	fileId := c.Params.ByName("fileId")

	var ElementsSequencia = make(map[string]models.Elements)
	var SystemSize = make(map[string]int)
	var err error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		url := host + ":" + port + "/api/v2/files/" + fileId + "/types/0/elements"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Erro ao Unmarshal endpoint ElementsSequencia:", err)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&ElementsSequencia); err != nil {
			fmt.Println("Erro ao Unmarshal endpoint ElementsSequencia", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		url := host + ":" + port + "/api/v2/files/" + fileId + "/size"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Erro ao Unmarshal endpoint SystemSize:", err)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&SystemSize); err != nil {
			fmt.Println("Erro ao Unmarshal endpoint SystemSize:", err)
			return
		}
	}()

	fmt.Println(err)

	if err != nil {
		fmt.Println(err)
		fmt.Println("ERRRROOO")
		return nil, nil, err
	}

	wg.Wait()

	fmt.Println(ElementsSequencia, SystemSize)

	return ElementsSequencia, SystemSize, nil
}

func jsonError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Error": err.Error(),
	})
}

func GetElements(c *gin.Context) (map[string]models.Elements, map[string]int, error) {
	line := c.Params.ByName("line")

	ElementsSequencia, SystemSize, err := GetData(c)

	if err != nil {
		return nil, nil, err
	}

	if c.Params.ByName("point") != "" {
		point, _ := strconv.Atoi(c.Params.ByName("point"))
		if err != nil {
			return nil, nil, err
		}

		ElementsSequencia["2"], err = adicionaBarraFicticia(line, point, ElementsSequencia["2"])
		if err != nil {
			return nil, nil, err
		}

		SystemSize["size"]++
	}

	if line != "" {
		for sequencia, elements := range ElementsSequencia {
			ElementsSequencia[sequencia] = deleteLine(line, elements)
		}
	}

	return ElementsSequencia, SystemSize, nil
}

func GetZbus(ElementsSequencia map[string]models.Elements, SystemSize map[string]int) (models.ZbusStr, map[string]zbus.Posicao_zbus, error) {
	zbus, barrasAdicionadas, err := zbus.MontaZbus(ElementsSequencia, SystemSize)
	return zbus, barrasAdicionadas, err
}
