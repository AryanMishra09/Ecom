package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AryanMishra09/Ecom/service/cart"
	"github.com/AryanMishra09/Ecom/service/order"
	"github.com/AryanMishra09/Ecom/service/product"
	"github.com/AryanMishra09/Ecom/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subRouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subRouter)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Printf("Listening on: %v", s.addr)

	return http.ListenAndServe(s.addr, router)
}
