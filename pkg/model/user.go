package model

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("User not found")
)
type LoginUserQuery struct {
	Username   string
	Password   string
	User       *User
	IpAddress  string
}
type GetUserByLoginQuery struct {
	LoginOrEmail string
	Result       *User
}
type User struct {
	Id            int64
	Email         string
	Name          string
	Login         string
	Password      string
	Salt          string
	Rands         string

}


// ------------------------
// DTO & Projections

type SignedInUser struct {
	UserId         int64
	OrgId          int64
	OrgName        string
	Login          string
	Name           string
	Email          string
	ApiKeyId       int64
	OrgCount       int
	IsMerkutAdmin bool
	IsAnonymous    bool
	LastSeenAt     time.Time
}


func (u *SignedInUser) NameOrFallback() string {
	if u.Name != "" {
		return u.Name
	} else if u.Login != "" {
		return u.Login
	} else {
		return u.Email
	}
}

type UpdateUserLastSeenAtCommand struct {
	UserId int64
}


type UserProfileDTO struct {
	Id             int64  `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Login          string `json:"login"`
	Theme          string `json:"theme"`
	OrgId          int64  `json:"orgId"`
}

type UserSearchHitDTO struct {
	Id            int64     `json:"id"`
	Name          string    `json:"name"`
	Login         string    `json:"login"`
	Email         string    `json:"email"`
	AvatarUrl     string    `json:"avatarUrl"`
	IsAdmin       bool      `json:"isAdmin"`
	LastSeenAt    time.Time `json:"lastSeenAt"`
	LastSeenAtAge string    `json:"lastSeenAtAge"`
}

type UserIdDTO struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}