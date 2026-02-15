package main

import (
	"context"
	"os"

	"github.com/WhiCu/school-museum/db"
	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/internal/config"
	"github.com/google/uuid"
)

// import "github.com/WhiCu/school-museum/cmd"

func main() {
	// cmd.Execute()
	os.Setenv("PATH_CONFIG", "./config/config.kdl")
	defer os.Unsetenv("PATH_CONFIG")

	cfg := config.MustLoad[config.Config](config.KDL)

	db, err := db.NewDB(context.Background(), cfg.Storage.DSN())
	if err != nil {
		panic(err)
	}

	news := &model.News{
		Title:   []byte("Hello, World!"),
		Content: []byte("This is a news item."),
	}
	err = db.NewInsert().Model(news).Returning("id").Scan(context.Background(), &news.ID)
	if err != nil {
		panic(err)
	}
	println(news.ID.String())

	id, err := uuid.Parse("2b7130c9-77ce-49e7-a0c7-b2d551c809a9")
	if err != nil {
		panic(err)
	}

	err = db.NewSelect().Model((*model.News)(nil)).Where("id = ?", id).Scan(context.Background(), news)
	if err != nil {
		panic(err)
	}
	println(string(news.Title))

}
