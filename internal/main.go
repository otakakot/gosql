package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/otakakot/gosql/internal/adapter/gateway"
	"github.com/otakakot/gosql/internal/domain/model"
)

func main() {
	config, err := pgx.ParseConfig("postgres://postgres:@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	ctx := context.Background()

	tx := gateway.Transactor{DB: db}

	uc := gateway.UserCommand{DB: db}

	uq := gateway.UserQuery{DB: db}

	mdl := model.User{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}

	qs, err := uc.Create(ctx, mdl)
	if err != nil {
		panic(err)
	}

	if err := tx.Transact(ctx, qs...); err != nil {
		panic(err)
	}

	got, err := uq.Find(ctx, mdl.ID)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("%+v", got))
}
