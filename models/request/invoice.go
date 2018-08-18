package mrequest

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"invoices/models/saft/go_SaftT104"
)

type InvoiceCreate struct {
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice
}

type InvoiceRead struct {
	ID                 objectid.ObjectID `json:"id,omitempty" bson:"_id"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice
}

type InvoiceUpdate struct {
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice
}

type InvoiceDelete struct {
	ID                 objectid.ObjectID `bson:"_id" json:"id,omitempty" valid:"required~Cannot be empty" bson:"_id"`
	go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice
}
