// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        int64
	GoogleID  string
	Email     string
	CreatedAt pgtype.Timestamptz
}
