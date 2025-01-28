package users

import (
	"fmt"
	"log"
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
		errs.ErrorWithJson(w, http.StatusBadRequest, "Content-Type header must be application/json")
		return
	}

	var info = SignupPool.Get().(*UserSignup)
	defer SignupPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Invalid JSON payload: "+err.Error())
		return
	}

	if !VerifyEmail(info.Email) {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Email field is invalid")
		return
	}

	if info.Password == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Password field is required")
		return
	}

	if info.Username == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Username field is required")
		return
	}

	if info.Role == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Role field is required")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Internal server error while processing password")
		return
	}

	err = dba.AddUser(s.Db, info.Username, info.Email, info.Role, string(hashedPassword))
	if err != nil {
		errs.ErrorWithJson(w, http.StatusBadRequest, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "application/json" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Content-Type header must be application/json")
		return
	}

	var info = LoginPool.Get().(*UserLogin)
	defer LoginPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Invalid JSON payload: "+err.Error())
		return
	}

	if info.Email == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Email field is invalid")
		return
	}

	if info.Password == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Password field is required")
		return
	}

	userID, passwordHash, err := dba.GetUserPasswordHash(s.Db, info.Email)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(info.Password)); err != nil {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	accessToken := auth.GenerateHMac(userID, ACCESS, time.Now().Add(15*time.Minute))
	refreshToken := auth.GenerateHMac(userID, REFRESH, time.Now().Add(48*time.Hour))

	if err := dba.AddRefreshToken(s.Db, refreshToken, userID); err != nil {
		log.Printf("Failed to add refresh token: %v", err)
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
	if err := json.NewEncoder(w).Encode(map[string]string{"acess_token": accessToken}); err != nil {
		log.Printf("Failed to encode response: %v", err)
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Failed to generate response")
		return
	}
}

func (s *UserService) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("CODELET-JWT-REFRESH-TOKEN")
	if err != nil {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "No refresh token provided")
		return
	}

	if cookie.Value == "" {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	userID, tokenType, err := auth.ValidateHmac(cookie.Value)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	if tokenType != REFRESH {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid token type")
		return
	}

	if userID == -1 {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "Invalid user token")
		return
	}

	dbToken, err := dba.GetRefreshToken(s.Db, userID)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	if cookie.Value != dbToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	accessToken := auth.GenerateHMac(userID, ACCESS, time.Now().Add(15*time.Minute))
	refreshToken := auth.GenerateHMac(userID, REFRESH, time.Now().Add(48*time.Hour))

	if err := dba.AddRefreshToken(s.Db, refreshToken, userID); err != nil {
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		errs.ErrorWithJson(w, http.StatusBadRequest, "Content-Type header must be application/json")
		return
	}

	defer r.Body.Close()

	var info = UpdatePasswordPool.Get().(*ChangePassword)
	defer UpdatePasswordPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "Invalid JSON payload: "+err.Error())
		return
	}

	if info.OldPassword == "" || info.NewPassword == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Both old and new passwords are required")
		return
	}

	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "User ID not found in request")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	passwordHash, err := dba.GetUserPasswordHashViaID(s.Db, userID)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusUnauthorized, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(info.OldPassword)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(info.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Failed to process new password")
		return
	}

	err = dba.UpdatePassword(s.Db, string(hashedNewPassword), userID)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		errs.ErrorWithJson(w, http.StatusInternalServerError, "Failed to update password in database")
		return
	}

	w.WriteHeader(http.StatusOK)
}
