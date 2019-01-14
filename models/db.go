package models

import (
	"fmt"
	"log"
	"github.com/epigos/newsbot/utils"
	"os"
	"reflect"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

var (
	// DS database session
	DS *DataStore
	// Kinds datastore kinds
	Kinds = map[string]string{
		"User":         "Users",
		"Article":      "Articles",
		"Audit":        "AuditRequest",
		"Message":      "Messages",
		"Topic":        "Topics",
		"Subscription": "Subscriptions",
		"SentItem":     "SentItem",
	}
)

// DataStore a struct for database collection
type DataStore struct {
	Client    *datastore.Client
	Context   context.Context
	ProjectID string
	Logger    *utils.Logger
}

// Filter struct for query filters
type Filter struct {
	key   string
	value interface{}
}

// Query daatastore query options
type Query struct {
	Kind    string
	Filters []*Filter
	Limit   int
	Offset  int
	Order   []string
}

// EntitySpec an interface for all entities
type EntitySpec interface {
	Key() *datastore.Key
	SetID(id *datastore.Key)
}

// NewFilter creates a new filter
func NewFilter(key string, value interface{}) *Filter {
	return &Filter{key, value}
}

func (f *Filter) String() string {
	return fmt.Sprintf("%s %s", f.key, f.value)
}

// NewQuery returns new query
func NewQuery(kind string, fs []*Filter, limit, page int, sort ...string) *Query {

	offset := limit * (page - 1)
	if offset < 1 {
		offset = 0
	}
	return &Query{
		Kind:    kind,
		Filters: fs,
		Limit:   limit,
		Offset:  offset,
		Order:   sort,
	}
}

// NewBaseQuery returns new query without limt, sort and skip
func NewBaseQuery(k string, fs []*Filter) *Query {
	return &Query{
		Kind:    k,
		Filters: fs,
	}
}

// Connect to database server
func Connect() {
	// Creates a client.
	projectID := os.Getenv("PROJECT_ID")
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	DS = &DataStore{client, ctx, projectID, utils.NewLogger("models")}

	DS.Logger.Info("Successfully connected to datastore")
}

// NewKey creates a new datastore key
func (d *DataStore) NewKey(kind string) *datastore.Key {
	return datastore.IncompleteKey(kind, nil)
}

// Key builds a new datastore key based on provided id
func (d *DataStore) Key(kind string, id string) *datastore.Key {
	return datastore.NameKey(kind, id, nil)
}

// DecodeKey decode datastore key based on provided id
func (d *DataStore) DecodeKey(id string) *datastore.Key {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		DS.Logger.Error(err)
	}
	return key
}

// GetByKey retrive entity by datastore key
func (d *DataStore) GetByKey(entity EntitySpec) error {
	err := d.Client.Get(d.Context, entity.Key(), entity)
	return err
}

// GetAll retrieves all entities based on given query
func (d *DataStore) GetAll(opts *Query, entities interface{}) ([]*datastore.Key, error) {

	query := datastore.NewQuery(opts.Kind).Offset(opts.Offset)
	if opts.Limit != 0 {
		query = query.Limit(opts.Limit)
	}
	for _, filter := range opts.Filters {
		query = query.Filter(filter.key, filter.value)
	}
	for _, order := range opts.Order {
		query = query.Order(order)
	}

	DS.Logger.Debugf("%+v", query)

	keys, err := d.Client.GetAll(d.Context, query, entities)

	return keys, err
}

// Save saves query
func (d *DataStore) Save(doc EntitySpec) *datastore.Key {

	key := doc.Key()
	val := reflect.ValueOf(doc).Elem()
	now := reflect.ValueOf(time.Now())

	if key.Incomplete() {
		val.FieldByName("Created").Set(now)
	}
	val.FieldByName("Updated").Set(now)

	key, err := d.Client.Put(d.Context, key, doc)
	if err != nil {
		DS.Logger.Panic(err)
	}
	// set id
	doc.SetID(key)

	return key
}

// Delete deletes an entity from its kind
func (d *DataStore) Delete(key *datastore.Key) error {
	return d.Client.Delete(d.Context, key)
}

// Close closes a datastore client
func Close() error {
	return DS.Client.Close()
}
