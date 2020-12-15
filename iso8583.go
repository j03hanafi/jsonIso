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
	"strings"
	"time"
)

func convertIso(transaction Transaction) iso8583.IsoStruct {

	cardAcceptorTerminalId := transaction.CardAcceptorData.CardAcceptorTerminalId
	cardAcceptorName := transaction.CardAcceptorData.CardAcceptorName
	cardAcceptorCity := transaction.CardAcceptorData.CardAcceptorCity
	cardAcceptorCountryCode := transaction.CardAcceptorData.CardAcceptorCountryCode

	if len(transaction.CardAcceptorData.CardAcceptorTerminalId) < 16 {
		cardAcceptorTerminalId = rightPad(transaction.CardAcceptorData.CardAcceptorTerminalId, 16, " ")
	}
	if len(transaction.CardAcceptorData.CardAcceptorName) < 25 {
		cardAcceptorName = rightPad(transaction.CardAcceptorData.CardAcceptorName, 25, " ")
	}
	if len(transaction.CardAcceptorData.CardAcceptorCity) < 13 {
		cardAcceptorCity = rightPad(transaction.CardAcceptorData.CardAcceptorCity, 13, " ")
	}
	if len(transaction.CardAcceptorData.CardAcceptorCountryCode) < 2 {
		cardAcceptorCountryCode = rightPad(transaction.CardAcceptorData.CardAcceptorCountryCode, 2, " ")
	}
	cardAcceptor := cardAcceptorName + cardAcceptorCity + cardAcceptorCountryCode

	trans := map[int64]string{2: transaction.Pan,
		3:  transaction.ProcessingCode,
		4:  strconv.Itoa(transaction.TotalAmount),
		5:  transaction.SettlementAmount,
		6:  transaction.CardholderBillingAmount,
		7:  transaction.TransmissionDateTime,
		9:  transaction.SettlementConversionRate,
		10: transaction.CardHolderBillingConvRate,
		11: transaction.Stan,
		12: transaction.LocalTransactionTime,
		13: transaction.LocalTransactionDate,
		17: transaction.CaptureDate,
		18: transaction.CategoryCode,
		22: transaction.PointOfServiceEntryMode,
		37: transaction.Refnum,
		41: cardAcceptorTerminalId,
		43: cardAcceptor,
		48: transaction.AdditionalData,
		49: transaction.Currency,
		50: transaction.SettlementCurrencyCode,
		51: transaction.CardHolderBillingCurrencyCode,
		57: transaction.AdditionalDataNational,
	}

	one := iso8583.NewISOStruct("spec1987.yml", false)
	spec, _ := specFromFile("spec1987.yml")

	if one.Mti.String() != "" {
		fmt.Printf("Empty generates invalid MTI")
	}

	one.AddMTI("0200")

	for field, data := range trans {

		fieldSpec := spec.fields[int(field)]

		if fieldSpec.LenType == "fixed" {
			lengthValidate, _ := iso8583.FixedLengthIntegerValidator(int(field), fieldSpec.MaxLen, data)

			if lengthValidate == false {
				if fieldSpec.ContentType == "n" {
					data = leftPad(data, fieldSpec.MaxLen, "0")
				} else {
					data = rightPad(data, fieldSpec.MaxLen, " ")
				}
			}
		}

		one.AddField(field, data)

		writeLog("Field[" + strconv.Itoa(int(field)) + "]: " + data + " (" + fieldSpec.Label + ")")

	}

	return one
}

func convertJson(message string) Json {
	var response Json

	header := message[0:4]
	data := message[4:]

	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)

	msg, err := isoStruct.Parse(data)
	if err != nil {
		fmt.Println(err)
	}

	response.Header, _ = strconv.Atoi(header)
	response.MTI = msg.Mti.String()
	response.Hex, _ = iso8583.BitMapArrayToHex(msg.Bitmap)

	emap := msg.Elements.GetElements()

	cardAcceptorTerminalId := strings.TrimRight(msg.Elements.GetElements()[41], " ")
	cardAcceptorName := strings.TrimRight(msg.Elements.GetElements()[43][:25], " ")
	cardAcceptorCity := strings.TrimRight(msg.Elements.GetElements()[43][25:38], " ")
	cardAcceptorCountryCode := strings.TrimRight(msg.Elements.GetElements()[43][38:], " ")

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
	response.Transaction.CardAcceptorData.CardAcceptorTerminalId = cardAcceptorTerminalId
	response.Transaction.CardAcceptorData.CardAcceptorName = cardAcceptorName
	response.Transaction.CardAcceptorData.CardAcceptorCity = cardAcceptorCity
	response.Transaction.CardAcceptorData.CardAcceptorCountryCode = cardAcceptorCountryCode
	response.Transaction.AdditionalData = emap[48]
	response.Transaction.Currency = emap[49]
	response.Transaction.SettlementCurrencyCode = emap[50]
	response.Transaction.CardHolderBillingCurrencyCode = emap[51]
	response.Transaction.AdditionalDataNational = emap[57]

	return response
}

func toIso8583(writer http.ResponseWriter, request *http.Request) {
	var response Iso8583
	var payment PaymentResponse
	err := pingDb(MysqlDB)
	processingCode := mux.Vars(request)["processingCode"]

	writeLog("New Request: JSON to ISO8583")

	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 500, DatabaseError
		jsonFormatter(writer, response, 500)
	}

	transaction, statusCode, err := selectPayment(payment, processingCode, MysqlDB)
	writeLog("Get transaction from database, processing code: " + transaction.ProcessingCode)

	isomsg := convertIso(transaction)

	msg, _ := isomsg.ToString()
	mti := isomsg.Mti.String()
	hex, _ := iso8583.BitMapArrayToHex(isomsg.Bitmap)
	header := fmt.Sprintf("%04d", uniseg.GraphemeClusterCount(msg))

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02_15.04.05")
	filename := CreateFile(transaction.ProcessingCode+"_"+date, header+msg)

	response.Iso8583 = header + msg
	response.MTI = mti
	response.Hex = hex
	response.FileName = filename
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = statusCode, "Success"

	writeLog("Header: " + header)
	writeLog("MTI: " + mti)
	writeLog("HEX: " + hex)
	writeLog("ISO Message: " + msg)
	writeLog("ISO Full Message: " + header + msg)
	writeLog("File: " + filename)
	writeLog("JSON to ISO8583 success (200)")

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
	writeLog("New Request: Extract Element ISO8583")

	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 500, DatabaseError
		jsonFormatter(writer, response, 500)
	}

	transaction, statusCode, err := selectPayment(payment, processingCode, MysqlDB)

	//isoSpec, _ := iso8583.SpecFromFile("spec1987.yml")
	isomsg := convertIso(transaction)
	writeLog("Get transaction from database, processing code: " + transaction.ProcessingCode)

	elementMap := isomsg.Elements.GetElements()
	data, exist := elementMap[elementInt64]

	var statusDesc string
	if exist {
		response.Data = data
		statusDesc = "Success"

		writeLog("Field[" + strconv.Itoa(elementInt) + "]: " + data)
		writeLog("Extract Element ISO8583 success (200)")
	} else {
		response.Data = "Data not found"
		statusCode = 404
		statusDesc = "Data not found"

		writeLog("Field[" + strconv.Itoa(elementInt) + "]: Data not found")
		writeLog("Extract Element ISO8583 failed (404)")
	}

	response.Element = elementInt64
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = statusCode, statusDesc

	jsonFormatter(writer, response, statusCode)
}

func jsonToIso(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var transaction Transaction
	var response Iso8583

	writeLog("New Request: JSON to ISO8583")

	json.Unmarshal(reqBody, &transaction)

	iso := convertIso(transaction)

	msg, _ := iso.ToString()
	mti := iso.Mti.String()
	hex, _ := iso8583.BitMapArrayToHex(iso.Bitmap)

	header := fmt.Sprintf("%04d", uniseg.GraphemeClusterCount(msg))

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02_15.04.05")
	filename := CreateFile(transaction.ProcessingCode+"_"+date, header+msg)

	response.Iso8583 = header + msg
	response.MTI = mti
	response.Hex = hex
	response.FileName = filename
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 200, "Success"

	writeLog("Header: " + header)
	writeLog("MTI: " + mti)
	writeLog("HEX: " + hex)
	writeLog("ISO Message: " + msg)
	writeLog("ISO Full Message: " + header + msg)
	writeLog("File: " + filename)
	writeLog("JSON to ISO8583 success (200)")

	jsonFormatter(writer, response, 200)
}

func isoToJson(writer http.ResponseWriter, request *http.Request) {
	var response Json
	reqBody, _ := ioutil.ReadAll(request.Body)

	writeLog("New Request: ISO8583 to JSON from file")

	req := string(reqBody)
	writeLog("ISO Message: " + req)

	response = convertJson(req)

	writeLog("ISO8583 to JSON success (200)")
	response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 200, "Success"
	jsonFormatter(writer, response, 200)
}

func isoToJsonFile(writer http.ResponseWriter, request *http.Request) {
	var response Json
	reqBody, _ := ioutil.ReadAll(request.Body)

	writeLog("New Request: ISO8583 to JSON")

	fileName := string(reqBody)
	fileName = "isoFiles/" + fileName

	if FuncCheckExist(fileName) {
		content := ReadFile(fileName)
		writeLog("ISO Message: " + content)
		response = convertJson(content)

		writeLog("ISO8583 to JSON success (200)")
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 200, "Success"

	} else {
		writeLog("File not found")
		writeLog("ISO8583 to JSON failed (404)")
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 404, "File Not Found"
	}

	jsonFormatter(writer, response, 200)

}
