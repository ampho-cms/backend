// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReplicaType uint8

const (
	ReplicaTypeRW ReplicaType = iota
	ReplicaTypeRO
)

type Replica struct {
	DSN  string
	Type ReplicaType

	pool *pgxpool.Pool
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (r *Replica) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

// RWReplica makes a new RW replica.
func RWReplica(dsn string) *Replica {
	return &Replica{DSN: dsn, Type: ReplicaTypeRW}
}

// ROReplica makes a new RO replica.
func ROReplica(dsn string) *Replica {
	return &Replica{DSN: dsn, Type: ReplicaTypeRO}
}
