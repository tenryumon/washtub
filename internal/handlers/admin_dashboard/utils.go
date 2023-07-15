package admin_dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/core/monitor"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

const (
	contextUserID    = "X-User-ID"
	contextAdminRole = "X-Role"
)

// Middleware to check whether User is login or not
//
// -------------------------
func (ad *AdminDashboard) UseMustLogin(h http.Handler) http.Handler {
	return http.HandlerFunc(ad.mustLogin(h))
}

func (ad *AdminDashboard) MustLogin(h http.HandlerFunc) http.HandlerFunc {
	return ad.mustLogin(http.HandlerFunc(h))
}

func (ad *AdminDashboard) mustLogin(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := ad.getSessionToken(r)

		loginID := ad.checkLogin(ctx, token)
		if loginID == 0 {
			reqErr := &models.Errors{}
			reqErr.ErrorNotLogin()
			writeResponseJSON(w, nil, reqErr.GetError())
			return
		}

		extend, rememberMe := ad.needExtend(r)
		if extend {
			session := ad.authorizationUC.ExtendLogin(ctx, token, rememberMe, loginID)
			if session.IsSuccess() {
				http.SetCookie(w, &http.Cookie{
					Name:    ad.cookieName,
					Value:   ad.setSessionToken(token, rememberMe),
					Path:    "/",
					Expires: session.ExpireAt,
				})
				http.SetCookie(w, &http.Cookie{
					Name:    ad.autoExtendCookie,
					Value:   "1",
					Path:    "/",
					Expires: time.Now().Add(3 * time.Hour),
				})
			}
		}

		ctx = context.WithValue(ctx, contextUserID, loginID)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}
}

// Middleware to add additional context value
//
// -------------------------
func (ad *AdminDashboard) UseAddContextValue(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqErr := &models.Errors{}
		reqErr.InitError()

		userID := ad.getUserID(ctx)
		roles, err := ad.authorizationUC.GetRolesByUserID(ctx, userID)
		if err != nil {
			writeResponseJSON(w, nil, reqErr.GetError())
			return
		}

		ctx = context.WithValue(ctx, contextAdminRole, strings.Join(roles, ";"))

		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (ad *AdminDashboard) MustNotLogin(h http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		loginID := ad.checkLogin(ctx, ad.getSessionToken(r))
		if loginID != 0 {
			reqErr := &models.Errors{}
			reqErr.ErrorAlreadyLogin()
			writeResponseJSON(w, nil, reqErr.GetError())
			return
		}

		h.ServeHTTP(w, r)
	}

	return fn
}

func (ad *AdminDashboard) checkLogin(ctx context.Context, token string) (userID int64) {
	if token == "" {
		return userID
	}

	result, err := ad.authorizationUC.CheckLogin(ctx, models.CheckLoginReq{Token: token})
	if err == nil && result.UserID != 0 {
		userID = result.UserID
	}

	return userID
}

func (ad *AdminDashboard) getSessionToken(r *http.Request) string {
	sessionToken, err := r.Cookie(ad.cookieName)
	if err != nil {
		return ""
	}

	s := strings.Split(sessionToken.Value, ":")
	return s[0]
}

func (ad *AdminDashboard) setSessionToken(token string, rememberMe bool) string {
	return fmt.Sprintf("%s:%t", token, rememberMe)
}

func (ad *AdminDashboard) hasExtExpiry(r *http.Request) bool {
	extExpiry, err := r.Cookie(ad.autoExtendCookie)
	if err != nil {
		return false
	}
	return extExpiry.Value == "1"
}

func (ad *AdminDashboard) needExtend(r *http.Request) (bool, bool) {
	extend := false
	rememberMe := false
	session, err := r.Cookie(ad.cookieName)
	if err != nil {
		return extend, rememberMe
	}

	s := strings.Split(session.Value, ":")
	if len(s) > 1 && s[1] == "true" {
		rememberMe = true
	}

	expires := 24 * time.Hour
	if rememberMe {
		expires = 7 * expires
	}

	if !ad.hasExtExpiry(r) {
		extend = true
	}

	return extend, rememberMe
}

func (ad *AdminDashboard) getUserID(ctx context.Context) int64 {
	userID, _ := ctx.Value(contextUserID).(int64)
	return userID
}

// Middleware to check whether has access or not
//
// -------------------------
func (ad *AdminDashboard) hasAccess(h http.HandlerFunc, roles []string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqErr := &models.Errors{}
		reqErr.InitError()

		userLoginRole := ad.getUserLoginRole(ctx)
		if !includeRole(userLoginRole, roles) {
			reqErr.ErrorNoAccess()
			writeResponseJSON(w, nil, reqErr.GetError())
			return
		}

		h.ServeHTTP(w, r)
	}

	return fn
}

func (ad *AdminDashboard) IsAdmin(h http.HandlerFunc) http.HandlerFunc {
	roles := []string{entities.RoleSuperAdmin}
	return ad.hasAccess(h, roles)
}

func (ad *AdminDashboard) HasAccess(h http.HandlerFunc, roles []string) http.HandlerFunc {
	newRoles := []string{entities.RoleSuperAdmin}
	newRoles = append(newRoles, roles...)
	return ad.hasAccess(h, newRoles)
}

func includeRole(adminRole []string, roles []string) bool {
	for _, sr := range adminRole {
		for _, r := range roles {
			if sr == r {
				return true
			}
		}
	}

	return false
}

func (ad *AdminDashboard) getUserLoginRole(ctx context.Context) []string {
	roleString, _ := ctx.Value(contextAdminRole).(string)
	return strings.Split(roleString, ";")
}

// Simple Get Body JSON
//
// -------------------------
func getBody(r *http.Request, request interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed Read Body: %s", err)
		return err
	}
	log.Debugf("Request: %s", string(body))

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Errorf("Failed Umarshal Body: %s", err)
	}
	return err
}

// General Response Payload for Admin Dashboard
//
// -------------------------
type Response struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
}

func writeResponseJSON(w http.ResponseWriter, data interface{}, err models.Errors) {
	errors := map[string]string{}
	for key, value := range err.ValidationErrors {
		errors[key] = value[0]
	}

	response := Response{
		Code:    err.StatusCode,
		Message: err.ErrorMessage,
		Errors:  errors,
		Data:    data,
	}

	result, _ := json.Marshal(response)
	status := err.GetStatusCode()
	log.Debugf("Response: %s", string(result))

	w.WriteHeader(status)
	w.Write(result)
}

func monitorAPI(apiName string, start time.Time, err models.Errors) {
	monitor.Record("api_call", start, map[string]string{"api": apiName, "status": err.StatusCode, "usecase": "admin"})
}

func (ad *AdminDashboard) getTime(ctx context.Context) time.Time {
	return time.Now()
}
