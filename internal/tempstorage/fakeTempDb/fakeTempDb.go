package faketempdb

import (
	"context"
	"log/slog"
)

type FakeTempDb struct {
	log         *slog.Logger
	database    map[uint64]interface{}
	lastTokenId uint64
}

func New(logger *slog.Logger) *FakeTempDb {
	return &FakeTempDb{
		log:         logger,
		database:    make(map[uint64]interface{}, 10),
		lastTokenId: 0,
	}
}

func (d *FakeTempDb) IsJWTRevoked(ctx context.Context, id uint64) bool {
	_, ok := d.database[id]
	return ok
}

func (d *FakeTempDb) RevokeJWT(ctx context.Context, id uint64) {
	var a interface{}
	d.database[id] = a
}

func (d *FakeTempDb) CreateNewId(ctx context.Context) uint64 {
	d.lastTokenId++
	return d.lastTokenId
}
