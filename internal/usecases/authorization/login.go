package authorization

import (
	"context"

	"github.com/nyelonong/boilerplate-go/core/jwt"
	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/core/validate"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

func validateEmail(email string) []string {
	return validate.String("Email", email, validate.GetStringCommon(models.ValidRequiredEmail)...)
}

func validatePassword(password string) []string {
	return validate.String("Password", password, validate.GetStringCommon(models.ValidRequiredPassword)...)
}

// DoLogin check email and password of a user, create a session in redis and return it to cookie
func (au *Authorization) DoLogin(ctx context.Context, param models.DoLoginReq) (resp models.DoLoginResp, err error) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase":  "DoLogin",
		"email":    param.Email,
		"remember": param.RememberMe,
	}

	// Validate User Input Parameter
	if failures := validateEmail(param.Email); len(failures) > 0 {
		resp.ErrorInvalidForm("email", failures)
	}
	if failures := validatePassword(param.Password); len(failures) > 0 {
		resp.ErrorInvalidForm("password", failures)
	}
	if resp.HaveInvalidForm() {
		log.Debugf("Failed to do Login because invalid form")
		return resp, nil
	}

	// Verify User Input Parameter
	user, err := au.user.GetUserByEmail(ctx, param.Email, param.UserType)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do GetUserByEmail because %s", err)
		return resp, nil
	}
	log.Debugf("User: %+v", user)
	if !user.Exist() {
		resp.ErrorGeneral(entities.MessageFailedLogin)
		return resp, nil
	}
	if !user.IsActive() {
		resp.ErrorGeneral(entities.MessageUserNotActive)
		return resp, nil
	}

	valid, exist, err := au.user.CheckPassword(ctx, user.ID, param.Password)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do CheckPassword because %s", err)
		return resp, nil
	}
	log.Debugf("Valid Password: %v. Exist Password: %v", valid, exist)
	if !exist {
		resp.ErrorGeneral(entities.MessagePasswordNotSet)
		return resp, nil
	}
	if !valid {
		resp.ErrorGeneral(entities.MessageFailedLogin)
		return resp, nil
	}

	// Start Usecase Action
	session, err := au.session.CreateSession(ctx, user, param.RememberMe)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do CreateSession because %s", err)
		return resp, nil
	}

	resp.Token = session.Token
	resp.ExpireAt = session.ExpireAt
	resp.Name = user.Name
	resp.UserID = user.ID
	resp.Success()

	// update user last action
	errLog := au.user.UpdateLastAction(ctx, user.ID)
	if errLog != nil {
		log.ErrorWithField(ucContext, "Failed to do UpdateLastAction because %s", err)
	}

	return resp, nil
}

// DoLogout will remove session from redis. This will always return success response
func (au *Authorization) DoLogout(ctx context.Context, sessionToken string) (resp models.DoLoginResp) {
	// Initialize Usecase Parameter
	resp.Success()
	ucContext := map[string]interface{}{
		"usecase": "DoLogout",
		"session": sessionToken,
	}

	// Validate User Input Parameter
	if sessionToken == "" {
		return resp
	}

	// Start Usecase Action
	err := au.session.DeleteSession(ctx, sessionToken)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do DeleteSession because %s", err)
	}

	return resp
}

// CheckLogin will check user info based session token from cookie
func (au *Authorization) CheckLogin(ctx context.Context, param models.CheckLoginReq) (resp models.CheckLoginResp, err error) {
	// Initialize Usecase Parameter
	// ucContext := map[string]interface{}{
	// 	"usecase": "CheckLogin",
	// 	"session": param.Token,
	// }

	// Validate User Input Parameter
	if param.Token == "" {
		return resp, nil
	}

	// Start Usecase Action
	var userID int64
	switch param.LoginDevice {
	case entities.LoginDeviceDesktop:
		userID = au.validateWithSession(ctx, param.Token)
	case entities.LoginDeviceMobileApp:
		userID = au.validateWithJWT(ctx, param.Token)
	default:
		userID = au.validateWithSession(ctx, param.Token)
	}

	if userID == 0 {
		return resp, nil
	}
	log.Debugf("UserID: %+v", userID)

	resp.UserID = userID

	return resp, nil
}

func (au *Authorization) ExtendLogin(ctx context.Context, token string, rememberMe bool, userID int64) (resp models.DoLoginResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase":  "ExtendLogin",
		"token":    token,
		"remember": rememberMe,
	}

	// Validate User Input Parameter
	if token == "" {
		return
	}

	session, err := au.session.ExtendSession(ctx, token, rememberMe)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to ExtendSession because %s", err)
		return
	}

	resp.Token = session.Token
	resp.ExpireAt = session.ExpireAt
	resp.Success()

	// update user last action
	errLog := au.user.UpdateLastAction(ctx, userID)
	if errLog != nil {
		log.ErrorWithField(ucContext, "Failed to do UpdateLastAction because %s", err)
	}

	return
}

// DoLoginFromApp check email and password of a user, create a jwt token and return it
func (au *Authorization) DoLoginFromApp(ctx context.Context, param models.DoLoginReq) (resp models.DoLoginFromAppResp, err error) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase":  "DoLoginApp",
		"email":    param.Email,
		"remember": param.RememberMe,
	}

	// Validate User Input Parameter
	if failures := validateEmail(param.Email); len(failures) > 0 {
		resp.ErrorInvalidForm("email", failures)
	}
	if failures := validatePassword(param.Password); len(failures) > 0 {
		resp.ErrorInvalidForm("password", failures)
	}
	if resp.HaveInvalidForm() {
		log.Debugf("Failed to do Login because invalid form")
		return resp, nil
	}

	// Verify User Input Parameter
	user, err := au.user.GetUserByEmail(ctx, param.Email, param.UserType)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do GetUserByEmail because %s", err)
		return resp, nil
	}
	log.Debugf("User: %+v", user)
	if !user.Exist() {
		resp.ErrorGeneral(entities.MessageFailedLogin)
		return resp, nil
	}
	if !user.IsActive() {
		resp.ErrorGeneral(entities.MessageUserNotActive)
		return resp, nil
	}

	valid, exist, err := au.user.CheckPassword(ctx, user.ID, param.Password)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do CheckPassword because %s", err)
		return resp, nil
	}
	log.Debugf("Valid Password: %v. Exist Password: %v", valid, exist)
	if !exist {
		resp.ErrorGeneral(entities.MessagePasswordNotSet)
		return resp, nil
	}
	if !valid {
		resp.ErrorGeneral(entities.MessageFailedLogin)
		return resp, nil
	}

	// Start Usecase Action
	tokenStr, err := jwt.CreateAccessToken(user.ID, user.Name)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do CreateToken JWT because %s", err)
	}
	log.Debugf("Token for UserID %d is %s", user.ID, tokenStr)
	resp.AccessToken = tokenStr

	if param.RememberMe {
		// Create Refresh Token
		refreshToken, err := jwt.CreateRefreshToken(user.ID, user.Name)
		if err != nil {
			log.ErrorWithField(ucContext, "Failed to do CreateRefreshToken JWT because %s", err)
		}
		log.Debugf("Refresh Token for UserID %d is %s", user.ID, refreshToken)
		resp.RefreshToken = refreshToken
	}
	resp.Success()

	// update user last action
	errLog := au.user.UpdateLastAction(ctx, user.ID)
	if errLog != nil {
		log.ErrorWithField(ucContext, "Failed to do UpdateLastAction because %s", err)
	}

	return resp, nil
}

func (au *Authorization) validateWithSession(ctx context.Context, token string) int64 {
	ucContext := map[string]interface{}{
		"usecase":    "CheckLogin",
		"validation": "session",
		"token":      token,
	}

	userID, err := au.session.GetSession(ctx, token)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to GetSession because %s", err)
		return 0
	}

	return userID
}

func (au *Authorization) validateWithJWT(ctx context.Context, token string) int64 {
	ucContext := map[string]interface{}{
		"usecase":    "CheckLogin",
		"validation": "jwt",
		"token":      token,
	}

	valid, err := jwt.IsValidAccessToken(token)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to validate token signature %s because %s", token, err)
		return 0
	}

	if !valid {
		return 0
	}

	// Get User Info
	claims, err := jwt.UnmarshalPayload(token)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to unmarshal token payload because %s", err)
		return 0
	}

	return claims.UserID
}

// DoRefreshAccessToken create a new access token if refresh token still valid
func (au *Authorization) DoRefreshAccessToken(ctx context.Context, param models.DoRefreshAccessTokenReq) (resp models.DoLoginFromAppResp, err error) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase":       "DoRefreshAccessToken",
		"refresh_token": param.RefreshToken,
	}

	// Validate User Input Parameter
	if param.RefreshToken == "" {
		resp.ErrorInvalidForm("refresh token", []string{"refresh token harus diisi"})
		return resp, nil
	}

	// Validate Token Signature
	valid, err := jwt.IsValidRefreshToken(param.RefreshToken)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to validate refresh token signature because %s", err)
		resp.ErrorInvalidToken(err.Error())
		return resp, nil
	}

	if !valid {
		resp.ErrorInvalidToken("")
		return resp, nil
	}

	// Get User Info
	claims, err := jwt.UnmarshalPayload(param.RefreshToken)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to unmarshal token payload because %s", err)
		resp.ErrorInvalidToken(err.Error())
		return resp, nil
	}

	if claims.UserID == 0 {
		resp.ErrorInvalidToken("")
		return resp, nil
	}
	log.Debugf("UserID: %d", claims.UserID)

	user, err := au.user.GetUserByID(ctx, claims.UserID)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to GetUserByID because %s", err)
		return resp, nil
	}

	// Start Usecase Action
	tokenStr, err := jwt.CreateAccessToken(user.ID, user.Name)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do CreateToken JWT because %s", err)
		resp.ErrorGeneral(err.Error())
		return resp, nil
	}
	log.Debugf("New access token for UserID %d is %s", user.ID, tokenStr)

	refreshToken, err := jwt.CreateRefreshToken(user.ID, user.Name)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do CreateRefreshToken JWT because %s", err)
		resp.ErrorGeneral(err.Error())
		return resp, nil
	}
	log.Debugf("New refresh Token for UserID %d is %s", user.ID, refreshToken)

	resp.AccessToken = tokenStr
	resp.RefreshToken = refreshToken
	resp.Success()

	return resp, nil
}
