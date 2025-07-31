package core

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}

func (u *User) CheckPassword(password string) bool {
	now := time.Now()
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	fmt.Println("Password check took", time.Since(now))

	return err == nil
}

func NewUser(username, password string, bcryptCost int) (*User, error) {
	now := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	fmt.Println("Password hashing took", time.Since(now))

	if err != nil {
		return nil, err
	}

	return &User{
		Username: username,
		Password: string(hash),
	}, nil
}
