package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"selforder.pramana.net/internal/data"
	"selforder.pramana.net/internal/validator"
)

func (app *application) addPesananHandler(w http.ResponseWriter, r *http.Request) {
	nomorMeja, err := app.readNomorMeja(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input data.PesananProdukInput

	err = r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Println(r.PostForm)

	produkID, err := strconv.ParseInt(r.FormValue("produk_id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	quantity, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	catatan := r.FormValue("catatan")

	input.ProdukID = &produkID
	input.Quantity = &quantity
	input.Catatan = &catatan

	v := validator.New()

	data.ValidatePesananProduk(v, &input)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var token string
	newCookie := false

	cookie, err := r.Cookie("transaction_id")

	if err != nil && err != http.ErrNoCookie {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err == nil && cookie.Value != "" {
		token = cookie.Value
	} else {
		token, err = generateRandomString()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		newCookie = true
	}

	ctx := r.Context()

	tx, err := app.models.Pesanan.DB.BeginTx(ctx, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback()

	noMeja := r.PostForm.Get("no_meja")

	var pesananID int64

	// Cari atau buat pesanan.
	err = tx.QueryRowContext(ctx, `
    INSERT INTO pesanan (no_meja, id_transaksi)
    VALUES ($1, $2)
    ON CONFLICT (id_transaksi) DO NOTHING
    RETURNING id
`, noMeja, token).Scan(&pesananID)

	isNew := err == nil

	if err != nil && err != sql.ErrNoRows {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Jika pesanan sudah ada, ambil ID dan kunci row.
	if !isNew {
		var status string

		err = tx.QueryRowContext(ctx, `
        SELECT id, status
        FROM pesanan
        WHERE id_transaksi = $1
        FOR UPDATE
    `, token).Scan(&pesananID, &status)

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if status != "pending_payment" {
			app.badRequestResponse(w, r, errors.New("pesanan tidak dapat diubah"))
			return
		}
	}

	// Buat timeline awal.
	if isNew {
		_, err = tx.ExecContext(ctx, `
        INSERT INTO status_timeline (pesanan_id, status)
        VALUES ($1, 'pending_payment')
    `, pesananID)

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// Tambahkan produk atau increment quantity.
	var detailID int64

	err = tx.QueryRowContext(ctx, `
    INSERT INTO pesanan_produk (
        pesanan_id,
        produk_id,
        qty,
        catatan,
        harga_saat_pesan
    )
    SELECT
        $1,
        p.id,
        $3,
        $4,
        p.harga
    FROM produk p
    WHERE p.id = $2

    ON CONFLICT (pesanan_id, produk_id)
    DO UPDATE SET
        qty = EXCLUDED.qty,
        catatan = EXCLUDED.catatan
    RETURNING id
`, pesananID, *input.ProdukID, *input.Quantity, input.Catatan).Scan(&detailID)

	if err != nil {
		if err == sql.ErrNoRows {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	// Hitung ulang total pesanan.
	_, err = tx.ExecContext(ctx, `
    UPDATE pesanan
    SET
        total_amount = (
            SELECT COALESCE(SUM(subtotal), 0)
            FROM pesanan_produk
            WHERE pesanan_id = $1
        ),
        updated_at = NOW()
    WHERE id = $1
`, pesananID)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Simpan semua perubahan.
	if err := tx.Commit(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if newCookie {
		http.SetCookie(w, &http.Cookie{
			Name:     "transaction_id",
			Value:    token,
			Path:     "/",
			MaxAge:   18000,
			HttpOnly: true,
			Secure:   false, //true if HTTPS/production
			SameSite: http.SameSiteLaxMode,
		})
	}

	http.Redirect(w, r, fmt.Sprintf("/meja/%d", nomorMeja), http.StatusSeeOther)
}

func (app *application) removePesananHandler(w http.ResponseWriter, r *http.Request) {
	nomorMeja, err := app.readNomorMeja(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	produkId, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	cookie, err := r.Cookie("transaction_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	token := cookie.Value

	pesanan, err := app.models.Pesanan.GetPesananByIDTransaksi(token)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Pesanan.DeletePesananProduk(produkId, pesanan.ID)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// err = app.writeJSON(w, http.StatusOK, envelope{"message": "pesanan produk successfully deleted"}, nil)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// }

	http.Redirect(w, r, fmt.Sprintf("/meja/%d", nomorMeja), http.StatusSeeOther)

}
