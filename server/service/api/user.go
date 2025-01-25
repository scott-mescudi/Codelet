package api

import (
	"net/http"
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

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Invalid content type, expected application/json"}`))
		return
	}

	var info = LoginPool.Get().(*UserLogin)
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

	userID, passwordHash, err := dba.GetUserPasswordHash(s.Db, info.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(info.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken := auth.GenerateHMac(userID, ACCESS, time.Now().Add(1 * time.Hour))
	refreshToken := auth.GenerateHMac(userID, REFRESH, time.Now().Add(48 * time.Hour))

	http.SetCookie(w, &http.Cookie{
		Name:     "CODELET-JWT-REFRESH-TOKEN",
		Value:    refreshToken,
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:  true,

	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"acess_token":accessToken}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) Refresh(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("CODELET-JWT-REFRESH-TOKEN")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userID, tokenType, err := auth.ValidateHmac(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if tokenType != REFRESH{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userID == -1 {
		w.WriteHeader(http.StatusBadRequest)
		return		
	}

	accessToken := auth.GenerateHMac(userID, ACCESS, time.Now().Add(1 * time.Hour))
	refreshToken := auth.GenerateHMac(userID, REFRESH, time.Now().Add(48 * time.Hour))

	http.SetCookie(w, &http.Cookie{
		Name:     "CODELET-JWT-REFRESH-TOKEN",
		Value:    refreshToken,
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:  true,

	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"acess_token":accessToken}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (s *Server) Logout(w http.ResponseWriter, r *http.Request){
	http.SetCookie(w, &http.Cookie{
		Name:     "CODELET-JWT-REFRESH-TOKEN",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:  true,
	})

	w.WriteHeader(http.StatusOK)
}

func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {

}
