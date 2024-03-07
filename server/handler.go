package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/devstik0407/shenron/pkg"
	"github.com/devstik0407/shenron/store"
	"github.com/gorilla/mux"
)

func GetItemHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		value, err := s.Get(key)
		if err != nil && errors.Is(err, store.ErrNotFound) {
			pkg.NewErrorResponse("shenron:not_found", "could not find given key", store.ErrNotFound.Error()).Write(w, http.StatusNotFound)
			return
		}
		if err != nil && errors.Is(err, store.ErrExpired) {
			pkg.NewErrorResponse("shenron:expired_key", "given key is expired", store.ErrExpired.Error()).Write(w, http.StatusBadRequest)
			return
		}
		if err != nil {
			pkg.NewErrorResponse("shenron:internal_error", "some server error occurred", err.Error()).Write(w, http.StatusInternalServerError)
			return
		}

		pkg.NewSuccessResponse(response{
			Key:   key,
			Value: value,
		}).Write(w, http.StatusOK)

		return
	}
}

func SetItemHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := setItemRequest{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			pkg.NewErrorResponse("shenron:invalid_request", "request body is invalid", err.Error()).Write(w, http.StatusBadRequest)
			return
		}

		if err := s.Set(requestBody.Key, requestBody.Value, time.Duration(requestBody.TTL)*time.Second); err != nil {
			pkg.NewErrorResponse("shenron:internal_error", fmt.Sprintf("failed to set value for key: %s", requestBody.Key), err.Error()).Write(w, http.StatusInternalServerError)
			return
		}

		pkg.NewSuccessResponse(response{
			Key:   requestBody.Key,
			Value: requestBody.Value,
		}).Write(w, http.StatusOK)

		return
	}
}

type response struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type setItemRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl_in_seconds"`
}
