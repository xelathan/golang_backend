package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xelathan/golang_backend/services/cart"
	"github.com/xelathan/golang_backend/services/order"
	"github.com/xelathan/golang_backend/services/product"
	"github.com/xelathan/golang_backend/services/user"
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
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subRouter)

	orderStore := order.NewStore(s.db)
	orderHandler := order.NewHandler(userStore, orderStore, productStore, s.db)
	orderHandler.RegisterRoutes(subRouter)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore, s.db)
	cartHandler.RegisterRoutes(subRouter)

	subRouter.HandleFunc("/", handleHome).Methods("GET")

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

func handleHome(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Success!")
}
