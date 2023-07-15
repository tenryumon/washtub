package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

var (
	querySelectPasswordRequest = "SELECT id, user_id, token, request_type, status, expired_time, created_by, created_time, updated_by, updated_time"
	queryGetPassRequestByToken = querySelectPasswordRequest + " FROM password_requests WHERE token = :token ORDER BY id desc LIMIT 1"
)

func (us *UserObject) CreateNewPasswordRequest(ctx context.Context, param entities.PasswordRequest) (entities.PasswordRequest, error) {
	result, err := us.db.Exec(ctx, queryInsertPasswordRequest, param)
	if err != nil {
		return param, err
	}

	param.ID, err = result.LastInsertId()
	return param, err
}

func (us *UserObject) GetPasswordRequestByToken(ctx context.Context, token string) (entities.PasswordRequest, error) {
	result := entities.PasswordRequest{}
	param := map[string]interface{}{
		"token": token,
	}

	err := us.db.Get(ctx, &result, queryGetPassRequestByToken, param)
	return result, err
}

func (us *UserObject) IsValidPasswordToken(ctx context.Context, token string, email string) (result entities.ValidPasswordRequest, err error) {
	passReq, err := us.GetPasswordRequestByToken(ctx, token)

	if err != nil {
		return result, fmt.Errorf("Failed to GetPasswordRequestByToken because %s", err)
	}
	if !passReq.Exist() {
		result.InvalidType = entities.ValidPassRequestNotFound
		return result, nil
	}
	if !passReq.IsActive() {
		result.InvalidType = entities.ValidPassRequestInactive
		return result, nil
	}
	if passReq.IsExpired() {
		result.InvalidType = entities.ValidPassRequestExpired
		return result, nil
	}

	user, err := us.GetUserByID(ctx, passReq.UserID)
	if err != nil {
		return result, fmt.Errorf("Failed to GetUserByID because %s", err)
	}
	if user.Email != email {
		result.InvalidType = entities.ValidPassRequestNotFound
		return result, nil
	}

	result.Valid = true
	result.PassReq = passReq
	return result, nil
}

func (us *UserObject) hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to GenerateFromPassword because %s", err)
	}

	return string(hashed), nil
}

func (us *UserObject) SetPassword(ctx context.Context, token string, email string, password string) (bool, bool, entities.ValidPasswordRequest, error) {
	result, err := us.IsValidPasswordToken(ctx, token, email)
	if err != nil {
		return false, false, result, fmt.Errorf("Failed to SetPassword because %s", err)
	}
	if !result.Valid {
		return false, false, result, nil
	}

	// Check whether logins row alread exist or not, either insert or update it.
	login, err := us.getLogin(ctx, result.PassReq.UserID)
	if err != nil {
		return false, false, result, fmt.Errorf("Failed to getLogin because %s", err)
	}

	// Hash Password before save it to database
	hashedPassword, err := us.hashPassword(password)
	if err != nil {
		return false, false, result, fmt.Errorf("Failed to hashPassword because %s", err)
	}

	// Updated By the user itself
	passReq := result.PassReq
	passReq.UpdatedBy = passReq.UserID

	tx, err := us.db.Begin()
	if err != nil {
		return false, false, result, fmt.Errorf("Failed to create db Tx because %s", err)
	}
	defer tx.Rollback()

	err = us.finishPasswordRequest(ctx, tx, passReq)
	if err != nil {
		return false, false, result, fmt.Errorf("Failed to finishPasswordRequest because %s", err)
	}

	var isNewUser bool
	if login.Exist() {
		login.Token = hashedPassword
		login.UpdatedBy = passReq.UserID
		err = us.updateLogin(ctx, tx, login)
		isNewUser = false
	} else {
		login = entities.Login{
			UserID:    passReq.UserID,
			LoginType: entities.LoginTypePassword,
			Token:     string(hashedPassword),
			CreatedBy: passReq.UserID,
			UpdatedBy: passReq.UserID,
		}
		err = us.insertNewLogin(ctx, tx, login)
		isNewUser = true
	}
	if err != nil {
		return false, false, result, fmt.Errorf("Failed to update/insert Login because %s", err)
	}

	err = us.activateUser(ctx, tx, passReq.UserID)
	if err != nil {
		return false, false, result, fmt.Errorf("Failed to activateUser because %s", err)
	}

	tx.Commit()

	return true, isNewUser, result, nil
}

func (us *UserObject) UpdatePassword(ctx context.Context, login entities.Login, newPassword string) error {
	// Hash Password before save it to database
	hashedPassword, err := us.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("Failed to hashPassword because %s", err)
	}
	login.Token = hashedPassword

	tx, err := us.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to create db Tx because %s", err)
	}
	defer tx.Rollback()

	err = us.updateLogin(ctx, tx, login)
	if err != nil {
		return fmt.Errorf("Failed to updateLogin because %s", err)
	}

	tx.Commit()
	return nil
}

var (
	queryFinishPasswordRequest = `UPDATE password_requests SET status = 2, updated_by = :updated_by WHERE id = :id`
)

func (us *UserObject) finishPasswordRequest(ctx context.Context, tx *database.Tx, param entities.PasswordRequest) error {
	_, err := tx.Exec(ctx, queryFinishPasswordRequest, param)
	return err
}

var (
	queryInsertNewLogin = `INSERT INTO logins (user_id, login_type, token, status, created_by, updated_by)
							VALUES (:user_id, :login_type, :token, 1, :created_by, :updated_by)`
	queryUpdateLogin = `UPDATE logins SET token = :token, updated_by = :updated_by
						WHERE user_id = :user_id AND login_type = :login_type`
)

func (us *UserObject) insertNewLogin(ctx context.Context, tx *database.Tx, param entities.Login) error {
	_, err := tx.Exec(ctx, queryInsertNewLogin, param)
	return err
}

func (us *UserObject) updateLogin(ctx context.Context, tx *database.Tx, param entities.Login) error {
	_, err := tx.Exec(ctx, queryUpdateLogin, param)
	return err
}

var (
	queryActivateUser = `UPDATE users SET status = 1, updated_by = :user_id WHERE id = :user_id`
)

func (us *UserObject) activateUser(ctx context.Context, tx *database.Tx, id int64) error {
	param := map[string]interface{}{
		"user_id": id,
	}

	_, err := tx.Exec(ctx, queryActivateUser, param)
	return err
}
