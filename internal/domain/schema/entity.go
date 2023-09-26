package schema

// データベースで扱う構造体
// CQRSにおけるCommand実行時に扱う

import "time"

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}

type UserName struct {
	ID        string
	UserID    string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}
