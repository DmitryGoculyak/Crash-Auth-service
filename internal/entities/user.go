package entities

import "time"

type User struct {
	Id        string    `db:"id"`
	FullName  string    `db:"full_name"`
	CreatedAt time.Time `db:"created_at"`
}
