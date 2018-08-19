package mresponse

import (
	"invoices/models/saft/go_SaftT104"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type Invoice struct {
	ID                                                                  objectid.ObjectID `json:"id,omitempty" bson:"_id"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice `bson:"inline"`
}

type InvoiceCreate struct {
	ID                                                                  string `json:"id,omitempty"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice `bson:"inline"`
}

type InvoiceRead struct {
	ID                                                                  string            `json:"id,omitempty"`
	IDdb                                                                objectid.ObjectID `json:"-" bson:"_id"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice `bson:"inline"`
}

type InvoiceSimple struct {
	ID          string            `json:"id,omitempty"`
	IDdb        objectid.ObjectID `bson:"-" json:"-" bson:"_id"`
	InvoiceNo   string            `json:"InvoiceNo"`
	InvoiceDate time.Time         `json:"InvoiceDate"`
	NetTotal    float64           `json:"NetTotal"`
	GrossTotal  float64           `json:"GrossTotal"`
	TaxPayable  float64           `json:"TaxPayable"`
}

type InvoiceList struct {
	Total   int64           `json:"total"`
	PerPage int64           `json:"per_page"`
	Page    int64           `json:"page"`
	Items   *[]*InvoiceSimple `json:"items"`
}
