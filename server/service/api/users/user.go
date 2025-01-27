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
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Invalid content type, expected application/json"}`))
		return
	}

	var info = SignupPool.Get().(*UserSignup)
	defer SignupPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Role == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = dba.AddUser(s.Db, info.Username, info.Email, info.Role, string(hashedPassword))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Invalid content type, expected application/json"}`))
		return
	}

	var info = LoginPool.Get().(*UserLogin)
	defer LoginPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, passwordHash, err := dba.GetUserPasswordHash(s.Db, info.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(info.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
}

func (s *UserService) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("CODELET-JWT-REFRESH-TOKEN")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if cookie.Value == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userID, tokenType, err := auth.ValidateHmac(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if tokenType != REFRESH {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userID == -1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbToken, err := dba.GetRefreshToken(s.Db, userID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if cookie.Value != dbToken {
		fmt.Println("here")
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
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Invalid content type, expected application/json"}`))
		return
	}

	var info = UpdatePasswordPool.Get().(*ChangePassword)
	defer UpdatePasswordPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if info.OldPassword == "" || info.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	passwordHash, err := dba.GetUserPasswordHashViaID(s.Db, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(info.OldPassword)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(info.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dba.UpdatePassword(s.Db, string(hashedNewPassword), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
