package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"selforder.pramana.net/internal/data"
	"selforder.pramana.net/internal/validator"
)

func (app *application) adminMainHandler(w http.ResponseWriter, r *http.Request) {
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

	produkList, metadata, err := app.models.Produk.GetAllByGroup(input.Nama, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	groups := make([]data.ProdukGroup, 0)

	groupIndex := make(map[int64]int)

	for _, produk := range produkList {
		index, exists := groupIndex[produk.GrupID]

		if !exists {
			index = len(groups)

			groupIndex[produk.GrupID] = index

			groups = append(groups, data.ProdukGroup{
				GrupID:   produk.GrupID,
				GrupNama: produk.GrupNama,
				Produk:   make([]*data.Produk, 0),
			})
		}

		groups[index].Produk = append(groups[index].Produk, produk)
	}

	viewData := map[string]any{
		"Groups":   groups,
		"Metadata": metadata,
	}

	result, err := json.MarshalIndent(viewData, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(result))

	templateName := "admin-main"

	err = app.render(w, r, http.StatusOK, templateName, viewData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
