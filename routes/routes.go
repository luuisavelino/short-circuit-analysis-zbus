package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/luuisavelino/short-circuit-analysis-zbus/controllers"
)

func HandleRequest() {
	r := mux.NewRouter()
	r.HandleFunc("/health/liveness", controllers.Liveness).Methods("Get")
	r.HandleFunc("/health/readiness", controllers.Readness).Methods("Get")
	r.HandleFunc("/api/files/{file}/zbus", controllers.AllZbus).Methods("Get")
	r.HandleFunc("/api/files/{file}/zbus/{seq}", controllers.ZbusSeq).Methods("Get")
	log.Fatal(http.ListenAndServe(":8081", r))
}
