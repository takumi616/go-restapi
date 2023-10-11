package routing

import (
	"github.com/gorilla/mux"
	"go-restapi/handler"
)

func RegistHandler() *mux.Router {

	router := mux.NewRouter()

	//regist handler to router
	router.HandleFunc("/nasdaq100", handler.GetAllDateStockIndex).Methods("GET")
	router.HandleFunc("/nasdaq100/{date}", handler.GetStockIndex).Methods("GET")
	router.HandleFunc("/nasdaq100", handler.CreateStockIndex).Methods("POST")
	router.HandleFunc("/nasdaq100/{date}", handler.UpdateStockIndex).Methods("PUT")
	router.HandleFunc("/nasdaq100/{date}", handler.DeleteStockIndex).Methods("DELETE")
	
	return router
}