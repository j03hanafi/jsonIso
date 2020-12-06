package main

import (
	"database/sql"
	"log"
)

func selectPayments(response PaymentsResponse, db *sql.DB) ([]Transaction, int, error) {
	rows, err := db.Query(GetPaymentsQuery)
	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 404, "Transaction Not Found"
		return nil, 404, err
	}

	var payments []Transaction

	for rows.Next() {
		var payment Transaction
		if err := rows.Scan(
			&payment.Pan,
			&payment.ProcessingCode,
			&payment.TotalAmount,
			&payment.AcquirerId,
			&payment.IssuerId,
			&payment.TransmissionDateTime,
			&payment.LocalTransactionTime,
			&payment.LocalTransactionDate,
			&payment.CaptureDate,
			&payment.AdditionalData,
			&payment.Stan,
			&payment.Refnum,
			&payment.Currency,
			&payment.TerminalId,
			&payment.AccountFrom,
			&payment.AccountTo,
			&payment.CategoryCode,
			&payment.SettlementAmount,
			&payment.CardholderBillingAmount,
			&payment.SettlementConversionRate,
			&payment.CardHolderBillingConvRate,
			&payment.PointOfServiceEntryMode,
			&payment.SettlementCurrencyCode,
			&payment.CardHolderBillingCurrencyCode,
			&payment.AdditionalDataNational,
		); err != nil {
			log.Print(err)
		} else {
			payment.CardAcceptorData = selectCardAcceptor(payment.ProcessingCode, db)
			payments = append(payments, payment)
		}
	}

	return payments, 200, err
}

func selectPayment(response PaymentResponse, processingCode string, db *sql.DB) (Transaction, int, error) {
	var payment Transaction
	rows, err := db.Query(GetPaymentQuery, processingCode)
	if err != nil {
		response.ResponseStatus.ReasonCode, response.ResponseStatus.ResponseDescription = 404, "Transaction Not Found"
		return Transaction{}, 0, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&payment.Pan,
			&payment.ProcessingCode,
			&payment.TotalAmount,
			&payment.AcquirerId,
			&payment.IssuerId,
			&payment.TransmissionDateTime,
			&payment.LocalTransactionTime,
			&payment.LocalTransactionDate,
			&payment.CaptureDate,
			&payment.AdditionalData,
			&payment.Stan,
			&payment.Refnum,
			&payment.Currency,
			&payment.TerminalId,
			&payment.AccountFrom,
			&payment.AccountTo,
			&payment.CategoryCode,
			&payment.SettlementAmount,
			&payment.CardholderBillingAmount,
			&payment.SettlementConversionRate,
			&payment.CardHolderBillingConvRate,
			&payment.PointOfServiceEntryMode,
			&payment.SettlementCurrencyCode,
			&payment.CardHolderBillingCurrencyCode,
			&payment.AdditionalDataNational,
		); err != nil {
			log.Print(err)
		}
		payment.CardAcceptorData = selectCardAcceptor(payment.ProcessingCode, db)
	}
	return payment, 200, err
}

func selectCardAcceptor(processingCode string, db *sql.DB) CardAcceptorData {
	var cardAcceptor CardAcceptorData
	rows, err := db.Query(GetCardAcceptorQuery, processingCode)
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(
			&cardAcceptor.CardAcceptorTerminalId,
			&cardAcceptor.CardAcceptorName,
			&cardAcceptor.CardAcceptorCity,
			&cardAcceptor.CardAcceptorCountryCode,
		); err != nil {
			log.Print(err)
		}
	}

	return cardAcceptor
}
