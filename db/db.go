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

func NewDB(ctx context.Context, dsn string) (db *bun.DB, err error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
	))
	db = bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

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
