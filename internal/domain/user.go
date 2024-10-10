package domain

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUserParams struct {
	Email    string
	Password string // Plain password
}

type CreateHashedUserParams struct {
	Email        string
	PasswordHash string
}

type CreateUser struct {
	Email        string
	PasswordHash string
}

type LoginUserParams struct {
	Email    string
	Password string
}

type UpdateUserParams struct {
	ID       int64
	Email    string
	Password string
}
