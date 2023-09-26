package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/otakakot/gosql/internal/domain/model"
	"github.com/otakakot/gosql/internal/domain/schema"
)

func main() {
	config, err := pgx.ParseConfig("postgres://postgres:@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	us := &schema.User{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	q1 := db.NewInsert().Model(us).String()

	slog.Info(q1)

	un := &schema.UserName{
		ID:        uuid.NewString(),
		UserID:    us.ID,
		Value:     uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	q2 := db.NewInsert().Model(un).String()

	slog.Info(q2)

	if err := Transact(db, q1, q2); err != nil {
		panic(err)
	}

	user := &model.User{}
	if err := db.NewSelect().
		Model(user).
		ColumnExpr("users.id AS id").
		ColumnExpr("user_names.value AS name").
		Table("users").
		Join("JOIN user_names ON user_names.user_id = users.id").
		Where("users.id = ?", us.ID).
		Scan(context.Background()); err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("%+v", user))
}

func Transact(db *bun.DB, queries ...string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
