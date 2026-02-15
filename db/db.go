package db

import (
	"context"
	"database/sql"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
)

type Option func(*bun.DB)

// WithOtel добавляет OpenTelemetry хук для трейсинга SQL-запросов.
func WithOtel() Option {
	return func(db *bun.DB) {
		db.AddQueryHook(bunotel.NewQueryHook(
			bunotel.WithFormattedQueries(true),
		))
	}
}

// WithDebug добавляет debug-хук для логирования SQL-запросов в stdout.
func WithDebug(verbose bool) Option {
	return func(db *bun.DB) {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(verbose),
		))
	}
}

func NewDB(ctx context.Context, dsn string, opts ...Option) (db *bun.DB, err error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
	))
	db = bun.NewDB(sqldb, pgdialect.New())

	for _, opt := range opts {
		opt(db)
	}

	_, err = db.NewCreateTable().
		Model((*model.News)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.NewCreateIndex().
		Index("news_id_idx").
		Model((*model.News)(nil)).
		Unique().
		Column("id").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.NewCreateTable().
		Model((*model.Exhibition)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	_, err = db.NewCreateIndex().
		Index("exhibition_id_idx").
		Model((*model.Exhibition)(nil)).
		Unique().
		Column("id").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.NewCreateTable().
		Model((*model.Exhibit)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	_, err = db.NewCreateIndex().
		Index("exhibit_id_idx").
		Model((*model.Exhibit)(nil)).
		Unique().
		Column("id").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
