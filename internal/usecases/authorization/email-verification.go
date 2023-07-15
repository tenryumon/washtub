package authorization

import (
	"context"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

// ShowEmailVerification will check whether combination of email and token is valid or not
func (au *Authorization) ShowEmailVerification(ctx context.Context, param models.EmailVerifReq) (resp models.EmailVerifResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "ShowEmailVerification",
		"email":   param.Email,
	}

	// Validate User Input Parameter
	if param.Email == "" {
		resp.ErrorCustom("400404", "Token tidak valid")
		return
	}
	if param.Token == "" {
		resp.ErrorCustom("400404", "Token tidak valid")
		return
	}

	// Start Usecase Action
	validPassReq, err := au.user.IsValidPasswordToken(ctx, param.Token, param.Email)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do IsValidPasswordToken because %s", err)
		return
	}
	if !validPassReq.Valid {
		if validPassReq.IsExpired() {
			resp.ErrorCustom("400444", "Token sudah kadaluarsa")
		} else {
			resp.ErrorCustom("400404", "Token tidak valid")
		}
		return
	}

	resp.Success()

	return
}

// SubmitEmailVerification will save password of a user based on combination of email and token
func (au *Authorization) SubmitEmailVerification(ctx context.Context, param models.EmailVerifReq) (resp models.EmailVerifResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "SubmitEmailVerification",
		"email":   param.Email,
	}

	// Validate User Input Parameter
	if passFail := validatePassword(param.Password); len(passFail) > 0 {
		resp.ErrorInvalidForm("password", passFail)
	}
	if resp.HaveInvalidForm() {
		log.Debugf("Failed to SubmitEmailVerification because invalid form")
		return
	}
	if param.Email == "" {
		resp.ErrorCustom("400404", "Token tidak valid")
		return
	}
	if param.Token == "" {
		resp.ErrorCustom("400404", "Token tidak valid")
		return
	}

	// Start Usecase Action
	success, isNewUser, validPassReq, err := au.user.SetPassword(ctx, param.Token, param.Email, param.Password)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do SetPassword because %s", err)
		return
	}
	log.Debugf("%+v", validPassReq)
	if !success && !validPassReq.Valid {
		if validPassReq.IsExpired() {
			resp.ErrorCustom("400444", "Token sudah kadaluarsa")
		} else {
			resp.ErrorCustom("400404", "Token tidak valid")
		}
		return
	}

	resp.IsNewUser = isNewUser
	resp.Success()

	return
}
