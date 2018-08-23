package repositories

import (
	"context"
	"invoices/models/request"
	"invoices/models/response"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
)

// InvoiceRepository performs CRUD operations on users resource
type InvoiceRepository struct {
	invoices MongoCollection
}

// NewInvoiceRepository is the constructor for InvoiceRepository
func NewInvoiceRepository(db *DBCollections) *InvoiceRepository {
	return &InvoiceRepository{invoices: db.Invoice}
}

// CreateOne saves provided model instance to database
func (this *InvoiceRepository) CreateOne(request *mrequest.InvoiceCreate) (*mongo.InsertOneResult, error) {

	return this.invoices.InsertOne(context.Background(), request)
}

// ReadOne returns a invoice based on InvoiceCode sent in request
// TODO: implement better query based on full request and not only the ProducCode
func (this *InvoiceRepository) ReadOne(p *mrequest.InvoiceRead) (*mresponse.Invoice, error) {
	result := this.invoices.FindOne(
		context.Background(),
		bson.NewDocument(bson.EC.String("InvoiceNo", string(p.InvoiceNo))),
	)

	res := mresponse.Invoice{}
	err := result.Decode(p)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// TODO: implement
func (this *InvoiceRepository) UpdateOne(p *mrequest.InvoiceUpdate) (*mresponse.Invoice, error) {
	return nil, nil
}

// TODO: implement
func (this *InvoiceRepository) DeleteOne(p *mrequest.InvoiceDelete) (*mresponse.Invoice, error) {
	return nil, nil
}

func (this *InvoiceRepository) InsertMany(request *[]*mrequest.InvoiceCreate) (*mongo.InsertManyResult, error) {
	// transform to []interface{} (https://golang.org/doc/faq#convert_slice_of_interface)
	s := make([]interface{}, len(*request))
	for i, v := range *request {
		s[i] = v
	}

	// { ordered: false } ordered is false in order to don't stop execution because an error ocurred on one of the inserts
	opt := insertopt.Ordered(false)
	return this.invoices.InsertMany(context.Background(), s, opt)
}

func (this *InvoiceRepository) List(req *mrequest.ListRequest) (int64, int64, int64, mongo.Cursor, error) {

	args := []*bson.Element{}

	for key, value := range req.Filters {
		if key != "_id" { // filter by text fields
			pattern := value.(string)
			elem := bson.EC.Regex(key, pattern, "i")
			args = append(args, elem)
		} else { // filter by _id
			elem := bson.EC.String(key, value.(string))
			args = append(args, elem)
		}
	}

	total, e := this.invoices.Count(
		context.Background(),
		bson.NewDocument(args...),
	)

	sorting := map[string]int{}
	var sortingValue int
	if req.Order == "reverse" {
		sortingValue = -1
	} else {
		sortingValue = 1
	}
	sorting[req.Sort] = sortingValue

	perPage := int64(req.PerPage)
	page := int64(req.Page)
	cursor, e := this.invoices.Find(
		context.Background(),
		bson.NewDocument(args...),
		findopt.Sort(sorting),
		findopt.Skip(int64(req.PerPage*(req.Page-1))),
		findopt.Limit(perPage),
	)

	return total, perPage, page, cursor, e
}
