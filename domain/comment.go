package domain

import "time"

type Comment struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	Author    User
}
