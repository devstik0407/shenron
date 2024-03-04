package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []Error     `json:"errors"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
		Errors:  nil,
	}
}

func NewErrorResponse(code, message, reason string) *Response {
	return &Response{
		Success: false,
		Data:    nil,
		Errors: []Error{
			{
				Code:    code,
				Message: message,
				Reason:  reason,
			},
		},
	}
}

func (r *Response) Write(w http.ResponseWriter, status int) {
	body, err := json.Marshal(r)
	if err != nil {
		//log here
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err = w.Write(body); err != nil {
		//log here
	}
}
