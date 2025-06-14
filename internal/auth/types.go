package auth

import "gorm.io/gorm"

type User struct {
	Id     string
	Name   string
	Tenant string
}

func GetCurrentUser(db *gorm.DB) *User {
	ctx := db.Statement.Context
	if user, flag := ctx.Value("user").(*User); flag {
		return user
	}

	return &User{}
}
