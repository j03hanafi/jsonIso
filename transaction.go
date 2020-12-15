package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func initHandler(writer http.ResponseWriter, request *http.Request) {
	jsonFormatter(writer, map[string]interface{}{
		"status":  200,
		"message": "Hello World!",
	}, 200)
}

func getPayments(writer http.ResponseWriter, request *http.Request) {
	var response PaymentsResponse
	err := pingDb(MysqlDB)

	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 500, DatabaseError
		jsonFormatter(writer, response, 500)
	}

	payments, statusCode, err := selectPayments(response, MysqlDB)
	response.TransactionData = payments
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = statusCode, "Success"

	jsonFormatter(writer, response, statusCode)

}

func getPayment(writer http.ResponseWriter, request *http.Request) {
	var response PaymentResponse
	err := pingDb(MysqlDB)
	processingCode := mux.Vars(request)["processingCode"]

	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 500, DatabaseError
		jsonFormatter(writer, response, 500)
	}

	payment, statusCode, err := selectPayment(response, processingCode, MysqlDB)
	response.TransactionData = payment
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = statusCode, "Success"

	jsonFormatter(writer, response, statusCode)
}

func createPayment(writer http.ResponseWriter, request *http.Request) {
	/*
		var response PaymentResponse
		var transaction TransactionResponse
		reqBody, _ := ioutil.ReadAll(request.Body)

		err := pingDb(MysqlDB)
		if err != nil {
			response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 500, DatabaseError
			jsonFormatter(writer, response, 500)
		} else {
			var transaction Transaction
			err := json.Unmarshal(reqBody, &transaction)

			if err != nil {
				fmt.Println(err)
			}


		}

	*/

}

func updatePayment(writer http.ResponseWriter, request *http.Request) {

}

func deletePayment(writer http.ResponseWriter, request *http.Request) {

}
