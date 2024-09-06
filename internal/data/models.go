package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Models struct contain the other models our application needs.
// For example: Users UserModel
type Models struct {
}

// NewModels returns a new Models struct.
// For example: Users: UserModel{DB:db}
func NewModels(db *pgxpool.Pool) Models {
	return Models{}
}
