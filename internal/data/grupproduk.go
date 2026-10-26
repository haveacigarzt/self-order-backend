package data

import (
	"database/sql"
	"time"

	"selforder.pramana.net/internal/validator"
)

type GrupProduk struct {
	ID        int64     `json:"id"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

func ValidateGrupProduk(v *validator.Validator, grupProduk *GrupProduk) {
	v.Check(grupProduk.Nama != "", "nama", "must be provided")
	v.Check(len(grupProduk.Nama) <= 500, "nama", "must not be more than 500 bytes long")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type GrupProdukModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m GrupProdukModel) Insert(grupProduk *GrupProduk) error {
	query := `
	INSERT INTO grup_produk (nama)
	VALUES ($1)
	RETURNING id, created_at, version`

	args := []any{grupProduk.Nama}

	return m.DB.QueryRow(query, args...).Scan(&grupProduk.ID, &grupProduk.CreatedAt, &grupProduk.Version)
}
