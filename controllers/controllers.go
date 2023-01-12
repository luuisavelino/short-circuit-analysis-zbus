package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
	"github.com/luuisavelino/short-circuit-analysis-zbus/pkg/zbus"
)

var host string = os.Getenv("elements_host")
var port string = os.Getenv("elements_port")

func AllZbus(c *gin.Context) {
	fileId := c.Params.ByName("fileId")

	responseData, err := GetAPI("api/v2/files/" + fileId + "/types/0/elements")
	json.Unmarshal(responseData, &models.Elements)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	responseData, err = GetAPI("api/v2/files/" + fileId + "/size")
	json.Unmarshal(responseData, &models.SystemSize)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	models.Zbus, _ = zbus.MontaZbus()
	c.JSON(http.StatusOK, models.Zbus)
}

func ZbusSeq(c *gin.Context) {
	fileId := c.Params.ByName("fileId")

	responseData, err := GetAPI("api/v2/files/" + fileId + "/types/0/elements")
	json.Unmarshal(responseData, &models.Elements)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	responseData, err = GetAPI("api/v2/files/" + fileId + "/size")
	json.Unmarshal(responseData, &models.SystemSize)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	models.Zbus, _ = zbus.MontaZbus()

	switch seq := c.Params.ByName("seq"); seq {
	case "positiva":
		c.JSON(http.StatusOK, models.Zbus.Positiva)
	case "negativa":
		c.JSON(http.StatusOK, models.Zbus.Negativa)
	case "zero":
		c.JSON(http.StatusOK, models.Zbus.Zero)
	default:
		c.JSON(http.StatusOK, models.Zbus)
	}
}

func GetAPI(endpoint string) ([]byte, error) {
	//response, err := http.Get(host + ":" + port + "/api/v2/files/" + fileId + "/types/0/elements")
	response, err := http.Get(host + ":" + port + endpoint)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return responseData, nil
}
