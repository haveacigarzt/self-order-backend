package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.Handler(
		http.MethodGet,
		"/static/*filepath",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./internal/static")),
		),
	)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/grupproduk", app.createGrupProdukHandler)

	router.HandlerFunc(http.MethodPost, "/v1/produk", app.requireActivatedUser(app.createProdukHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/produk/:id", app.requirePermission("produk:read", app.requireActivatedUser(app.showProdukHandler)))
	router.HandlerFunc(http.MethodPatch, "/v1/produk/:id", app.requireActivatedUser(app.updateProdukHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/produk/:id", app.requireActivatedUser(app.deleteProdukHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/produk", app.requireActivatedUser(app.listProdukHandler))

	router.HandlerFunc(http.MethodGet, "/meja/:nomor_meja", app.listProdukHandler)

	router.HandlerFunc(http.MethodPost, "/meja/:nomor_meja", app.addPesananHandler)
	router.HandlerFunc(http.MethodDelete, "/meja/:nomor_meja/:id", app.removePesananHandler)

	router.HandlerFunc(http.MethodGet, "/admin", app.adminMainHandler)
	router.HandlerFunc(http.MethodGet, "/admin/produk/:slug", app.showProdukHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
