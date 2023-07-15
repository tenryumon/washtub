package authorization

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

// NewPassword will save new password for user
func (au *Authorization) NewPassword(ctx context.Context, param models.NewPasswordReq) (resp models.NewPasswordResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "NewPassword",
	}

	// Validate User Input Parameter
	if passFail := validatePassword(param.Password); len(passFail) > 0 {
		resp.ErrorInvalidForm("password", passFail)
	}
	if resp.HaveInvalidForm() {
		log.Debugf("Failed to NewPassword because invalid form")
		return
	}

	// Verify User Input Parameter
	tokenDecoded, err := base64.StdEncoding.DecodeString(param.Token)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do DecodeString token because %s", err)
		return
	}
	var token entities.NewPasswordToken
	err = json.Unmarshal(tokenDecoded, &token)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do Marshal token because %s", err)
		return
	}
	// validate token content
	if token.Email == "" {
		resp.ErrorCustom("400404", "Token tidak valid")
		return
	}
	if token.ID == "" {
		resp.ErrorCustom("400404", "Token tidak valid")
		return
	}

	// Start Usecase Action
	success, isNewUser, validPassReq, err := au.user.SetPassword(ctx, token.ID, token.Email, param.Password)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do SetPassword because %s", err)
		return
	}

	if !success && !validPassReq.Valid {
		if validPassReq.IsExpired() {
			resp.ErrorCustom("400444", "Halaman telah kedaluwarsa, Anda dapat meminta pengaturan ulang kata sandi kembali.")
		} else {
			resp.ErrorCustom("400404", "Halaman tidak valid, Anda dapat meminta pengaturan ulang kata sandi kembali.")
		}
		return
	}

	resp.IsNewUser = isNewUser
	resp.Success()

	return
}
