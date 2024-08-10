package server

import (
	"fmt"
	"github.com/devstik0407/shenron/service"
	"github.com/gorilla/mux"
	"net/http"
)

func Start(cfg *Config, deps Dependencies) {
	r := mux.NewRouter()
	r.HandleFunc("/store/{key}", GetItemHandler(deps.Service)).Methods(http.MethodGet)
	r.HandleFunc("/store", SetItemHandler(deps.Service)).Methods(http.MethodPost)

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)

	fmt.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}

type Dependencies struct {
	Service service.Service
}

type Config struct {
	Address string
	Port    int
}
