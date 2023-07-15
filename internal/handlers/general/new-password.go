package general

import (
	"net/http"

	"github.com/nyelonong/boilerplate-go/internal/models"
)

type SendForgotPasswordReq struct {
	Email string `json:"email"`
}

func (sd *GeneralView) SendForgotPassword(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()

	defer func() {
		writeResponseJSON(w, nil, reqErr)
	}()

	request := SendForgotPasswordReq{}
	err := getBody(r, &request)
	if err != nil {
		return
	}

	result := sd.authorizationUC.SendForgotPassword(ctx, request.Email)

	reqErr = result.GetError()
	return
}

type NewPasswordReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (sd *GeneralView) NewPassword(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	var result models.NewPasswordResp
	ctx := r.Context()

	defer func() {
		writeResponseJSON(w, result, reqErr)
	}()

	request := NewPasswordReq{}
	err := getBody(r, &request)
	if err != nil {
		return
	}

	result = sd.authorizationUC.NewPassword(ctx, models.NewPasswordReq{
		Token:    request.Token,
		Password: request.Password,
	})

	reqErr = result.GetError()
	return
}
