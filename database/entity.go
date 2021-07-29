// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package database

import (
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

type Entity struct {
	ID        uint
	UUID      pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

// GetID returns entity ID.
func (e *Entity) GetID() uint {
	return e.ID
}

// GetUUID returns entity UUID.
func (e *Entity) GetUUID() string {
	v := e.UUID.Get()
	switch v.(type) {
	case [16]byte:
		vb := v.([16]byte)
		return fmt.Sprintf("%x-%x-%x-%x-%x", vb[0:4], vb[4:6], vb[6:8], vb[8:10], vb[10:16])
	default:
		return ""
	}
}

func (e *Entity) GetCreatedAt() time.Time {
	return e.CreatedAt.Time
}

func (e *Entity) GetUpdatedAt() time.Time {
	return e.UpdatedAt.Time
}

func (e *Entity) GetDeletedAt() time.Time {
	return e.DeletedAt.Time
}
