package main

import (
	"fmt"
	"net/http"
	"time"

	"selforder.pramana.net/internal/data"
	"selforder.pramana.net/internal/validator"
)

func (app *application) showProdukHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	produk := data.Produk{
		ID:        id,
		GroupID:   102,
		Nama:      "Americano",
		Slug:      "americano",
		Deskripsi: "test aja",
		Harga:     20000,
		// Stok:      1000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"produk": produk}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createProdukHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GroupID   int64  `json:"group_id"`
		Nama      string `json:"nama"`
		Slug      string `json:"slug"`
		Deskripsi string `json:"deskripsi"`
		Harga     uint32 `json:"harga"`
		Stok      uint16 `json:"stok"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	produk := &data.Produk{
		GroupID:   input.GroupID,
		Nama:      input.Nama,
		Slug:      input.Slug,
		Deskripsi: input.Deskripsi,
		Harga:     input.Harga,
		Stok:      input.Stok,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	if data.ValidateProduk(v, produk); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
