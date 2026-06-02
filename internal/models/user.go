package models

type User struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
