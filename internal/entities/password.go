package entities

import "time"

type UserPass struct {
	id             string    `db:"id"`
	UserId         string    `db:"user_id"`
	HashedPassword string    `db:"hash"`
	CreatedAt      time.Time `db:"created_at"`
}
