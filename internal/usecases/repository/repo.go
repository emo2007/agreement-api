package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

type DBTX interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

var TxCtxKey = "dbtx"

func (r *Repository) Transaction(ctx context.Context, fn func(context.Context) error) (err error) {
	var tx *sql.Tx = new(sql.Tx)

	hasExternalTx := hasExternalTransaction(ctx)

	defer func() {
		if hasExternalTx {
			if err != nil {
				err = fmt.Errorf("error perform operation. %w", err)
				return
			}

			return
		}

		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(fmt.Errorf("error rollback transaction. %w", rbErr), err)
				return
			}

			err = fmt.Errorf("error execute transactional operation. %w", err)

			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(fmt.Errorf("error rollback transaction. %w", rbErr), commitErr, err)

				return
			}

			err = fmt.Errorf("error commit transaction. %w", err)
		}
	}()

	if !hasExternalTx {
		tx, err = r.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
		})
		if err != nil {
			return fmt.Errorf("error begin transaction. %w", err)
		}

		ctx = context.WithValue(ctx, TxCtxKey, tx)
	}

	return fn(ctx)
}

func (s *Repository) Conn(ctx context.Context) DBTX {
	if tx, ok := ctx.Value(TxCtxKey).(*sql.Tx); ok {
		return tx
	}

	return s.db
}

func hasExternalTransaction(ctx context.Context) bool {
	if _, ok := ctx.Value(TxCtxKey).(*sql.Tx); ok {
		return true
	}

	return false
}
