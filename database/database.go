// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Database represents a database replicas pool.
type Database struct {
	rwReplicas []*Replica
	roReplicas []*Replica
	rwIndex    int
	roIndex    int
}

// GetReplica returns a next replica of type t from a pool.
func (d *Database) GetReplica(t ReplicaType) *Replica {
	var (
		list  []*Replica
		index *int
	)

	if len(d.roReplicas) == 0 {
		list = d.rwReplicas
		index = &d.rwIndex
	} else if len(d.rwReplicas) == 0 {
		list = d.roReplicas
		index = &d.roIndex
	} else {
		switch t {
		case ReplicaTypeRW:
			list = d.rwReplicas
			index = &d.rwIndex
		default:
			list = d.roReplicas
			index = &d.roIndex
		}
	}

	*index++
	if *index == len(list) {
		*index = 0
	}

	return list[*index]
}

// Exec executes a non-SELECT query using an RW replica.
func (d *Database) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.GetReplica(ReplicaTypeRW).pool.Exec(ctx, sql, args...)
}

// Query executes a SELECT query using an RO replica.
func (d *Database) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return d.GetReplica(ReplicaTypeRO).pool.Query(ctx, sql, args...)
}

func (d *Database) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return d.GetReplica(ReplicaTypeRO).pool.QueryRow(ctx, sql, args...)
}

func (d *Database) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return d.GetReplica(ReplicaTypeRO).pool.QueryFunc(ctx, sql, args, scans, f)
}

func (d *Database) SendBatch(ctx context.Context, t ReplicaType, b *pgx.Batch) pgx.BatchResults {
	return d.GetReplica(t).pool.SendBatch(ctx, b)
}

func (d *Database) Begin(ctx context.Context, t ReplicaType) (pgx.Tx, error) {
	return d.GetReplica(t).pool.Begin(ctx)
}

func (d *Database) BeginTx(ctx context.Context, t ReplicaType, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return d.GetReplica(t).pool.BeginTx(ctx, txOptions)
}

func (d *Database) BeginFunc(ctx context.Context, t ReplicaType, f func(pgx.Tx) error) error {
	return d.GetReplica(t).pool.BeginFunc(ctx, f)
}

func (d *Database) BeginTxFunc(ctx context.Context, t ReplicaType, txOptions pgx.TxOptions, f func(pgx.Tx) error) error {
	return d.GetReplica(t).pool.BeginTxFunc(ctx, txOptions, f)
}

func (d *Database) Migrate(source migration.Source, direction migration.Direction, max int) (int, error) {
	driver, err := postgres.New(d.GetReplica(ReplicaTypeRW).DSN)
	if err != nil {
		return 0, err
	}

	return migration.Migrate(driver, source, direction, max)
}

// SelectOne performs a SELECT query and scans a single row into a struct or map.
func (d *Database) SelectOne(ctx context.Context, dest interface{}, sql string, args ...interface{}) error {
	rows, err := d.Query(ctx, sql, args...)
	if err != nil {
		return err
	}

	if err = pgxscan.ScanOne(dest, rows); err != nil {
		return err
	}

	return nil
}

// SelectAll performs a SELECT query and scans all rows into a slice of structs or maps.
func (d *Database) SelectAll(ctx context.Context, dest interface{}, sql string, args ...interface{}) error {
	rows, err := d.Query(ctx, sql, args...)
	if err != nil {
		return err
	}

	if err = pgxscan.ScanAll(dest, rows); err != nil {
		return err
	}

	return nil
}

// New creates a new database replicas pool.
func New(ctx context.Context, replica ...*Replica) (*Database, error) {
	if len(replica) == 0 {
		return nil, errors.New("failed to initialize a database with zero replicas")
	}

	db := Database{}

	for _, r := range replica {
		pool, err := pgxpool.Connect(ctx, r.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize connection pool: %v", err)
		}

		r.pool = pool
		switch r.Type {
		case ReplicaTypeRW:
			db.rwReplicas = append(db.rwReplicas, r)
		case ReplicaTypeRO:
			db.roReplicas = append(db.roReplicas, r)
		}
	}

	return &db, nil
}

func NewFromDSN(ctx context.Context, dsn ...string) (*Database, error) {
	var replicas []*Replica

	for _, d := range dsn {
		rt := ReplicaTypeRW
		if strings.HasPrefix(d, "ro:") {
			rt = ReplicaTypeRO
			d = strings.Replace(d, "ro:", "", 1)
		}

		replicas = append(replicas, &Replica{Type: rt, DSN: d})
	}

	return New(ctx, replicas...)
}
