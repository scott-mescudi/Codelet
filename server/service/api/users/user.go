package users

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	dba "github.com/scott-mescudi/codelet/service/data_access"
	auth "github.com/scott-mescudi/codelet/shared/auth"
	errs "github.com/scott-mescudi/codelet/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const ACCESS = 0
const REFRESH = 1

var SignupPool = &sync.Pool{
	New: func() any {
		return &UserSignup{}
	},
}

var LoginPool = &sync.Pool{
	New: func() any {
		return &UserLogin{}
	},
}

var UpdatePasswordPool = &sync.Pool{
	New: func() any {
		return &ChangePassword{}
	},
}

func (s *UserService) Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "application/json" {
		s.Logger.Warn().Str("function", "Signup").Str("origin", r.RemoteAddr).Msg("Invalid Content-Type, expected application/json")
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Content-Type header must be application/json")
		return
	}

	var info = SignupPool.Get().(*UserSignup)
	defer SignupPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		s.Logger.Warn().Str("function", "Signup").Str("origin", r.RemoteAddr).Msg("Failed to decode body into json")
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Invalid JSON payload: "+err.Error())
		return
	}

	if !VerifyEmail(info.Email) {
		s.Logger.Warn().Str("function", "Signup").Str("origin", r.RemoteAddr).Msg("Invalid email")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Email field is invalid")
		return
	}

	if info.Password == "" {
		s.Logger.Warn().Str("function", "Signup").Str("origin", r.RemoteAddr).Msg("Invalid password")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Password field is required")
		return
	}

	if info.Username == "" {
		s.Logger.Warn().Str("function", "Signup").Str("origin", r.RemoteAddr).Msg("Invalid Username")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Username field is required")
		return
	}

	if info.Role == "" {
		s.Logger.Warn().Str("function", "Signup").Str("origin", r.RemoteAddr).Msg("Invalid Role")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Role field is required")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Warn().Str("function", "Signup").Str("origin", r.RemoteAddr).Msg("Failed to hash password")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Internal server error while processing password")
		return
	}

	err = dba.AddUser(s.Db, info.Username, info.Email, info.Role, string(hashedPassword))
	if err != nil {
		s.Logger.Error().Str("function", "Signup").Str("origin", r.RemoteAddr).Err(err)
		errs.ErrorWithJson(w, http.StatusBadRequest, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	s.Logger.Info().Str("function", "Signup").Str("origin", r.RemoteAddr).Str("user", info.Username).Msg("Created new user")
	w.WriteHeader(http.StatusCreated)
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "application/json" {
		s.Logger.Warn().Str("function", "Login").Str("origin", r.RemoteAddr).Msg("Invalid Content-Type, expected application/json")
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Content-Type header must be application/json")
		return
	}

	var info = LoginPool.Get().(*UserLogin)
	defer LoginPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		s.Logger.Warn().Str("function", "Login").Str("origin", r.RemoteAddr).Msg("Failed to decode body into json")
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Invalid JSON payload: "+err.Error())
		return
	}

	if info.Email == "" {
		s.Logger.Warn().Str("function", "Login").Str("origin", r.RemoteAddr).Msg("Invalid email")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Email field is invalid")
		return
	}

	if info.Password == "" {
		s.Logger.Warn().Str("function", "Login").Str("origin", r.RemoteAddr).Msg("Invalid password")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Password field is required")
		return
	}

	userID, passwordHash, last_login, err := dba.GetUserPasswordHashAndLastLogin(s.Db, info.Email)
	if err != nil {
		s.Logger.Warn().Str("function", "Login").Str("origin", r.RemoteAddr).Err(err).Msg("Failed to retrieve user password hash")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if last_login != nil && time.Since(*last_login) < 30*time.Second {
		s.Logger.Warn().Str("function", "Login").Str("origin", r.RemoteAddr).Int("Userid", userID).Msg("Login attempt blocked: user must wait before trying again")
		errs.ErrorWithJson(w, http.StatusTooManyRequests, "Too many login attempts. Please wait a 30 seconds and try again.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(info.Password)); err != nil {
		s.Logger.Warn().Str("function", "Login").Str("origin", r.RemoteAddr).Msg("Invalid password comparison")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	accessToken := auth.GenerateHMac(userID, ACCESS, time.Now().Add(120*time.Minute))
	refreshToken := auth.GenerateHMac(userID, REFRESH, time.Now().Add(48*time.Hour))

	if err := dba.UpdateTokenAndLoginTime(s.Db, refreshToken, time.Now(), userID); err != nil {
		s.Logger.Error().Str("function", "Login").Str("origin", r.RemoteAddr).Err(err).Msg("Failed to add refresh token")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Failed to complete login process")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "CODELET-JWT-REFRESH-TOKEN",
		Value:    refreshToken,
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken}); err != nil {
		s.Logger.Error().Str("function", "Login").Err(err).Msg("Failed to encode response")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Failed to generate response")
		return
	}

	s.Logger.Info().Str("function", "Login").Str("origin", r.RemoteAddr).Msg("User logged in successfully")
}

func (s *UserService) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("CODELET-JWT-REFRESH-TOKEN")
	if err != nil {
		s.Logger.Warn().Str("function", "Refresh").Str("origin", r.RemoteAddr).Err(err).Msg("No refresh token provided")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "No refresh token provided")
		return
	}

	if cookie.Value == "" {
		s.Logger.Warn().Str("function", "Refresh").Str("origin", r.RemoteAddr).Msg("Invalid refresh token")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	userID, tokenType, err := auth.ValidateHmac(cookie.Value)
	if err != nil {
		s.Logger.Warn().Str("function", "Refresh").Str("origin", r.RemoteAddr).Err(err).Msg("Invalid or expired refresh token")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	if tokenType != REFRESH {
		s.Logger.Warn().Str("function", "Refresh").Str("origin", r.RemoteAddr).Msg("Invalid token type")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid token type")
		return
	}

	if userID == -1 {
		s.Logger.Warn().Str("function", "Refresh").Str("origin", r.RemoteAddr).Msg("Invalid user token")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid user token")
		return
	}

	dbToken, err := dba.GetRefreshToken(s.Db, userID)
	if err != nil {
		s.Logger.Error().Str("function", "Refresh").Str("origin", r.RemoteAddr).Err(err).Msg("Failed to retrieve refresh token from database")
		errs.ErrorWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	if cookie.Value != dbToken {
		s.Logger.Warn().Str("function", "Refresh").Str("origin", r.RemoteAddr).Msg("Refresh token mismatch")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	accessToken := auth.GenerateHMac(userID, ACCESS, time.Now().Add(120*time.Minute))
	refreshToken := auth.GenerateHMac(userID, REFRESH, time.Now().Add(48*time.Hour))

	if err := dba.AddRefreshToken(s.Db, refreshToken, userID); err != nil {
		s.Logger.Error().Str("function", "Refresh").Str("origin", r.RemoteAddr).Err(err).Msg("Failed to add new refresh token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "CODELET-JWT-REFRESH-TOKEN",
		Value:    refreshToken,
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"acess_token": accessToken}); err != nil {
		s.Logger.Error().Str("function", "Refresh").Err(err).Msg("Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.Logger.Info().Str("function", "Refresh").Str("origin", r.RemoteAddr).Msg("Refresh token successfully renewed")
	w.WriteHeader(http.StatusOK)
}

func (s *UserService) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "CODELET-JWT-REFRESH-TOKEN",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	if err := dba.AddRefreshToken(s.Db, "", userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserService) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		s.Logger.Warn().Str("function", "ChangePassword").Str("origin", r.RemoteAddr).Msg("Invalid Content-Type, expected application/json")
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Content-Type header must be application/json")
		return
	}

	defer r.Body.Close()

	var info = UpdatePasswordPool.Get().(*ChangePassword)
	defer UpdatePasswordPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		s.Logger.Warn().Str("function", "ChangePassword").Str("origin", r.RemoteAddr).Err(err).Msg("Failed to decode body into json")
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Invalid JSON payload: "+err.Error())
		return
	}

	if info.OldPassword == "" || info.NewPassword == "" {
		s.Logger.Warn().Str("function", "ChangePassword").Str("origin", r.RemoteAddr).Msg("Both old and new passwords are required")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Both old and new passwords are required")
		return
	}

	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		s.Logger.Warn().Str("function", "ChangePassword").Str("origin", r.RemoteAddr).Msg("User ID not found in request")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "User ID not found in request")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		s.Logger.Error().Str("function", "ChangePassword").Err(err).Msg("Failed to parse user ID")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	passwordHash, err := dba.GetUserPasswordHashViaID(s.Db, userID)
	if err != nil {
		s.Logger.Warn().Str("function", "ChangePassword").Str("origin", r.RemoteAddr).Msg("User not found")
		errs.ErrorWithJson(w, http.StatusUnauthorized, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(info.OldPassword)); err != nil {
		s.Logger.Warn().Str("function", "ChangePassword").Str("origin", r.RemoteAddr).Msg("Old password does not match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(info.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error().Str("function", "ChangePassword").Err(err).Msg("Failed to hash new password")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Failed to process new password")
		return
	}

	err = dba.UpdatePassword(s.Db, string(hashedNewPassword), time.Now(), userID)
	if err != nil {
		s.Logger.Error().Str("function", "ChangePassword").Err(err).Msg("Failed to update password in database")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Failed to update password in database")
		return
	}

	s.Logger.Info().Str("function", "ChangePassword").Str("origin", r.RemoteAddr).Msg("Password changed successfully")
	w.WriteHeader(http.StatusOK)
}
