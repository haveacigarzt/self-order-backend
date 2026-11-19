package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"selforder.pramana.net/internal/validator"
)

type PesananStatus string

const (
	PesananPendingPayment PesananStatus = "pending_payment"
	PesananPaid           PesananStatus = "paid"
	PesananCooking        PesananStatus = "cooking"
	PesananServed         PesananStatus = "served"
	PesananConfirmed      PesananStatus = "confirmed"
)

type Pesanan struct {
	ID          int64         `json:"id"`
	NoMeja      string        `json:"no_meja"`
	Status      PesananStatus `json:"status"`
	IDTransaksi *string       `json:"id_transaksi"`
	TotalAmount float64       `json:"total_amount"`
	NamaPemesan *string       `json:"nama_pemesan"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Catatan     *string       `json:"catatan"`
}

type PesananProduk struct {
	ID             int64           `json:"id"`
	ProdukID       int64           `json:"prooduk_id"`
	Qty            uint16          `json:"qty"`
	HargaSaatPesan decimal.Decimal `json:"harga_saat_pesan"`
	Subtotal       decimal.Decimal `json:"subtotal"`
	Catatan        string          `json:"catatan"`
	Nama           string          `json:"nama"`
	Slug           string          `json:"slug"`
}

type PesananProdukInput struct {
	ProdukID *int64  `json:"produk_id"`
	Catatan  *string `json:"catatan"`
	Quantity *int64  `json:"quantity"`
}

func ValidatePesananProduk(v *validator.Validator, input *PesananProdukInput) {
	v.Check(input.ProdukID != nil, "produk_id", "harus diisi")
	v.Check(input.ProdukID == nil || *input.ProdukID > 0, "produk_id", "harus lebih besar dari 0")

	v.Check(input.Quantity != nil, "quantity", "harus diisi")
	v.Check(input.Quantity == nil || *input.Quantity > 0, "quantity", "harus lebih besar dari 0")

	if input.Catatan != nil {
		v.Check(len(*input.Catatan) <= 1000, "catatan", "maksimal 1000 karakter")
	}
}

type PesananModel struct {
	DB *sql.DB
}

func (m PesananModel) GetPesananProdukByIDTransaksi(token string) (*[]PesananProduk, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id
		FROM pesanan
		WHERE id_transaksi = $1
	`

	var pesanan Pesanan

	err := m.DB.QueryRowContext(ctx, query, token).Scan(
		&pesanan.ID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &[]PesananProduk{}, nil
	}

	if err != nil {
		return nil, err
	}

	query = `
		SELECT pp.id, pp.produk_id, pp.qty, pp.catatan, pp.harga_saat_pesan, pp.subtotal, p.nama, p.slug
		FROM pesanan_produk pp
		JOIN produk p ON p.id = pp.produk_id
		WHERE pp.pesanan_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, pesanan.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pesananProduk := make([]PesananProduk, 0)

	for rows.Next() {
		var produk PesananProduk

		err := rows.Scan(
			&produk.ID,
			&produk.ProdukID,
			&produk.Qty,
			&produk.Catatan,
			&produk.HargaSaatPesan,
			&produk.Subtotal,
			&produk.Nama,
			&produk.Slug,
		)
		if err != nil {
			return nil, err
		}

		pesananProduk = append(pesananProduk, produk)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pesananProduk, nil
}

func (m PesananModel) GetPesananByIDTransaksi(token string) (*Pesanan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id
		FROM pesanan
		WHERE id_transaksi = $1
	`

	var pesanan Pesanan

	err := m.DB.QueryRowContext(ctx, query, token).Scan(
		&pesanan.ID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &pesanan, nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m PesananModel) DeletePesananProduk(produkId int64, pesananId int64) error {
	if produkId < 1 {
		return ErrRecordNotFound
	}

	if pesananId < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM pesanan_produk
	WHERE produk_id = $1 AND pesanan_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, produkId, pesananId)
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
