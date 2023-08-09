package response

import "backend/internal/app/model"

type UserResponse struct {
	User model.User `json:"user"`
}

type LoginResponse struct {
	User      model.User `json:"user"`
	Token     string     `json:"token"`
	ExpiresAt int64      `json:"expiresAt"`
}

type UserInfo struct {
	User model.User `json:"user"`
}
