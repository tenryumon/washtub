package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/models"
)

type AuthorizationUC interface {
	DoLogin(ctx context.Context, param models.DoLoginReq) (resp models.DoLoginResp, err error)
	DoLogout(ctx context.Context, token string) models.DoLoginResp
	CheckLogin(ctx context.Context, param models.CheckLoginReq) (resp models.CheckLoginResp, err error)
	ExtendLogin(ctx context.Context, token string, rememberMe bool, userID int64) (resp models.DoLoginResp)

	GetRolesByUserID(ctx context.Context, userID int64) (result []string, err error)

	ShowEmailVerification(ctx context.Context, param models.EmailVerifReq) models.EmailVerifResp
	SubmitEmailVerification(ctx context.Context, param models.EmailVerifReq) models.EmailVerifResp

	SendForgotPassword(ctx context.Context, email string) (resp models.LoginForgotPasswordResp)
	NewPassword(ctx context.Context, param models.NewPasswordReq) (resp models.NewPasswordResp)

	DoLoginFromApp(ctx context.Context, param models.DoLoginReq) (resp models.DoLoginFromAppResp, err error)
	DoRefreshAccessToken(ctx context.Context, param models.DoRefreshAccessTokenReq) (resp models.DoLoginFromAppResp, err error)
}
