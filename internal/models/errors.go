package models

import (
	"fmt"
	"strconv"

	errsPkg "github.com/nyelonong/boilerplate-go/pkg/errors"
)

type Errors struct {
	StatusCode       string
	ValidationErrors map[string][]string
	ErrorMessage     string
}

func NewErrors(err error) (modelErrs Errors) {
	statusCode := errsPkg.HTTPCode(err).String()
	modelErrs = Errors{}
	modelErrs.InitError()
	modelErrs.StatusCode = statusCode
	if statusCode != "500000" {
		modelErrs.ErrorMessage = err.Error()
	}
	modelErrs.ValidationErrors = errsPkg.Wrap(err).Forms().ToMapStringArray()
	return
}

func (e *Errors) Success() {
	e.StatusCode = "200000"
	e.ErrorMessage = "Success"
}

func (e *Errors) Created() {
	e.StatusCode = "201000"
	e.ErrorMessage = "Created"
}

func (e *Errors) SuccessWithMessage(msg string) {
	e.StatusCode = "200000"
	e.ErrorMessage = msg
}

func (e *Errors) ErrorInvalidForm(key string, errors []string) {
	e.StatusCode = "400001"
	e.ErrorMessage = "Maaf, terdapat kesalahan dalam form. Harap melakukan pengecekan ulang."

	// Create ValidationErrors if still nil
	if e.ValidationErrors == nil {
		e.ValidationErrors = map[string][]string{}
	}

	// Append Errors to keys if already exist
	if _, ok := e.ValidationErrors[key]; ok {
		e.ValidationErrors[key] = append(e.ValidationErrors[key], errors...)
	} else {
		e.ValidationErrors[key] = errors
	}
}

func (e *Errors) HaveInvalidForm() bool {
	if len(e.ValidationErrors) > 0 {
		return true
	}
	return false
}

func (e *Errors) ErrorNotLogin() {
	e.StatusCode = "401000"
	e.ErrorMessage = "Silakan login terlebih dahulu."
}

func (e *Errors) ErrorAlreadyLogin() {
	e.StatusCode = "401001"
	e.ErrorMessage = "Pengguna tidak valid."
}

func (e *Errors) ErrorInvalidToken(msg string) {
	e.StatusCode = "401002"
	e.ErrorMessage = fmt.Sprintf("Token tidak valid. %s", msg)
}

func (e *Errors) ErrorNoAccess() {
	e.StatusCode = "403000"
	e.ErrorMessage = "Maaf, pengguna tidak memiliki izin untuk mengakses halaman ini."
}

func (e *Errors) ErrorNotFound() {
	e.StatusCode = "404000"
	e.ErrorMessage = "Maaf, Item tidak ditemukan."
}

func (e *Errors) InitError() {
	e.StatusCode = "500000"
	e.ErrorMessage = "Maaf, sistem mengalami kegagalan."
}

func (e *Errors) ErrorGeneral(msg string) {
	e.StatusCode = "400000"
	e.ErrorMessage = msg
}

func (e *Errors) ErrorCustom(code, msg string) {
	e.StatusCode = code
	e.ErrorMessage = msg
}

func (e *Errors) GetError() Errors {
	if e.StatusCode == "" {
		e.InitError()
	}
	return *e
}

func (e *Errors) IsSuccess() bool {
	return e.GetStatusCode() == 200
}

func (e *Errors) GetStatusCode() int {
	if e.StatusCode == "" {
		return 500
	}

	code := e.StatusCode[0:3]
	statusCode, err := strconv.Atoi(code)
	if err != nil {
		return 500
	}
	return statusCode
}
