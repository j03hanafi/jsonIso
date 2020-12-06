package main

const (
	GetPaymentsQuery = `
		SELECT
			pan,
			processingCode,
			totalAmount,
			acquirerId,
			issuerId,
			transmissionDateTime,
			localTransactionTime,
			localTransactionDate,
			captureDate,
			additionalData,
			stan,
			refnum,
			currency,
			terminalId,
			accountFrom,
			accountTo,
			categoryCode,
			settlementAmount,
			cardholderBillingAmount,
			settlementConversionRate,
			cardHolderBillingConvRate,
			pointOfServiceEntryMode,
			settlementCurrencyCode,
			cardHolderBillingCurrencyCode,
			additionalDataNational
		FROM
			transaction
	`

	GetCardAcceptorQuery = `
		SELECT
			cardAcceptorTerminalID,
			cardAcceptorName,
			cardAcceptorCity,
			cardAcceptorCountryCode
		FROM
			transaction
		WHERE processingCode = ?
	`

	GetPaymentQuery = `
		SELECT
			pan,
			processingCode,
			totalAmount,
			acquirerId,
			issuerId,
			transmissionDateTime,
			localTransactionTime,
			localTransactionDate,
			captureDate,
			additionalData,
			stan,
			refnum,
			currency,
			terminalId,
			accountFrom,
			accountTo,
			categoryCode,
			settlementAmount,
			cardholderBillingAmount,
			settlementConversionRate,
			cardHolderBillingConvRate,
			pointOfServiceEntryMode,
			settlementCurrencyCode,
			cardHolderBillingCurrencyCode,
			additionalDataNational
		FROM
			transaction
		WHERE processingCode = ?
	`
)
