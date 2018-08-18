package controllers

import (
	"encoding/json"
	"invoices/models/request"
	"invoices/services"
	"invoices/util/errors"

	"github.com/gin-gonic/gin"
)

type (
	// InvoiceController represents the controller for operating on the invoices resource
	InvoiceController struct {
		InvoiceService services.InvoiceServiceContract
	}
)

// NewInvoiceController is the constructor of InvoiceController
func NewInvoiceController(ps *services.InvoiceService) *InvoiceController {
	return &InvoiceController{
		InvoiceService: ps,
	}
}

// CreateAction creates a new invoice
func (pc InvoiceController) CreateAction(c *gin.Context) {
	iReq := mrequest.InvoiceCreate{}
	json.NewDecoder(c.Request.Body).Decode(&iReq)

	e := errors.ValidateRequest(&iReq)
	if e != nil {
		c.JSON(e.HttpCode, e)
		return
	}

	iRes, err := pc.InvoiceService.CreateOne(&iReq)

	if err != nil {
		c.JSON(err.HttpCode, err)
		return
	}

	c.JSON(200, iRes)
}

// ListAction list invoices
func (pc InvoiceController) ListAction(c *gin.Context) {
	validSorts := map[string]string{}
	validSorts["InvoiceNo"] = "InvoiceNo"
	validSorts["MovementStartTime"] = "MovementStartTime"
	validSorts["_id"] = "_id"

	validFilters := map[string]string{}
	validFilters["InvoiceNo"] = "InvoiceNo"
	validFilters["CustomerID"] = "CustomerID"
	validFilters["_id"] = "_id"

	qValues := c.Request.URL.Query()
	req := mrequest.NewListRequest(qValues, validSorts, validFilters)

	res, err := pc.InvoiceService.List(req)

	if err != nil {
		c.JSON(err.HttpCode, err)
		return
	}

	c.JSON(200, res)
}
