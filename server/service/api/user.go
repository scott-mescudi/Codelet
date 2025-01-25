package api

import (
	"fmt"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
	dba "github.com/scott-mescudi/codelet/service/data_access"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var SignupPool = &sync.Pool{
	New: func() any {
		return &UserSignup{}
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

	err := dba.AddUser(s.Db, info.Username, info.Email, "user", info.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(err)

}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

}


func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {

}
