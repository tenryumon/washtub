package staff_dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/nyelonong/boilerplate-go/core/helper"
	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/core/monitor"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

const (
	contextUserID    = "X-User-ID"
	contextOrgID     = "X-Organization-ID"
	contextStaffRole = "X-Staff-Role"
)

// Middleware to check whether User is login or not
//
// -------------------------
func (sd *StaffDashboard) UseMustLogin(h http.Handler) http.Handler {
	return http.HandlerFunc(sd.mustLogin(h))
}

func (sd *StaffDashboard) MustLogin(h http.HandlerFunc) http.HandlerFunc {
	return sd.mustLogin(http.HandlerFunc(h))
}

func (sd *StaffDashboard) mustLogin(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var loginDevice int

		token := sd.getAuthorizationToken(r)
		if token != "" {
			// Use JWT Auth
			loginDevice = entities.LoginDeviceMobileApp
		} else {
			// Use old auth
			token = sd.getSessionToken(r)
			loginDevice = entities.LoginDeviceDesktop
		}

		loginID := sd.checkLogin(ctx, token, loginDevice)
		if loginID == 0 {
			reqErr := &models.Errors{}
			reqErr.ErrorNotLogin()
			writeResponseJSON(w, nil, reqErr.GetError())
			return
		}

		extend, rememberMe := sd.needExtend(r)
		if extend {
			session := sd.authorizationUC.ExtendLogin(ctx, token, rememberMe, loginID)
			if session.IsSuccess() {
				http.SetCookie(w, &http.Cookie{
					Name:    sd.cookieName,
					Value:   sd.setSessionToken(token, rememberMe),
					Path:    "/",
					Expires: session.ExpireAt,
				})
				http.SetCookie(w, &http.Cookie{
					Name:    sd.autoExtendCookie,
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

func (sd *StaffDashboard) MustNotLogin(h http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var loginDevice int

		token := sd.getAuthorizationToken(r)
		if token != "" {
			loginDevice = entities.LoginDeviceMobileApp
		} else {
			token = sd.getSessionToken(r)
			loginDevice = entities.LoginDeviceDesktop
		}

		loginID := sd.checkLogin(ctx, token, loginDevice)
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

func (sd *StaffDashboard) checkLogin(ctx context.Context, token string, loginDevice int) (userID int64) {
	if token == "" {
		return
	}

	req := models.CheckLoginReq{
		Token:       token,
		LoginDevice: loginDevice,
	}

	if result, err := sd.authorizationUC.CheckLogin(ctx, req); err == nil {
		if result.UserID != 0 {
			userID = result.UserID
		}
	}

	return
}

func (sd *StaffDashboard) getSessionToken(r *http.Request) string {
	sessionToken, err := r.Cookie(sd.cookieName)
	if err != nil {
		return ""
	}

	s := strings.Split(sessionToken.Value, ":")
	return s[0]
}

func (sd *StaffDashboard) setSessionToken(token string, rememberMe bool) string {
	return fmt.Sprintf("%s:%t", token, rememberMe)
}

func (sd *StaffDashboard) hasExtExpiry(r *http.Request) bool {
	extExpiry, err := r.Cookie(sd.autoExtendCookie)
	if err != nil {
		return false
	}
	return extExpiry.Value == "1"
}

func (sd *StaffDashboard) needExtend(r *http.Request) (bool, bool) {
	extend := false
	rememberMe := false
	session, err := r.Cookie(sd.cookieName)
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

	if !sd.hasExtExpiry(r) {
		extend = true
	}

	return extend, rememberMe
}

func (sd *StaffDashboard) getUserID(ctx context.Context) int64 {
	userID, _ := ctx.Value(contextUserID).(int64)
	return userID
}

// Middleware to check whether has access or not
//
// -------------------------
func (sd *StaffDashboard) hasAccess(h http.HandlerFunc, roles []string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqErr := &models.Errors{}
		reqErr.InitError()

		staffRole := sd.getStaffRole(ctx)
		if !includeRole(staffRole, roles) {
			reqErr.ErrorNoAccess()
			writeResponseJSON(w, nil, reqErr.GetError())
			return
		}

		h.ServeHTTP(w, r)
	}

	return fn
}

func (sd *StaffDashboard) IsOwner(h http.HandlerFunc) http.HandlerFunc {
	roles := []string{entities.RoleOwner}
	return sd.hasAccess(h, roles)
}

func (sd *StaffDashboard) IsAdmin(h http.HandlerFunc) http.HandlerFunc {
	roles := []string{entities.RoleOwner, entities.RoleOrgAdmin}
	return sd.hasAccess(h, roles)
}

func (sd *StaffDashboard) HasAccess(h http.HandlerFunc, roles []string) http.HandlerFunc {
	newRoles := []string{entities.RoleOwner}
	newRoles = append(newRoles, roles...)
	return sd.hasAccess(h, newRoles)
}

func includeRole(staffRole []string, roles []string) bool {
	for _, sr := range staffRole {
		for _, r := range roles {
			if sr == r {
				return true
			}
		}
	}

	return false
}

func (sd *StaffDashboard) getOrganizationID(ctx context.Context) int64 {
	orgID, _ := ctx.Value(contextOrgID).(int64)
	return orgID
}

func (sd *StaffDashboard) getStaffRole(ctx context.Context) []string {
	roleString, _ := ctx.Value(contextStaffRole).(string)
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

func getQuery(r *http.Request, request interface{}) (err error) {
	if err = helper.QueryParser(r, request); err != nil {
		log.Errorf("Failed Parse Query: %s", err)
		return
	}
	return
}

// General Response Payload for Staff Dashboard
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
	w.Header().Add("Content-type", "application/json")

	w.WriteHeader(status)
	w.Write(result)
}

func monitorAPI(apiName string, start time.Time, err models.Errors) {
	monitor.Record("api_call", start, map[string]string{"api": apiName, "status": err.StatusCode, "usecase": "staff"})
}

func (sd *StaffDashboard) getTime(ctx context.Context) time.Time {
	return time.Now()
}

func (sd *StaffDashboard) getAuthorizationToken(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return ""
	}

	token := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	return token
}
