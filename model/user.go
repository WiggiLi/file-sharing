package model

import (
	"time"
	"context"
	"github.com/google/uuid"
)

// User is a JSON user
type User struct {
	ID        uuid.UUID  `json:"id" pg:"id,notnull,pk"`
	Name 	  string     `json:"name" validate:"required" pg:"name,notnull"`
	Email     string 	 `json:"email" pg:"email,notnull"`
	Password  string 	 `json:"password" pg:"password,notnull"`
	Token     string 	 `json:"token";sql:"-" pg:"token"`
	CreatedAt time.Time  `json:"created_at" pg:"created_at,notnull"`
}

type FilesOfUsers struct {
	User_ID   uuid.UUID  `pg:"id_user,type:uuid,notnull"`
	File_ID	  uuid.UUID  `pg:"id_file,type:uuid,notnull"`
}


// UserService is a service for users
//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	GetFileNamesByUserID(context.Context, uuid.UUID) (*[]File, error)
	CreateUser(context.Context, *User) (*User, error)
	LoginUser(context.Context, *User) (*User, error)
	//DeleteUser(context.Context, uuid.UUID) error
}

// UserRepo is a store for users
//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepoPg interface {
	GetFileNamesByUserID(context.Context, uuid.UUID) (*[]File, error)
	CreateUser(context.Context, *User) (*User, error)
	LoginUser(context.Context, *User) (*User, error)
	//DeleteUser(context.Context, uuid.UUID) error
}

// UserRepoRedis is a store for users
type UserRepoRedis interface {
	GetFileNamesByUserID(context.Context, string) (*[]File, error)
	CreateFileNamesByUserID(context.Context, string, int, *[]File) error
	DeleteFileNamesByUserID(context.Context, string) error
}


