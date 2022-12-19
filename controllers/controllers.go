package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/luuisavelino/short-circuit-analysis-zbus/models"
	"github.com/luuisavelino/short-circuit-analysis-zbus/pkg/zbus"
)

//var url string= os.Getenv("URL")
//var port string = os.Getenv("PORT")

func Readness(w http.ResponseWriter, r *http.Request) {}

func Liveness(w http.ResponseWriter, r *http.Request) {}
func AllZbus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["file"]
	
	elementos_tipo_1 := GetAPI(file)["1"]
	elementos_tipo_2_3 := GetAPI(file)["2"]

	models.Zbus, _ = zbus.MontaZbus(elementos_tipo_1, elementos_tipo_2_3, 6)

	json.NewEncoder(w).Encode(models.Zbus)
}

func ZbusSeq(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["file"]

	elementos_tipo_1 := GetAPI(file)["1"]
	elementos_tipo_2_3 := GetAPI(file)["2"]

	models.Zbus, _ = zbus.MontaZbus(elementos_tipo_1, elementos_tipo_2_3, 6)

	switch seq := vars["seq"]; seq {
		case "positiva":
			json.NewEncoder(w).Encode(models.Zbus.Positiva)
		case "negativa":
			json.NewEncoder(w).Encode(models.Zbus.Negativa)
		case "zero":
			json.NewEncoder(w).Encode(models.Zbus.Zero)
		default:
			json.NewEncoder(w).Encode(models.Zbus)
	}
}

func GetAPI(file string) (map[string]map[string]models.Element) {
	url := "localhost"
	port := "8080"

	response, err := http.Get("http://" + url + ":" + port + "/api/files/" + file + "/elements")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}


	var responseObject = make(map[string]map[string]models.Element)
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}