package view

import (
	"encoding/json"
	"errors"
	"github.com/rithikjain/SocialRecipe/pkg"
	"net/http"
)

type ErrView struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

//noinspection ALL
var (
	ErrMethodNotAllowed = errors.New("Error: Method is not allowed")
	ErrInvalidToken     = errors.New("Error: Invalid Authorization token")
	ErrUserExists       = errors.New("Error: User already exists")
	ErrFile             = errors.New("Error: Something wrong with file")
	ErrUpload           = errors.New("Error: Upload failed")
)

var ErrHTTPStatusMap = map[string]int{
	pkg.ErrNotFound.Error():     http.StatusNotFound,
	pkg.ErrInvalidSlug.Error():  http.StatusBadRequest,
	pkg.ErrExists.Error():       http.StatusConflict,
	pkg.ErrNoContent.Error():    http.StatusNotFound,
	pkg.ErrDatabase.Error():     http.StatusInternalServerError,
	pkg.ErrUnauthorized.Error(): http.StatusUnauthorized,
	pkg.ErrForbidden.Error():    http.StatusForbidden,
	pkg.ErrEmail.Error():        http.StatusBadRequest,
	pkg.ErrPassword.Error():     http.StatusBadRequest,
	ErrMethodNotAllowed.Error(): http.StatusMethodNotAllowed,
	ErrInvalidToken.Error():     http.StatusBadRequest,
	ErrUserExists.Error():       http.StatusConflict,
	ErrFile.Error():             http.StatusBadRequest,
}

func Wrap(err error, w http.ResponseWriter) {
	msg := err.Error()
	code := ErrHTTPStatusMap[msg]

	if code == 0 {
		code = http.StatusInternalServerError
	}

	errView := ErrView{
		Message: msg,
		Status:  code,
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errView)
}
