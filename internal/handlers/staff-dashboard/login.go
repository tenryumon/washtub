package staff_dashboard

import (
	"net/http"

	"github.com/nyelonong/boilerplate-go/internal/models"
)

type DoLoginReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type DoRefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (sd *StaffDashboard) DoLogin(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()
	start := sd.getTime(ctx)
	var result models.DoLoginResp

	defer func() {
		writeResponseJSON(w, result, reqErr)
		monitorAPI("DoLogin", start, reqErr)
	}()

	request := DoLoginReq{}
	err := getBody(r, &request)
	if err != nil {
		return
	}

	result, err = sd.authorizationUC.DoLogin(ctx, models.DoLoginReq{
		Email:      request.Email,
		Password:   request.Password,
		RememberMe: request.RememberMe,
		UserType:   sd.userType,
	})
	if err != nil {
		return
	}

	reqErr = result.GetError()
	if reqErr.IsSuccess() {
		http.SetCookie(w, &http.Cookie{
			Name:    sd.cookieName,
			Value:   sd.setSessionToken(result.Token, request.RememberMe),
			Path:    "/",
			Expires: result.ExpireAt,
		})
	}

	return
}

func (sd *StaffDashboard) DoLogout(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()
	start := sd.getTime(ctx)

	defer func() {
		writeResponseJSON(w, nil, reqErr)
		monitorAPI("DoLogout", start, reqErr)
	}()

	result := sd.authorizationUC.DoLogout(ctx, sd.getSessionToken(r))
	reqErr = result.GetError()

	http.SetCookie(w, &http.Cookie{
		Name:    sd.cookieName,
		Value:   "",
		Path:    "/",
		Expires: result.ExpireAt,
		MaxAge:  -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    sd.autoExtendCookie,
		Value:   "",
		Path:    "/",
		Expires: result.ExpireAt,
		MaxAge:  -1,
	})

	return
}

func (sd *StaffDashboard) DoLoginFromApp(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()
	start := sd.getTime(ctx)
	var result models.DoLoginFromAppResp

	defer func() {
		writeResponseJSON(w, result, reqErr)
		monitorAPI("DoLoginFromApp", start, reqErr)
	}()

	request := DoLoginReq{}
	err := getBody(r, &request)
	if err != nil {
		return
	}

	result, err = sd.authorizationUC.DoLoginFromApp(ctx, models.DoLoginReq{
		Email:      request.Email,
		Password:   request.Password,
		RememberMe: request.RememberMe,
		UserType:   sd.userType,
	})
	if err != nil {
		return
	}

	reqErr = result.GetError()
	return
}

func (sd *StaffDashboard) DoRefreshToken(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()
	start := sd.getTime(ctx)
	var result models.DoLoginFromAppResp

	defer func() {
		writeResponseJSON(w, result, reqErr)
		monitorAPI("DoRefreshToken", start, reqErr)
	}()

	request := DoRefreshTokenReq{}
	err := getBody(r, &request)
	if err != nil {
		return
	}

	result, err = sd.authorizationUC.DoRefreshAccessToken(ctx, models.DoRefreshAccessTokenReq{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		return
	}

	reqErr = result.GetError()
	return
}
