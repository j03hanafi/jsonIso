package main

type CardAcceptorData struct {
	CardAcceptorTerminalId  string `json:"cardAcceptorTerminalID"`
	CardAcceptorName        string `json:"cardAcceptorName"`
	CardAcceptorCity        string `json:"cardAcceptorCity"`
	CardAcceptorCountryCode string `json:"cardAcceptorCountryCode"`
}

type Transaction struct {
	Pan                           string           `json:"pan"`
	ProcessingCode                string           `json:"processingCode"`
	TotalAmount                   int              `json:"totalAmount"`
	AcquirerId                    string           `json:"acquirerId"`
	IssuerId                      string           `json:"issuerId"`
	TransmissionDateTime          string           `json:"transmissionDateTime"`
	LocalTransactionTime          string           `json:"localTransactionTime"`
	LocalTransactionDate          string           `json:"localTransactionDate"`
	CaptureDate                   string           `json:"captureDate"`
	AdditionalData                string           `json:"additionalData"`
	Stan                          string           `json:"stan"`
	Refnum                        string           `json:"refnum"`
	Currency                      string           `json:"currency"`
	TerminalId                    string           `json:"terminalId"`
	AccountFrom                   string           `json:"accountFrom"`
	AccountTo                     string           `json:"accountTo"`
	CategoryCode                  string           `json:"categoryCode"`
	SettlementAmount              string           `json:"settlementAmount"`
	CardholderBillingAmount       string           `json:"cardholderBillingAmount"`
	SettlementConversionRate      string           `json:"settlementConversionRate"`
	CardHolderBillingConvRate     string           `json:"cardHolderBillingConvRate"`
	PointOfServiceEntryMode       string           `json:"pointOfServiceEntryMode"`
	CardAcceptorData              CardAcceptorData `json:"cardAcceptorData"`
	SettlementCurrencyCode        string           `json:"settlementCurrencyCode"`
	CardHolderBillingCurrencyCode string           `json:"cardHolderBillingCurrencyCode"`
	AdditionalDataNational        string           `json:"additionalDataNational"`
}

type Response struct {
	ResponseCode        int    `json:"responseCode"`
	ReasonCode          int    `json:"reasonCode"`
	ResponseDescription string `json:"responseDescription"`
}

type PaymentsResponse struct {
	TransactionData []Transaction `json:"transactionData"`
	ResponseStatus  Response      `json:"responseStatus"`
}

type PaymentResponse struct {
	TransactionData Transaction `json:"transactionData"`
	ResponseStatus  Response    `json:"responseStatus"`
}

type Iso8583 struct {
	Iso8583        string   `json:"iso_8583"`
	MTI            string   `json:"mti"`
	Hex            string   `json:"hex"`
	ResponseStatus Response `json:"responseStatus"`
}

type ExtractElem struct {
	Element        int64    `json:"element"`
	Data           string   `json:"data"`
	Header         string   `json:"header"`
	ResponseStatus Response `json:"responseStatus"`
}

type LogIso struct {
	ResponseStatus Response `json:"responseStatus"`
}