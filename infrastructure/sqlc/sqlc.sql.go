// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: sqlc.sql

package sqlc

import (
	"context"
	"database/sql"
)

const existsUser = `-- name: ExistsUser :one
SELECT
  EXISTS(
    SELECT
      1
    FROM
      user
    WHERE
      email = ?
  )
`

func (q *Queries) ExistsUser(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRowContext(ctx, existsUser, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getLastInsertID = `-- name: GetLastInsertID :one
SELECT
  LAST_INSERT_ID()
`

func (q *Queries) GetLastInsertID(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLastInsertID)
	var last_insert_id int64
	err := row.Scan(&last_insert_id)
	return last_insert_id, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
  user_id,
  user_name,
  password
FROM
  user
WHERE
  email = ?
`

type GetUserByEmailRow struct {
	UserID   uint32
	UserName string
	Password string
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(&i.UserID, &i.UserName, &i.Password)
	return i, err
}

const insertUser = `-- name: InsertUser :execresult
INSERT INTO
  user (user_name, email, password)
VALUES
  (?, ?, ?)
`

type InsertUserParams struct {
	UserName string
	Email    string
	Password string
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertUser, arg.UserName, arg.Email, arg.Password)
}
