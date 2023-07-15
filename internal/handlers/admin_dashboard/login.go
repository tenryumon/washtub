package admin_dashboard

import (
	"net/http"

	"github.com/nyelonong/boilerplate-go/internal/models"
)

type DoLoginReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

func (ad *AdminDashboard) DoLogin(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()
	start := ad.getTime(ctx)
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

	result, err = ad.authorizationUC.DoLogin(ctx, models.DoLoginReq{
		Email:      request.Email,
		Password:   request.Password,
		RememberMe: request.RememberMe,
		UserType:   models.UserTypeAdmin,
	})
	if err != nil {
		return
	}

	reqErr = result.GetError()
	if reqErr.IsSuccess() {
		http.SetCookie(w, &http.Cookie{
			Name:    ad.cookieName,
			Value:   ad.setSessionToken(result.Token, request.RememberMe),
			Path:    "/",
			Expires: result.ExpireAt,
		})
	}

	return
}

func (ad *AdminDashboard) DoLogout(w http.ResponseWriter, r *http.Request) {
	var reqErr models.Errors
	ctx := r.Context()
	start := ad.getTime(ctx)

	defer func() {
		writeResponseJSON(w, nil, reqErr)
		monitorAPI("DoLogout", start, reqErr)
	}()

	result := ad.authorizationUC.DoLogout(ctx, ad.getSessionToken(r))
	reqErr = result.GetError()

	http.SetCookie(w, &http.Cookie{
		Name:    ad.cookieName,
		Value:   "",
		Path:    "/",
		Expires: result.ExpireAt,
		MaxAge:  -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    ad.autoExtendCookie,
		Value:   "",
		Path:    "/",
		Expires: result.ExpireAt,
		MaxAge:  -1,
	})

	return
}
