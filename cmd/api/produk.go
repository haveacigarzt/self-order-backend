package main

import (
	"fmt"
	"net/http"
	"time"

	"selforder.pramana.net/internal/data"
)

func (app *application) createProdukHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new Produk")
}

func (app *application) showProdukHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
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
	err = app.writeJSON(w, http.StatusOK, produk, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
