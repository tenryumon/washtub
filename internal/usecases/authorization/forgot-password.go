package authorization

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nyelonong/boilerplate-go/core/environment"
	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

const (
	forgotPasswordTemplate = "forgot-password"
	forgotPasswordSubject  = "Atur ulang Kata Sandi Anda"
)

// SendForgotPassword will send user email to reset their password
func (au *Authorization) SendForgotPassword(ctx context.Context, email string) (resp models.LoginForgotPasswordResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "SendForgotPassword",
		"email":   email,
	}

	// Validate User Input Parameter
	if failures := validateEmail(email); len(failures) > 0 {
		resp.ErrorInvalidForm("email", failures)
	}

	// Verify User Input Parameter
	user, err := au.user.GetUserByEmail(ctx, email, entities.UserTypeStaff)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do GetUserByEmail because %s", err)
		return
	}
	if !user.Exist() {
		resp.ErrorGeneral("Alamat email tidak ditemukan.")
		return
	}
	if !user.IsActive() {
		resp.ErrorGeneral("Email belum diaktifkan, mohon periksa email terlebih dahulu.")
		return
	}
	// Start Usecase Action
	passwordRequest, err := au.user.CreateNewPasswordRequest(ctx, entities.PasswordRequest{
		UserID:      user.ID,
		Token:       uuid.New().String(),
		RequestType: entities.PasswordRequestTypeForgotPass,
		ExpiredTime: time.Now().Add(time.Minute * 20),
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	})
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do CreateNewPasswordRequest because %s", err)
		return
	}
	token := entities.NewPasswordToken{
		ID:    passwordRequest.Token,
		Email: email,
	}
	base64Token, err := token.GetBase64()
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do token GetBase64 because %s", err)
		return
	}
	forgotPasswordEndpoint := "%s/new-password?token=%s"
	metadata := map[string]interface{}{
		"name":       user.Name,
		"email":      email,
		"action_url": fmt.Sprintf(forgotPasswordEndpoint, environment.GetBaseDomain(), base64Token),
	}
	channel := entities.GetChannelEmailSupport(email, forgotPasswordSubject, forgotPasswordTemplate)

	err = au.notification.Send(ctx, channel, metadata)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do SendForgotPassword because %s", err)
		return
	}

	resp.Success()
	return
}
