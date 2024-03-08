package dto

import "github.com/nseptio/food-recipes/internal/entity"

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required,max=100"` // This field can be either username or email
	Password   string `json:"password" validate:"required,max=100"`
}

type UserRegisterRequest struct {
	Username  string `json:"username" validate:"required,max=50"`
	Email     string `json:"email" validate:"required,max=100"`
	Password  string `json:"password" validate:"required,max=100"`
	FirstName string `json:"firstName" validate:"required,max=50"`
	LastName  string `json:"lastName" validate:"required,max=50"`
}

type UserResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type AuthResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func UserToResponse(user entity.User) *UserResponse {
	return &UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
