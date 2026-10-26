package data

import (
	"database/sql"
	"time"

	"selforder.pramana.net/internal/validator"
)

type Produk struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	Nama      string    `json:"nama"`
	Slug      string    `json:"slug"`
	Deskripsi string    `json:"deskripsi"`
	Harga     uint32    `json:"harga"`
	Stok      uint16    `json:"stok"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

func ValidateProduk(v *validator.Validator, produk *Produk) {
	v.Check(produk.GroupID != 0, "group_id", "must be provided")
	v.Check(produk.GroupID > 0, "group_id", "must be a positive integer")

	v.Check(produk.Nama != "", "nama", "must be provided")
	v.Check(len(produk.Nama) <= 500, "nama", "must not be more than 500 bytes long")

	v.Check(produk.Slug != "", "slug", "must be provided")
	v.Check(len(produk.Slug) <= 500, "slug", "must not be more than 500 bytes long")

	v.Check(produk.Deskripsi != "", "deskripsi", "must be provided")

	v.Check(produk.Harga > 0, "harga", "must be a positive integer")

	v.Check(produk.Stok > 0, "stok", "must be a positive integer")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type ProdukModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m ProdukModel) Insert(produk *Produk) error {
	return nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m ProdukModel) Get(id int64) (*Produk, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m ProdukModel) Update(produk *Produk) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m ProdukModel) Delete(id int64) error {
	return nil
}
