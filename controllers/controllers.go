package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
)

var host string = os.Getenv("elements_host")
var port string = os.Getenv("elements_port")

func AllZbus(c *gin.Context) {
	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Zbus)
}

func ZbusSeq(c *gin.Context) {
	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
		return
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
		return
	}

	line := c.Params.ByName("line")
	deleteLine(line)

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Zbus)
}

func ShortCircuit(c *gin.Context) {
	line := c.Params.ByName("line")
	point, err := strconv.Atoi(c.Params.ByName("point"))
	if err != nil {
		jsonError(c, err)
		return
	}

	err = GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	err = adicionaBarraFicticia(line, point)
	if err != nil {
		jsonError(c, err)
		return
	}

	deleteLine(line)

	err = GetZbus()
	if err != nil {
		jsonError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Zbus)
}
