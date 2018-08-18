package mresponse

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"invoices/models/saft/go_SaftT104"
)

type Invoice struct {
	ID                 objectid.ObjectID `json:"id,omitempty" bson:"_id"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice
}

type InvoiceCreate struct {
	ID string `json:"id,omitempty"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice
}

type InvoiceRead struct {
	ID                 string            `json:"id,omitempty"`
	IDdb               objectid.ObjectID `json:"-" bson:"_id"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice
}

type InvoiceList struct {
	Total   int64           `json:"total"`
	PerPage int64           `json:"per_page"`
	Page    int64           `json:"page"`
	Items   *[]*InvoiceRead `json:"items"`
}
