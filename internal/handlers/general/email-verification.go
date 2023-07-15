package general

import (
	"net/http"

	"github.com/nyelonong/boilerplate-go/internal/models"
)

type EmailVerifReq struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (sd *GeneralView) ShowEmailVerification(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()

	defer func() {
		writeResponseJSON(w, nil, reqErr)
	}()

	request := EmailVerifReq{}
	err := getBody(r, &request)
	if err != nil {
		return
	}

	result := sd.authorizationUC.ShowEmailVerification(ctx, models.EmailVerifReq{
		Token: request.Token,
		Email: request.Email,
	})

	reqErr = result.GetError()
	return
}

func (sd *GeneralView) SubmitEmailVerification(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()

	defer func() {
		writeResponseJSON(w, nil, reqErr)
	}()

	request := EmailVerifReq{}
	err := getBody(r, &request)
	if err != nil {
		return
	}

	result := sd.authorizationUC.SubmitEmailVerification(ctx, models.EmailVerifReq{
		Token:    request.Token,
		Email:    request.Email,
		Password: request.Password,
	})

	reqErr = result.GetError()
	return
}
