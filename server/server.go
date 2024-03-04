package server

import (
	"fmt"
	"github.com/devstik0407/shenron/store"
	"github.com/gorilla/mux"
	"net/http"
)

func Start(cfg *Config, deps Dependencies) {
	r := mux.NewRouter()
	r.HandleFunc("/store/{key}", GetItemHandler(deps.Store)).Methods(http.MethodGet)
	r.HandleFunc("/store", SetItemHandler(deps.Store)).Methods(http.MethodPost)

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)

	fmt.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}

type Dependencies struct {
	Store store.Store
}

type Config struct {
	Address string
	Port    int
}
