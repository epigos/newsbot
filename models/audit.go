package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// AuditRequestKind collection name for AuditRequest
const AuditRequestKind = "AuditRequests"

// AuditRequest logs incoming request
type AuditRequest struct {
	ID         string         `json:"id" datastore:"-"`
	UserKey    *datastore.Key `json:"user_id"`
	Path       string         `json:"request_path"`
	Proto      string         `datastore:",noindex" json:"protocol"`
	StatusCode int            `json:"status_code"`
	UserAgent  string         `datastore:",noindex" json:"user_agent"`
	IPAddress  string         `datastore:",noindex" json:"ip_address"`
	Duration   string         `datastore:",noindex" json:"exec_time"`
	Method     string         `json:"method"`
	Query      string         `json:"query" datastore:",noindex"`
	Body       string         `json:"body" datastore:",noindex"`
	Referrer   string         `datastore:",noindex" json:"referrer"`
	IsRobot    bool           `datastore:",noindex" json:"is_robot"`
	Size       int            `datastore:",noindex" json:"size"`
	Created    time.Time      `json:"created"`
	Updated    time.Time      `json:"updated"`
}

// Key method
func (m *AuditRequest) Key() *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if m.ID == "" {
		return DS.NewKey(AuditRequestKind)
	}

	// if Id is already set, we'll just build the Key based
	// on the one provided.
	return DS.Key(AuditRequestKind, m.ID)
}

// SetID set id
func (m *AuditRequest) SetID(key *datastore.Key) {
	m.ID = key.Encode()
}

// Save users
func (m *AuditRequest) Save() {
	DS.Save(m)
}

// GetAllAuditRequest returns audit request
func GetAllAuditRequest() ([]*AuditRequest, error) {
	var results []*AuditRequest

	query := NewQuery(AuditRequestKind, []*Filter{}, 10, 1, "-Created")

	keys, err := DS.GetAll(query, &results)

	for i, key := range keys {
		results[i].SetID(key)
	}

	return results, err
}
