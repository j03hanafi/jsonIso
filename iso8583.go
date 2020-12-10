package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mofax/iso8583"
	"github.com/rivo/uniseg"
	"io/ioutil"
	"net/http"
	"strconv"
)

func convertIso(transaction Transaction) iso8583.IsoStruct {

	cardAcceptor := transaction.CardAcceptorData.CardAcceptorName + transaction.CardAcceptorData.CardAcceptorCity + transaction.CardAcceptorData.CardAcceptorCountryCode

	one := iso8583.NewISOStruct("spec1987.yml", false)

	if one.Mti.String() != "" {
		fmt.Printf("Empty generates invalid MTI")
	}

	one.AddMTI("0200")
	one.AddField(2, transaction.Pan)
	one.AddField(3, transaction.ProcessingCode)
	one.AddField(4, strconv.Itoa(transaction.TotalAmount))
	one.AddField(5, transaction.SettlementAmount)
	one.AddField(6, transaction.CardholderBillingAmount)
	one.AddField(7, transaction.TransmissionDateTime)
	one.AddField(9, transaction.SettlementConversionRate)
	one.AddField(10, transaction.CardHolderBillingConvRate)
	one.AddField(11, transaction.Stan)
	one.AddField(12, transaction.LocalTransactionTime)
	one.AddField(13, transaction.LocalTransactionDate)
	one.AddField(17, transaction.CaptureDate)
	one.AddField(18, transaction.CategoryCode)
	one.AddField(22, transaction.PointOfServiceEntryMode)
	one.AddField(37, transaction.Refnum)
	one.AddField(41, transaction.CardAcceptorData.CardAcceptorTerminalId)
	one.AddField(43, cardAcceptor)
	one.AddField(48, transaction.AdditionalData)
	one.AddField(49, transaction.Currency)
	one.AddField(50, transaction.SettlementCurrencyCode)
	one.AddField(51, transaction.CardHolderBillingCurrencyCode)
	one.AddField(57, transaction.AdditionalDataNational)

	return one
}

func toIso8583(writer http.ResponseWriter, request *http.Request) {
	var response Iso8583
	var payment PaymentResponse
	err := pingDb(MysqlDB)
	processingCode := mux.Vars(request)["processingCode"]

	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 500, DatabaseError
		jsonFormatter(writer, response, 500)
	}

	transaction, statusCode, err := selectPayment(payment, processingCode, MysqlDB)

	isomsg := convertIso(transaction)

	msg, _ := isomsg.ToString()
	mti := isomsg.Mti.String()
	bitmap := isomsg.Bitmap
	header := fmt.Sprintf("%04d", uniseg.GraphemeClusterCount(msg))

	response.Iso8583 = header + msg
	response.MTI = mti
	response.Bitmap = bitmap
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = statusCode, "Success"

	jsonFormatter(writer, response, statusCode)
}

func extractElem(writer http.ResponseWriter, request *http.Request) {
	var response ExtractElem
	var payment PaymentResponse
	err := pingDb(MysqlDB)
	processingCode := mux.Vars(request)["processingCode"]
	element := mux.Vars(request)["element"]
	elementInt, _ := strconv.Atoi(element)
	elementInt64 := int64(elementInt)

	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 500, DatabaseError
		jsonFormatter(writer, response, 500)
	}

	transaction, statusCode, err := selectPayment(payment, processingCode, MysqlDB)

	//isoSpec, _ := iso8583.SpecFromFile("spec1987.yml")
	isomsg := convertIso(transaction)

	elementMap := isomsg.Elements.GetElements()
	data, exist := elementMap[elementInt64]

	var statusDesc string
	if exist {
		response.Data = data
		statusDesc = "Success"
	} else {
		response.Data = "Data not found"
		statusCode = 404
		statusDesc = "Data not found"
	}

	response.Element = elementInt64
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = statusCode, statusDesc

	jsonFormatter(writer, response, statusCode)
}

func jsonToIso(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var transaction Transaction
	var response Iso8583

	json.Unmarshal(reqBody, &transaction)

	iso := convertIso(transaction)

	msg, _ := iso.ToString()
	mti := iso.Mti.String()
	bmap := iso.Bitmap

	header := fmt.Sprintf("%04d", uniseg.GraphemeClusterCount(msg))

	response.Iso8583 = header + msg
	response.MTI = mti
	response.Bitmap = bmap
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 200, "Success"

	jsonFormatter(writer, response, 200)
}

func isoToJson(writer http.ResponseWriter, request *http.Request) {
	var response Json
	reqBody, _ := ioutil.ReadAll(request.Body)

	req := string(reqBody)

	header := req[:4]
	mti := req[4:8]
	hex := req[8:24]
	de := req[24:]

	bmap, _ := iso8583.HexToBitmapArray(hex)
	spec, _ := iso8583.SpecFromFile("spec1987.yml")

	element, _ := iso8583.UnpackElements(bmap, de, spec)
	q := iso8583.IsoStruct{
		Spec:     spec,
		Bitmap:   bmap,
		Elements: element,
	}
	emap := q.Elements.GetElements()

	response.Header, _ = strconv.Atoi(header)
	response.MTI = mti
	response.Bitmap = bmap

	response.Transaction.Pan = emap[2]
	response.Transaction.ProcessingCode = emap[3]
	response.Transaction.TotalAmount, _ = strconv.Atoi(emap[4])
	response.Transaction.SettlementAmount = emap[5]
	response.Transaction.CardholderBillingAmount = emap[6]
	response.Transaction.TransmissionDateTime = emap[7]
	response.Transaction.SettlementConversionRate = emap[9]
	response.Transaction.CardHolderBillingConvRate = emap[10]
	response.Transaction.Stan = emap[11]
	response.Transaction.LocalTransactionTime = emap[12]
	response.Transaction.LocalTransactionDate = emap[13]
	response.Transaction.CaptureDate = emap[17]
	response.Transaction.CategoryCode = emap[18]
	response.Transaction.PointOfServiceEntryMode = emap[22]
	response.Transaction.Refnum = emap[37]
	response.Transaction.CardAcceptorData.CardAcceptorTerminalId = emap[41]
	response.Transaction.AdditionalData = emap[48]
	response.Transaction.Currency = emap[49]
	response.Transaction.SettlementCurrencyCode = emap[50]
	response.Transaction.CardHolderBillingCurrencyCode = emap[51]
	response.Transaction.AdditionalDataNational = emap[57]

	jsonFormatter(writer, response, 200)

	//header mti msg
}
