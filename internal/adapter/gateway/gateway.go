package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/otakakot/gosql/internal/domain/model"
	"github.com/otakakot/gosql/internal/domain/repository"
	"github.com/otakakot/gosql/internal/domain/schema"
)

var _ repository.Transactor = (*Transactor)(nil)

type Transactor struct {
	DB *bun.DB
}

func (tx *Transactor) Transact(
	ctx context.Context,
	queries ...string,
) error {
	db, err := tx.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, query := range queries {
		if _, err := db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	if err := db.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

var _ repository.UserQuery = (*UserQuery)(nil)

type UserQuery struct {
	DB *bun.DB
}

func (uq *UserQuery) Find(
	ctx context.Context,
	id string,
) (*model.User, error) {
	res := &model.User{}
	if err := uq.DB.NewSelect().
		Model(res).
		ColumnExpr("users.id AS id").
		ColumnExpr("user_names.value AS name").
		Table("users").
		Join("JOIN user_names ON user_names.user_id = users.id").
		Where("users.id = ?", id).
		Scan(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return res, nil
}

var _ repository.UserCommand = (*UserCommand)(nil)

type UserCommand struct {
	DB *bun.DB
}

func (uc *UserCommand) Create(
	ctx context.Context,
	user model.User,
) ([]string, error) {
	now := time.Now()

	q1 := uc.DB.NewInsert().Model(&schema.User{
		ID:        user.ID,
		CreatedAt: now,
		UpdatedAt: now,
		IsDeleted: false,
	}).String()

	q2 := uc.DB.NewInsert().Model(&schema.UserName{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Value:     user.Name,
		CreatedAt: now,
		UpdatedAt: now,
		IsDeleted: false,
	}).String()

	return []string{
		q1,
		q2,
	}, nil
}
