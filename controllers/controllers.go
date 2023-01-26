package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
	"github.com/luuisavelino/short-circuit-analysis-zbus/pkg/zbus"
)

var host string = os.Getenv("elements_host")
var port string = os.Getenv("elements_port")

func AllZbus(c *gin.Context) {
	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
	}

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
	}

	c.JSON(http.StatusOK, models.Zbus)
}

func ZbusSeq(c *gin.Context) {
	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
	}

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
	}

	switch seq := c.Params.ByName("seq"); seq {
	case "positiva":
		c.JSON(http.StatusOK, models.Zbus.Positiva)
	case "negativa":
		c.JSON(http.StatusOK, models.Zbus.Negativa)
	case "zero":
		c.JSON(http.StatusOK, models.Zbus.Zero)
	default:
		jsonError(c, errors.New("sequencia nao encontrada"))
	}
}

func Atuacao(c *gin.Context) {

	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
	}

	line := c.Params.ByName("line")
	deleteLine(line)

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
	}

	c.JSON(http.StatusOK, models.Zbus)
}

func ShortCircuit(c *gin.Context) {
	line := c.Params.ByName("line")
	point, err := strconv.Atoi(c.Params.ByName("point"))
	if err != nil {
		jsonError(c, err)
	}

	err = GetElements(c)
	if err != nil {
		jsonError(c, err)
	}

	err = adicionaBarraFicticia(line, point)
	if err != nil {
		jsonError(c, err)
	}

	deleteLine(line)

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
	}

	c.JSON(http.StatusOK, models.Zbus)
}

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
		Z_zero:     fmt.Sprint(z_zero * complex(float64(point), 0)/ 100),
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

func GetElements(c *gin.Context) error {
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

func GetZbus() error {
	_, err := zbus.MontaZbus()

	return err
}

func jsonError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Error": err.Error(),
	})
}
