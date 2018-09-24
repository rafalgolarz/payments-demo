package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

var jsonPOST = []byte(`
	{
		"type": "Payment",
		"id": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
		"version": 0,
		"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		  "amount": "100.21",
		  "beneficiary_party": {
			"account_name": "W Owens",
			"account_number": "31926819",
			"account_number_code": "BBAN",
			"address": "1 The Beneficiary Localtown SE2",
			"bank_id": "403000",
			"bank_id_code": "GBDSC",
			"name": "Wilfred Jeremiah Owens"
		  },
		  "charges_information": {
			"bearer_code": "SHAR",
			"sender_charges_amount": "5.00",
			 "sender_charges_currency": "GBP",
			"receiver_charges_amount": "1.00",
			"receiver_charges_currency": "USD"
		  },
		  "currency": "GBP",
		  "debtor_party": {
			"account_name": "EJ Brown Black",
			"account_number": "GB29XABC10161234567801",
			"account_number_code": "IBAN",
			"address": "10 Debtor Crescent Sourcetown NE1",
			"bank_id": "203301",
			"bank_id_code": "GBDSC",
			"name": "Emelia Jane Brown"
		  },
		  "fx": {
			"contract_reference": "FX123",
			"exchange_rate": "2.00000",
			"original_amount": "200.42",
			"original_currency": "USD"
		  },
		  "payment_purpose": "Paying for goods/services",
		  "payment_type": "Credit",
		  "processing_date": "2017-01-18"
	  }
}`)

var jsonPUT = []byte(`
	{
		"type": "Payment",
		"id": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
		"version": 1,
		"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		  "amount": "1000",
		  "beneficiary_party": {
			"account_name": "W Owens",
			"account_number": "31926819",
			"account_number_code": "BBAN",
			"address": "1 The Beneficiary Localtown SE2",
			"bank_id": "403000",
			"bank_id_code": "GBDSC",
			"name": "Wilfred Jeremiah Owens"
		  },
		  "charges_information": {
			"bearer_code": "SHAR",
			"sender_charges_amount": "5.00",
			 "sender_charges_currency": "GBP",
			"receiver_charges_amount": "1.00",
			"receiver_charges_currency": "USD"
		  },
		  "currency": "GBP",
		  "debtor_party": {
			"account_name": "EJ Brown Black",
			"account_number": "GB29XABC10161234567801",
			"account_number_code": "IBAN",
			"address": "10 Debtor Crescent Sourcetown NE1",
			"bank_id": "203301",
			"bank_id_code": "GBDSC",
			"name": "Emelia Jane Brown"
		  },
		  "fx": {
			"contract_reference": "FX123",
			"exchange_rate": "2.00000",
			"original_amount": "200.42",
			"original_currency": "USD"
		  },
		  "payment_purpose": "Paying for goods/services",
		  "payment_type": "Credit",
		  "processing_date": "2017-01-18"
	  }
}`)

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

// Helper function to create a router during testing
func getRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func TestGetPayments(t *testing.T) {
	r := getRouter()
	r.GET("/v1/payments", getPayments)

	req, _ := http.NewRequest("GET", "/v1/payments", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestAddPayment(t *testing.T) {
	r := getRouter()
	r.POST("/v1/payments", addPayment)

	req, _ := http.NewRequest("POST", "/v1/payments", bytes.NewReader(jsonPOST))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonPOST)))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestUpdatePaymentByID(t *testing.T) {
	r := getRouter()

	r.PUT("/v1/payments", updatePaymentByID)
	req, _ := http.NewRequest("PUT", "/v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", bytes.NewReader(jsonPUT))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonPUT)))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusNotFound
		return statusOK
	})
}

func TestGetPaymentByID(t *testing.T) {
	r := getRouter()
	r.GET("/v1/payments/:id", getPaymentByID)

	req, _ := http.NewRequest("GET", "/v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestDeletePaymentByID(t *testing.T) {
	r := getRouter()
	r.DELETE("/v1/payments/:id", deletePaymentByID)

	req, _ := http.NewRequest("DELETE", "/v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}
