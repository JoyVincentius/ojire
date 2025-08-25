package model

import (
	"database/sql"
	"ojire/db"
)

func BeginTx() (*sql.Tx, error) {
	return db.DB.Begin()
}

// clear cart inside an existing transaction
func ClearCartTx(tx *sql.Tx, userID int64) error {
	_, err := tx.Exec(`
        DELETE ci FROM cart_items ci
        JOIN carts c ON ci.cart_id = c.id
        WHERE c.user_id = ?
    `, userID)
	return err
}
