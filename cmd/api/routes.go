package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/grupproduk", app.createGrupProdukHandler)

	router.HandlerFunc(http.MethodPost, "/v1/produk", app.createProdukHandler)
	router.HandlerFunc(http.MethodGet, "/v1/produk/:id", app.showProdukHandler)

	return app.recoverPanic(router)
}
