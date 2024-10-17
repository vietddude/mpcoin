package domain

import "github.com/google/uuid"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignupResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type LoginResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
