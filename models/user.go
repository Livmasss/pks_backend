package models

type User struct {
	UserID   string `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
}
