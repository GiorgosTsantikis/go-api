// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package database

import (
	"context"
)

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, name, email, "emailVerified", image, "createdAt", "updatedAt", role, banned, "banReason", "banExpires", "phoneNumber", "phoneNumberVerified" FROM "user"
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.EmailVerified,
			&i.Image,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Role,
			&i.Banned,
			&i.BanReason,
			&i.BanExpires,
			&i.PhoneNumber,
			&i.PhoneNumberVerified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, "emailVerified", image, "createdAt", "updatedAt", role, banned, "banReason", "banExpires", "phoneNumber", "phoneNumberVerified" FROM "user" WHERE email=$1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.EmailVerified,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
		&i.Banned,
		&i.BanReason,
		&i.BanExpires,
		&i.PhoneNumber,
		&i.PhoneNumberVerified,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, "emailVerified", image, "createdAt", "updatedAt", role, banned, "banReason", "banExpires", "phoneNumber", "phoneNumberVerified" FROM "user" WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.EmailVerified,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
		&i.Banned,
		&i.BanReason,
		&i.BanExpires,
		&i.PhoneNumber,
		&i.PhoneNumberVerified,
	)
	return i, err
}

const getUserBySession = `-- name: GetUserBySession :one
SELECT u.id, u.name, u.email, u."emailVerified", u.image, u."createdAt", u."updatedAt", u.role, u.banned, u."banReason", u."banExpires", u."phoneNumber", u."phoneNumberVerified" FROM "user" u JOIN "session" s ON
    u.id = s."userId" WHERE s.token = $1
`

func (q *Queries) GetUserBySession(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserBySession, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.EmailVerified,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
		&i.Banned,
		&i.BanReason,
		&i.BanExpires,
		&i.PhoneNumber,
		&i.PhoneNumberVerified,
	)
	return i, err
}
