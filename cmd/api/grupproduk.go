package main

import (
	"fmt"
	"net/http"

	"selforder.pramana.net/internal/data"
	"selforder.pramana.net/internal/validator"
)

func (app *application) createGrupProdukHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nama string `json:"nama"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	grupProduk := &data.GrupProduk{
		Nama: input.Nama,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	if data.ValidateGrupProduk(v, grupProduk); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.GrupProduk.Insert(grupProduk)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/grupproduk/%d", grupProduk.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"grup_produk": grupProduk}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
