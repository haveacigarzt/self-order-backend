package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
	"selforder.pramana.net/internal/data"
	"selforder.pramana.net/internal/validator"
)

// func (app *application) showProdukHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	produk, err := app.models.Produk.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	// Encode the struct to JSON and send it as the HTTP response.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"produk": produk}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}

// }

func (app *application) showProdukHandler(w http.ResponseWriter, r *http.Request) {
	slug, err := app.readStringParam(r, "slug")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	produk, err := app.models.Produk.Get(slug)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	gambar, err := app.models.Gambar.GetAllByProdukID(produk.ID)

	// Inisialisasi 5 slot gambar
	produk.Gambar = make([]data.ProdukGambar, 5)

	for i := range produk.Gambar {
		produk.Gambar[i].Urutan = int16(i + 1)
	}

	// Isi gambar yang tersedia dari database
	for _, g := range gambar {
		if g.Urutan < 1 || g.Urutan > 5 {
			continue
		}

		produk.Gambar[g.Urutan-1].FilePath = g.FilePath
	} // Inisialisasi 5 slot gambar
	produk.Gambar = make([]data.ProdukGambar, 5)

	for i := range produk.Gambar {
		produk.Gambar[i].Urutan = int16(i + 1)
	}

	// Isi gambar yang tersedia dari database
	for _, g := range gambar {
		if g.Urutan < 1 || g.Urutan > 5 {
			continue
		}

		produk.Gambar[g.Urutan-1].FilePath = g.FilePath
	}

	viewData := map[string]any{
		"Produk": produk,
	}

	result, err := json.MarshalIndent(viewData, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(result))

	err = app.render(w, r, http.StatusOK, "single-produk", viewData)
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

	produk, err := app.models.Produk.GetById(id)
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
	nomorMeja, err := app.readNomorMeja(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

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

	// Ambil pesanan
	var pesananProduk *[]data.PesananProduk
	var total decimal.Decimal
	var totalQty uint16
	produkQty := make(map[int64]uint16)
	produkCatatan := make(map[int64]string)

	cookie, err := r.Cookie("transaction_id")
	if err == nil && cookie.Value != "" {
		pesananProduk, err = app.models.Pesanan.GetPesananProdukByIDTransaksi(cookie.Value)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		for _, item := range *pesananProduk {
			total = total.Add(item.Subtotal)
			totalQty += item.Produk.Qty
			produkQty[item.Produk.ID] = item.Produk.Qty
			produkCatatan[item.Produk.ID] = item.Produk.Catatan
		}
	} else {
		pesananProduk = &[]data.PesananProduk{}
	}

	groups := make([]data.ProdukGroup, 0)

	groupIndex := make(map[int64]int)

	for _, produk := range produkList {
		index, exists := groupIndex[produk.GrupID]
		produk.Qty = produkQty[produk.ID]
		produk.Catatan = produkCatatan[produk.ID]

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
		"Groups":        groups,
		"Metadata":      metadata,
		"NomorMeja":     nomorMeja,
		"PesananProduk": pesananProduk,
		"Total":         formatRupiah(total),
		"TotalQty":      totalQty,
	}

	result, err := json.MarshalIndent(pesananProduk, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(result))

	templateName := "produk"

	err = app.render(w, r, http.StatusOK, templateName, viewData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateGambarProdukHandler(w http.ResponseWriter, r *http.Request) {
	idProduk, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	urutan, err := app.readInt16Param(r, "urutan")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	urutan += 1

	if urutan < 1 || urutan > 5 {
		app.badRequestResponse(w, r, errors.New("urutan must be between 1 and 5"))
		return
	}

	//upload gambar
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	file, header, err := r.FormFile("gambar")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	// Validasi ukuran
	if header.Size > 5<<20 {
		app.badRequestResponse(w, r, errors.New("image exceeds 5 MB"))
		return
	}

	contentType, err := validateImage(file, header)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	filename := uuid.New().String()

	fmt.Println("Produk Id:", idProduk)
	fmt.Println("UU Id:", filename)
	fmt.Println("Content Type:", contentType)

	oldPath, err := app.models.Gambar.GetByProdukIDAndUrutan(idProduk, urutan)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Println("Old path:", oldPath)

	if oldPath != "" {
		err := os.Remove(filepath.Join("internal/static/uploads", oldPath))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			app.logger.Error("failed to delete old image", "error", err)
		}
	}

	// Ambil dimensi gambar
	img, format, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("Decode error:", err)
		app.badRequestResponse(w, r, errors.New("invalid image file"))
		return
	}

	fmt.Println("Format:", format)
	fmt.Println("Width:", img.Width)
	fmt.Println("Height:", img.Height)

	if img.Width > math.MaxInt16 || img.Height > math.MaxInt16 {
		app.badRequestResponse(w, r, errors.New("image dimensions are too large"))
		return
	}

	width := int16(img.Width)
	height := int16(img.Height)

	fmt.Println("width:", width)

	// Kembalikan posisi file ke awal
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 6. Simpan gambar baru ke local storage
	ext := filepath.Ext(header.Filename)
	newPath := filepath.Join("produk", filename+ext)

	fullPath := filepath.Join("internal/static/uploads", newPath)

	err = os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Simpan file baru
	dst, err := os.Create(fullPath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_, copyErr := io.Copy(dst, file)
	closeErr := dst.Close()

	if copyErr != nil {
		os.Remove(fullPath)
		app.serverErrorResponse(w, r, copyErr)
		return
	}

	if closeErr != nil {
		os.Remove(fullPath)
		app.serverErrorResponse(w, r, closeErr)
		return
	}

	// 7. UPSERT data gambar ke PostgreSQL
	gambar := &data.Gambar{
		ProdukID: idProduk,
		FilePath: newPath,
		Width: sql.NullInt16{
			Int16: width,
			Valid: true,
		},
		Height: sql.NullInt16{
			Int16: height,
			Valid: true,
		},
		Urutan:      urutan,
		IsThumbnail: urutan == 1,
	}

	err = app.models.Gambar.Upsert(gambar)
	if err != nil {
		os.Remove(fullPath) // Bersihkan file baru jika database gagal
		app.serverErrorResponse(w, r, err)
		return
	}

	// 8. Hapus gambar lama setelah UPSERT berhasil
	if oldPath != "" && oldPath != newPath {
		err := os.Remove(filepath.Join("internal/static/uploads", oldPath))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			app.logger.Error("failed to delete old image", "error", err)
		}
	}
}
