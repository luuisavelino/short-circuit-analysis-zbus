package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
)

//var host string = os.Getenv("elements_host")
var host string = "http://localhost"
var port string = "8080"

//var port string = os.Getenv("elements_port")

func AllZbus(c *gin.Context) {

	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}
	_, err = GetZbus()
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
	_, err = GetZbus()
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

func Bars(c *gin.Context) {
	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	barrasAdicionadas, err := GetZbus()
	if err != nil {
		jsonError(c, err)
		return
	}

	c.JSON(http.StatusOK, barrasAdicionadas)
}


func Bar(c *gin.Context) {
	bar := c.Params.ByName("bar")

	err := GetElements(c)
	if err != nil {
		jsonError(c, err)
		return
	}

	barrasAdicionadas, err := GetZbus()
	if err != nil {
		jsonError(c, err)
		return
	}

	for barra, _ := range barrasAdicionadas {
		if barra == bar {
			c.JSON(http.StatusOK, barrasAdicionadas[barra])
			return
		}
	}

	jsonError(c, errors.New("barra nao encontrada"))
}
