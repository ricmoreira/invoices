package services

import (
	"context"
	"invoices/models/request"
	"invoices/models/response"
	"invoices/repositories"
	"invoices/util/errors"

	"log"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// InvoiceServiceContract is the abstraction for service layer on roles resource
type InvoiceServiceContract interface {
	CreateOne(*mrequest.InvoiceCreate) (*mresponse.InvoiceCreate, *mresponse.ErrorResponse)
	ReadOne(*mrequest.InvoiceRead) (*mresponse.Invoice, *mresponse.ErrorResponse)
	UpdateOne(*mrequest.InvoiceUpdate) (*mresponse.Invoice, *mresponse.ErrorResponse)
	DeleteOne(*mrequest.InvoiceDelete) (*mresponse.Invoice, *mresponse.ErrorResponse)
	CreateMany(*[]*mrequest.InvoiceCreate) (*[]*mresponse.InvoiceCreate, *mresponse.ErrorResponse)
	List(request *mrequest.ListRequest) (*mresponse.InvoiceList, *mresponse.ErrorResponse)
}

// InvoiceService is the layer between http client and repository for Invoice resource
type InvoiceService struct {
	InvoiceRepository *repositories.InvoiceRepository
}

// NewInvoiceService is the constructor of InvoiceService
func NewInvoiceService(pr *repositories.InvoiceRepository) *InvoiceService {
	return &InvoiceService{
		InvoiceRepository: pr,
	}
}

// CreateOne saves provided model instance to database
func (this *InvoiceService) CreateOne(request *mrequest.InvoiceCreate) (*mresponse.InvoiceCreate, *mresponse.ErrorResponse) {

	// validate request
	e := errors.ValidateRequest(request)
	if e != nil {
		return nil, e
	}

	res, err := this.InvoiceRepository.CreateOne(request)

	if err != nil {
		errR := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
		return nil, errR
	}

	id := res.InsertedID.(objectid.ObjectID)
	p := mresponse.InvoiceCreate{
		ID: id.Hex(),
	}

	return &p, nil
}

// TODO: implement
func (this *InvoiceService) ReadOne(p *mrequest.InvoiceRead) (*mresponse.Invoice, *mresponse.ErrorResponse) {
	return nil, nil
}

// TODO: implement
func (this *InvoiceService) UpdateOne(p *mrequest.InvoiceUpdate) (*mresponse.Invoice, *mresponse.ErrorResponse) {
	return nil, nil
}

// TODO: implement
func (this *InvoiceService) DeleteOne(p *mrequest.InvoiceDelete) (*mresponse.Invoice, *mresponse.ErrorResponse) {
	return nil, nil
}

// CreateMany saves many Invoices in one bulk operation
func (this *InvoiceService) CreateMany(request *[]*mrequest.InvoiceCreate) (*[]*mresponse.InvoiceCreate, *mresponse.ErrorResponse) {

	res, err := this.InvoiceRepository.InsertMany(request)

	if err != nil {
		mngBulkError := err.(mongo.BulkWriteError)
		writeErrors := mngBulkError.WriteErrors
		for _, err := range writeErrors {
			log.Println(err)
		}
	}

	result := make([]*mresponse.InvoiceCreate, len(res.InsertedIDs))
	for i, insertedID := range res.InsertedIDs {
		id := insertedID.(objectid.ObjectID)
		result[i] = &mresponse.InvoiceCreate{
			ID: id.Hex(),
		}
	}

	return &result, nil
}

// List returns a list of Invoices with pagination and filtering options
func (this *InvoiceService) List(request *mrequest.ListRequest) (*mresponse.InvoiceList, *mresponse.ErrorResponse) {

	total, perPage, page, cursor, err := this.InvoiceRepository.List(request)

	if err != nil {
		e := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
		return nil, e
	}

	docs := []*mresponse.InvoiceRead{}

	for cursor.Next(context.Background()) {
		doc := mresponse.InvoiceRead{}
		if err := cursor.Decode(&doc); err != nil {
			errR := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
			return nil, errR
		}

		doc.ID = doc.IDdb.Hex()

		docs = append(docs, &doc)
	}

	resp := mresponse.InvoiceList{
		Total:   total,
		PerPage: perPage,
		Page:    page,
		Items:   &docs,
	}
	return &resp, nil
}
