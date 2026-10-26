package main

import (
	"errors"
	"fmt"
	"net/http"

	"selforder.pramana.net/internal/data"
	"selforder.pramana.net/internal/validator"
)

func (app *application) showProdukHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	produk, err := app.models.Produk.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"produk": produk}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createProdukHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nama      string `json:"nama"`
		Deskripsi string `json:"deskripsi"`
		GrupID    int64  `json:"grup_id"`
		Harga     uint32 `json:"harga"`
		Stok      uint16 `json:"stok"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	produk := &data.Produk{
		Nama:      input.Nama,
		Slug:      app.generateSlug(input.Nama),
		Deskripsi: input.Deskripsi,
		GrupID:    input.GrupID,
		Harga:     input.Harga,
		Stok:      input.Stok,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	if data.ValidateProduk(v, produk); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Produk.Insert(produk)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/produk/%d", produk.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"produk": produk}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProdukHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	produk, err := app.models.Produk.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Nama      *string `json:"nama"`
		Deskripsi *string `json:"deskripsi"`
		GrupID    *int64  `json:"grup_id"`
		Harga     *uint32 `json:"harga"`
		Stok      *uint16 `json:"stok"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Nama != nil {
		produk.Nama = *input.Nama
		produk.Slug = app.generateSlug(*input.Nama)
	}

	if input.Deskripsi != nil {
		produk.Deskripsi = *input.Deskripsi
	}

	if input.GrupID != nil {
		produk.GrupID = *input.GrupID
	}

	if input.Harga != nil {
		produk.Harga = *input.Harga
	}

	if input.Stok != nil {
		produk.Stok = *input.Stok
	}

	v := validator.New()
	if data.ValidateProduk(v, produk); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Produk.Update(produk)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"produk": produk}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteProdukHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Produk.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "produk successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listProdukHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nama string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Nama = app.readString(qs, "nama", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "nama", "-id", "-nama"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
