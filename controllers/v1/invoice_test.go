package controllers

import (
	"invoices/models/request"
	"invoices/models/response"
	"invoices/util/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// stub InvoiceService behaviour
type MockInvoiceService struct{}

// mocked behaviour for CreateOne
func (ps *MockInvoiceService) CreateOne(pReq *mrequest.InvoiceCreate) (*mresponse.InvoiceCreate, *mresponse.ErrorResponse) {
	// validate request
	err := errors.ValidateRequest(pReq)
	if err != nil {
		return nil, err
	}

	pRes := mresponse.InvoiceCreate{}
	pRes.ID = "some-unique-id"

	return &pRes, nil
}

// mocked behaviour for ReadOne
func (ps *MockInvoiceService) ReadOne(p *mrequest.InvoiceRead) (*mresponse.Invoice, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

// mocked behaviour for UpdateOne
func (ps *MockInvoiceService) UpdateOne(p *mrequest.InvoiceUpdate) (*mresponse.Invoice, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

// mocked behaviour for DeleteOne
func (ps *MockInvoiceService) DeleteOne(p *mrequest.InvoiceDelete) (*mresponse.Invoice, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

func (ps *MockInvoiceService) CreateMany(*[]*mrequest.InvoiceCreate) (*[]*mresponse.InvoiceCreate, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

func (ps *MockInvoiceService) List(*mrequest.ListRequest) (*mresponse.InvoiceList, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}
func TestCreateInvoiceAction(t *testing.T) {

	// Mock the server

	// Switch to test mode in order to don't get such noisy output
	gin.SetMode(gin.TestMode)

	pps := &MockInvoiceService{}

	pc := InvoiceController{
		InvoiceService: pps,
	}

	r := gin.Default()

	r.POST("/api/v1/invoice", pc.CreateAction)

	// TEST SUCCESS



	// Mock a request
	body := mrequest.InvoiceCreate{}
	body.InvoiceNo = "FS 1/1861"
	body.Period = 2
	body.CustomerID = "113725"
	body.InvoiceType = "FS"

	jsonValue, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/invoice", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder in order to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Do asssertions
	if w.Code != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(w.Body)
		bodyString := string(bodyBytes)

		t.Fatalf("Expected to get status %d but instead got %d\nResponse body:\n%s", http.StatusOK, w.Code, bodyString)
	}
}
