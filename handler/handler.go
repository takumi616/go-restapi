package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"go-restapi/entity"
	"go-restapi/repository"
)

//Get all stock index data
func GetAllDateStockIndex(w http.ResponseWriter, r *http.Request) {
	allDateStockIndex, err := repository.SelectAllDateStockIndex(r.Context())
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error happened in repository package: %v", err)))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(allDateStockIndex)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to encode response data: %v", err)))
		return 	
	}
}

//get specific date stock index data
func GetStockIndex(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	
	stockIndex, err := repository.SelectStockIndex(r.Context(), param["date"])
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error happened in repository package: %v", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stockIndex)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to encode response data: %v", err)))
		return
	}
}

//create new stock index data
func CreateStockIndex(w http.ResponseWriter, r *http.Request) {
	var stockIndex entity.StockIndex
	err := json.NewDecoder(r.Body).Decode(&stockIndex)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Failed to decode request data: %v", err)))
		return
	}

	rowsAffectedNumber, err := repository.InsertStockIndex(r.Context(), &stockIndex)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error happened in repository package: %v", err)))
		return 
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Successfully insert data and inserted record number is %d.", rowsAffectedNumber)))
}

//update specific date stock index data
func UpdateStockIndex(w http.ResponseWriter, r *http.Request) {
	var stockIndex entity.StockIndex
	err := json.NewDecoder(r.Body).Decode(&stockIndex)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error message: %v", err.Error())))
		return 
	}

	param := mux.Vars(r)
	rowsAffectedNumber, err := repository.UpdateStockIndex(r.Context(), &stockIndex, param["date"])
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error happened in repository package: %v", err)))
		return 
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Successfully update data and updated record number is %d.", rowsAffectedNumber)))
}

//delete specific date stock index data 
func DeleteStockIndex(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	
	rowsAffectedNumber, err := repository.DeleteStockIndex(r.Context(), param["date"])
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error happened in repository package: %v", err)))
		return
	}
	
	w.Write([]byte(fmt.Sprintf("Successfully delete data and deleted record number is %v.", rowsAffectedNumber)))
}