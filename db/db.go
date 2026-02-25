package db

import (
	"context"
	"database/sql"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Option func(*bun.DB)

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

	_, err = db.NewCreateTable().
		Model((*model.Visitor)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Unique index on visitor IP for UPSERT
	_, err = db.NewCreateIndex().
		Index("visitors_ip_idx").
		Model((*model.Visitor)(nil)).
		Unique().
		Column("ip").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Migrate image_url -> image_urls for news and exhibits tables.
	// Adds new column if missing, copies data, then drops old column.
	for _, table := range []string{"news", "exhibits"} {
		_, _ = db.ExecContext(ctx,
			"ALTER TABLE "+table+" ADD COLUMN IF NOT EXISTS image_urls text[] DEFAULT '{}'")
		_, _ = db.ExecContext(ctx,
			"UPDATE "+table+" SET image_urls = ARRAY[image_url] WHERE image_url IS NOT NULL AND image_url != '' AND (image_urls IS NULL OR image_urls = '{}')")
		_, _ = db.ExecContext(ctx,
			"ALTER TABLE "+table+" DROP COLUMN IF EXISTS image_url")
	}

	return db, nil
}
