package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
	"github.com/luuisavelino/short-circuit-analysis-zbus/pkg/zbus"
)

func adicionaBarraFicticia(line string, point int) error {

	nums := strings.Split(line, "-")
	de, para := nums[0], nums[1]

	z_pos, err := strconv.ParseComplex(models.Elements["2"][line].Z_positiva, 64)
	if err != nil {
		return err
	}

	z_zero, err := strconv.ParseComplex(models.Elements["2"][line].Z_zero, 64)
	if err != nil {
		return err
	}

	models.Elements["2"][de+"-ficticia"] = models.Element{
		De:         de,
		Para:       "ficticia",
		Nome:       "Bara ficticia",
		Z_positiva: fmt.Sprint(z_pos * complex(float64(point), 0) / 100),
		Z_zero:     fmt.Sprint(z_zero * complex(float64(point), 0) / 100),
	}

	models.Elements["2"]["ficticia-"+para] = models.Element{
		De:         "ficticia",
		Para:       para,
		Nome:       "Barra ficticia",
		Z_positiva: fmt.Sprint(z_pos * complex(float64(100-point), 0) / 100),
		Z_zero:     fmt.Sprint(z_zero * complex(float64(100-point), 0) / 100),
	}

	models.SystemSize["size"]++

	return nil
}

func deleteLine(line string) {
	for tipo, element := range models.Elements {
		delete(element, line)
		models.Elements[tipo] = element
	}
}

func GetAPI(endpoint string) ([]byte, error) {
	response, err := http.Get(host + ":" + port + endpoint)

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func GetData(c *gin.Context) error {
	fileId := c.Params.ByName("fileId")

	responseData, err := GetAPI("/api/v2/files/" + fileId + "/types/0/elements")
	json.Unmarshal(responseData, &models.Elements)

	if err != nil {
		return err
	}

	responseData, err = GetAPI("/api/v2/files/" + fileId + "/size")
	json.Unmarshal(responseData, &models.SystemSize)

	if err != nil {
		return err
	}

	return nil
}

func jsonError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Error": err.Error(),
	})
}

func GetElements(c *gin.Context) error {
	line := c.Params.ByName("line")

	err := GetData(c)
	if err != nil {
		return err
	}

	if c.Params.ByName("point") != "" {
		point, _ := strconv.Atoi(c.Params.ByName("point"))
		if err != nil {
			return err
		}

		err = adicionaBarraFicticia(line, point)
		if err != nil {
			return err
		}
	}

	if line != "" {
		deleteLine(line)
	}

	return nil
}

func GetZbus() (map[string]zbus.Posicao_zbus, error) {
	barrasAdicionadas, err := zbus.MontaZbus()
	return barrasAdicionadas, err
}
