package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

//var host string = os.Getenv("elements_host")
var host string = "http://localhost"
var port string = "8080"

//var port string = os.Getenv("elements_port")

func AllZbus(c *gin.Context) {

	ElementsSequencia, SystemSize, err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	zbus, _, err := GetZbus(ElementsSequencia, SystemSize)
	if err != nil {
		jsonError(c, err)
		return
	}

	c.JSON(http.StatusOK, zbus)
}

func ZbusSeq(c *gin.Context) {
	ElementsSequencia, SystemSize, err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	zbus, _, err := GetZbus(ElementsSequencia, SystemSize)
	if err != nil {
		jsonError(c, err)
		return
	}

	switch seq := c.Params.ByName("seq"); seq {
	case "positiva":
		c.JSON(http.StatusOK, zbus.Positiva)
	case "negativa":
		c.JSON(http.StatusOK, zbus.Negativa)
	case "zero":
		c.JSON(http.StatusOK, zbus.Zero)
	default:
		jsonError(c, errors.New("sequencia nao encontrada"))
	}
}

func Bars(c *gin.Context) {
	ElementsSequencia, SystemSize, err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	_, barrasAdicionadas, err := GetZbus(ElementsSequencia, SystemSize)
	if err != nil {
		jsonError(c, err)
		return
	}

	c.JSON(http.StatusOK, barrasAdicionadas)
}

func Bar(c *gin.Context) {
	bar := c.Params.ByName("bar")

	ElementsSequencia, SystemSize, err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	_, barrasAdicionadas, err := GetZbus(ElementsSequencia, SystemSize)
	if err != nil {
		jsonError(c, err)
		return
	}

	for barra := range barrasAdicionadas {
		if barra == bar {
			c.JSON(http.StatusOK, barrasAdicionadas[barra])
			return
		}
	}

	jsonError(c, errors.New("barra nao encontrada"))
}
