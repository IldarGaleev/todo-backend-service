// Package faketempdb implements fake temp storage provider
package faketempdb

import (
	"context"
	"log/slog"
)

type FakeTempDB struct {
	log         *slog.Logger
	database    map[uint64]interface{}
	lastTokenID uint64
}

func New(logger *slog.Logger) *FakeTempDB {
	return &FakeTempDB{
		log:         logger,
		database:    make(map[uint64]interface{}, 10),
		lastTokenID: 0,
	}
}

func (d *FakeTempDB) IsJWTRevoked(ctx context.Context, id uint64) bool {
	_, ok := d.database[id]
	return ok
}

func (d *FakeTempDB) RevokeJWT(ctx context.Context, id uint64) {
	var a interface{}
	d.database[id] = a
}

func (d *FakeTempDB) CreateNewID(ctx context.Context) uint64 {
	d.lastTokenID++
	return d.lastTokenID
}
