package main

import (
	"github.com/gorilla/mux"
)

func server() *mux.Router {
	router := mux.NewRouter()

	//endpoints
	//CRUD
	router.HandleFunc("/", initHandler)
	router.HandleFunc("/payment", getPayments).Methods("GET")
	router.HandleFunc("/payment/{processingCode}", getPayment).Methods("GET")
	router.HandleFunc("/payment/", createPayment).Methods("POST")
	router.HandleFunc("/payment/{processingCode}", updatePayment).Methods("PUT")
	router.HandleFunc("/payment/{processingCode}", deletePayment).Methods("DELETE")

	//iso8583
	router.HandleFunc("/payment/{processingCode}/iso", toIso8583).Methods("GET")
	router.HandleFunc("/payment/{processingCode}/iso/{element}", extractElem).Methods("GET")
	router.HandleFunc("/payment/iso", jsonToIso).Methods("POST")
	router.HandleFunc("/payment/json", isoToJson).Methods("POST")

	return router
}
