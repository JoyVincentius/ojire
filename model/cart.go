package model

import (
	"errors"
	"fmt"
	"ojire/db"
)

type CartItem struct {
	ProductID  int64  `db:"product_id" json:"product_id"`
	Name       string `db:"name" json:"name"`
	Quantity   int    `db:"quantity" json:"quantity"`
	PriceCents int64  `db:"price_cents" json:"price_cents"`
	TotalCents int64  `db:"total_cents" json:"total_cents"`
}

type productInfo struct {
	Name       string `db:"name"`
	PriceCents int64  `db:"price_cents"`
}

// Helper to obtain (or create) a cart ID for a user
func getOrCreateCartID(userID int64) (int64, error) {
	var cartID int64
	err := db.DB.Get(&cartID, "SELECT id FROM carts WHERE user_id = ?", userID)
	if err == nil {
		return cartID, nil
	}
	// not found, create one
	res, err := db.DB.Exec("INSERT INTO carts (user_id) VALUES (?)", userID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func AddItemToCart(userID, productID int64, qty int) error {
	cartID, err := getOrCreateCartID(userID)
	if err != nil {
		return fmt.Errorf("cannot get/create cart: %w", err)
	}

	var stock int
	if err := db.DB.Get(&stock, "SELECT stock FROM products WHERE id = ?", productID); err != nil {
		return fmt.Errorf("stock lookup failed: %w", err)
	}
	if stock < qty {
		return errors.New("insufficient stock")
	}

	var p productInfo
	if err := db.DB.Get(&p, `
            SELECT name, price_cents FROM products WHERE id = ?`, productID); err != nil {
		return fmt.Errorf("product lookup failed: %w", err)
	}

	_, err = db.DB.Exec(`
        INSERT INTO cart_items
            (cart_id, product_id, name, quantity, price_cents, total_cents)
        VALUES
            (?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            quantity = quantity + VALUES(quantity),
            price_cents = VALUES(price_cents),
            total_cents = quantity * price_cents
        `, cartID, productID, p.Name, qty, p.PriceCents, p.PriceCents*int64(qty))

	if err != nil {
		return fmt.Errorf("cannot insert cart item: %w", err)
	}
	return nil
}

// Get cart items with product details
func GetCartItems(userID int64) ([]CartItem, error) {
	var items []CartItem
	query := `
        SELECT ci.product_id, p.name, ci.quantity,
               p.price_cents,
               (ci.quantity * p.price_cents) AS total_cents
        FROM cart_items ci
        JOIN carts c ON ci.cart_id = c.id
        JOIN products p ON ci.product_id = p.id
        WHERE c.user_id = ?
    `
	err := db.DB.Select(&items, query, userID)
	return items, err
}

// Remove specific product from cart
func RemoveItemFromCart(userID, productID int64) error {
	_, err := db.DB.Exec(`
        DELETE ci FROM cart_items ci
        JOIN carts c ON ci.cart_id = c.id
        WHERE c.user_id = ? AND ci.product_id = ?
    `, userID, productID)
	return err
}

// Clear cart completely
func ClearCart(userID int64) error {
	_, err := db.DB.Exec(`
        DELETE ci FROM cart_items ci
        JOIN carts c ON ci.cart_id = c.id
        WHERE c.user_id = ?
    `, userID)
	return err
}
