package models

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash *string
	Provider     string
	CreatedAt    time.Time
}

type RefreshToken struct {
	ID         int64
	UserID     int64
	TokenHash  string
	DeviceInfo *string
	ExpiredAt  time.Time
	Revoked    bool
	CreatedAt  time.Time
}
