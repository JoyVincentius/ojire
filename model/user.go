package model

import (
	"errors"

	"ojire/db" // ← Replace "your_module_path" with your actual module name from go.mod
)

// User represents the database table schema
type User struct {
	ID       uint64 `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"` // bcrypt hash; not returned in JSON
}

// GetUserByEmail queries the database for a user by email
func GetUserByEmail(email string) (*User, error) {
	// Fail fast if DB isn't initialized
	if db.DB == nil {
		return nil, errors.New("database not initialized")
	}

	const query = `
		SELECT id, email, password
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	var u User
	// sqlx.Get automatically maps result to struct via db tags
	if err := db.DB.Get(&u, query, email); err != nil {
		return nil, err // Include sql.ErrNoRows
	}
	return &u, nil
}
