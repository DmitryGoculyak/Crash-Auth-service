package entities

import "time"

type UserEmail struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}
