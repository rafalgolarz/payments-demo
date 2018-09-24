package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

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

func getPaymentPOSTPayload() string {
	params := url.Values{}
	params.Add("id", "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43")
	params.Add("type", "Payment")

	return params.Encode()
}

func getPaymentPUTPayload() string {
	params := url.Values{}
	params.Add("version", "1")

	return params.Encode()
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

	postBody := getPaymentPOSTPayload()
	req, _ := http.NewRequest("POST", "/v1/payments", strings.NewReader(postBody))
	req.Header.Add("Content-Type", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestUpdatePaymentByID(t *testing.T) {
	r := getRouter()
	r.PUT("/v1/payments", updatePaymentByID)

	putBody := getPaymentPUTPayload()
	req, _ := http.NewRequest("PUT", "/v1/payments", strings.NewReader(putBody))
	req.Header.Add("Content-Type", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestGetPaymentByID(t *testing.T) {
	r := getRouter()
	r.GET("/v1/payments/:id", getPaymentByID)

	req, _ := http.NewRequest("GET", "/v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestDeletePaymentByID(t *testing.T) {
	r := getRouter()
	r.DELETE("/v1/payments/:id", deletePaymentByID)

	req, _ := http.NewRequest("DELETE", "/v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK
		return statusOK
	})

	TestAddPayment(t)
}
