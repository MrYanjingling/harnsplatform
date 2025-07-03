package auth

import (
	"gorm.io/gorm"
	"harnsplatform/internal/common"
)

type User struct {
	Id     string
	Name   string
	Tenant string
}

func GetCurrentUser(db *gorm.DB) *User {
	ctx := db.Statement.Context
	if user, flag := ctx.Value(common.USER).(*User); flag {
		return user
	}

	return &User{}
}
