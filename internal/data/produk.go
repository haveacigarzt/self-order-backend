package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"selforder.pramana.net/internal/validator"
)

type Produk struct {
	ID        int64     `json:"id"`
	GrupID    int64     `json:"grup_id"`
	Nama      string    `json:"nama"`
	Slug      string    `json:"slug"`
	Deskripsi string    `json:"deskripsi"`
	Harga     uint32    `json:"harga"`
	Stok      uint16    `json:"stok"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
	GrupNama  string    `json:"grup_nama,omitempty"`
}

func ValidateProduk(v *validator.Validator, produk *Produk) {
	v.Check(produk.GrupID != 0, "grup_id", "must be provided")
	v.Check(produk.GrupID > 0, "grup_id", "must be a positive integer")

	v.Check(produk.Nama != "", "nama", "must be provided")
	v.Check(utf8.RuneCountInString(produk.Nama) <= 100, "nama", "must not be more than 100 characters long")

	v.Check(produk.Slug != "", "slug", "must be provided")
	v.Check(utf8.RuneCountInString(produk.Slug) <= 120, "slug", "must not be more than 120 characters long")

	v.Check(produk.Deskripsi != "", "deskripsi", "must be provided")
	v.Check(utf8.RuneCountInString(produk.Deskripsi) <= 255, "deskripsi", "must not be more than 255 characters long")

	v.Check(produk.Harga > 0, "harga", "must be a positive integer")
	v.Check(produk.Harga <= 2000000, "harga", "must not be more than 2000000")

	v.Check(produk.Stok > 0, "stok", "must be a positive integer")
	v.Check(produk.Stok <= 30000, "stok", "must not be more than 30000")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type ProdukModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m ProdukModel) Insert(produk *Produk) error {
	query := `
	INSERT INTO produk (nama, deskripsi, grup_id, slug, stok, harga)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, version`

	args := []any{produk.Nama, produk.Deskripsi, produk.GrupID, produk.Slug, produk.Stok, produk.Harga}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&produk.ID, &produk.CreatedAt, &produk.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m ProdukModel) Get(id int64) (*Produk, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT p.id, p.created_at, p.nama, p.deskripsi, p.grup_id, gp.nama AS grup_nama, p.slug, p.stok, p.harga, p.version
	FROM produk p JOIN grup_produk gp ON p.grup_id = gp.id
	WHERE p.id = $1`

	var produk Produk

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&produk.ID,
		&produk.CreatedAt,
		&produk.Nama,
		&produk.Deskripsi,
		&produk.GrupID,
		&produk.GrupNama,
		&produk.Slug,
		&produk.Stok,
		&produk.Harga,
		&produk.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &produk, nil
}

func (m ProdukModel) Update(produk *Produk) error {
	query := `
	UPDATE produk
	SET nama = $1, deskripsi = $2, grup_id = $3, slug = $4, stok = $5, harga = $6, version = version + 1, updated_at = now()
	WHERE id = $7 AND version = $8
	RETURNING version`

	args := []any{
		produk.Nama,
		produk.Deskripsi,
		produk.GrupID,
		produk.Slug,
		produk.Stok,
		produk.Harga,
		produk.ID,
		produk.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&produk.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

// Add a placeholder method for deleting a specific record from the movies table.
func (m ProdukModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM produk
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

func (m ProdukModel) GetAll(nama string, filters Filters) ([]*Produk, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), p.id, p.created_at, p.nama, p.deskripsi, p.grup_id, gp.nama AS grup_nama, p.slug, p.stok, p.harga, p.version
		FROM produk p JOIN grup_produk gp ON p.grup_id = gp.id
		WHERE (to_tsvector('simple', p.nama) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{nama, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	produkList := []*Produk{}

	for rows.Next() {
		var produk Produk

		err := rows.Scan(
			&totalRecords,
			&produk.ID,
			&produk.CreatedAt,
			&produk.Nama,
			&produk.Deskripsi,
			&produk.GrupID,
			&produk.GrupNama,
			&produk.Slug,
			&produk.Stok,
			&produk.Harga,
			&produk.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		produkList = append(produkList, &produk)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return produkList, metadata, nil
}
