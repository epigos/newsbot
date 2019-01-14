package web

import (
	"context"

	"cloud.google.com/go/datastore"
	"golang.org/x/crypto/acme/autocert"
)

const (
	certsKind = "Certs"
)

type (
	letsEncryptCert struct {
		Data []byte `datastore:"data,noindex"`
	}
	// DatastoreCertCache datastore cache
	DatastoreCertCache struct {
		client *datastore.Client
	}
)

// NewDatastoreCertCache returns new datastore cache
func NewDatastoreCertCache(client *datastore.Client) *DatastoreCertCache {
	return &DatastoreCertCache{
		client: client,
	}
}

// Get get cache data
func (d *DatastoreCertCache) Get(ctx context.Context, key string) ([]byte, error) {
	var cert letsEncryptCert
	k := datastore.NameKey(certsKind, key, nil)
	if err := d.client.Get(ctx, k, &cert); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}
	return cert.Data, nil
}

// Put stores the data in the cache under the specified key.
// Underlying implementations may use any data storage format,
// as long as the reverse operation, Get, results in the original data.
func (d *DatastoreCertCache) Put(ctx context.Context, key string, data []byte) error {
	k := datastore.NameKey(certsKind, key, nil)
	cert := letsEncryptCert{
		Data: data,
	}
	if _, err := d.client.Put(ctx, k, &cert); err != nil {
		return err
	}
	return nil
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (d *DatastoreCertCache) Delete(ctx context.Context, key string) error {
	k := datastore.NameKey(certsKind, key, nil)
	return d.client.Delete(ctx, k)
}
