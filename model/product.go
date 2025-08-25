package model

import "ojire/db"

type Product struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description,omitempty"`
	PriceCents  int64  `db:"price_cents" json:"price_cents"`
	Stock       int    `db:"stock" json:"stock"`
}

func GetAllProducts() ([]Product, error) {
	var products []Product
	err := db.DB.Select(&products, "SELECT id, name, description, price_cents, stock FROM products")
	return products, err
}
