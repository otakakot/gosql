package model

// アプリケーションで使う構造体
// CQRSにおけるQuery実行時に扱う

type User struct {
	ID   string
	Name string
}
