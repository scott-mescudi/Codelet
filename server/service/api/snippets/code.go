package snippets

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	dba "github.com/scott-mescudi/codelet/service/data_access"
	errs "github.com/scott-mescudi/codelet/shared/errors"
)

var SnippetPool = &sync.Pool{
	New: func() any {
		return &Snippet{}
	},
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *SnippetService) AddSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "Content-Type must be 'application/json'")
		return
	}

	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'X-USERID' header")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, "invalid 'X-USERID' header format")
		return
	}

	var info = SnippetPool.Get().(*Snippet)
	defer SnippetPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		errs.ErrorWithJson(w, http.StatusBadRequest, "unable to parse request body")
		return
	}

	if info.Title == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing title")
		return
	}

	if info.Language == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing language")
		return
	}

	if info.Code == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing code text")
		return
	}

	if len(info.Code) > 3072 {
		errs.ErrorWithJson(w, http.StatusNotAcceptable, "code too large")
		return
	}

	if err := dba.AddSnippet(s.Db, userID, info.Language, info.Description, info.Title, info.Code, info.Private, info.Tags, time.Now(), time.Now()); err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to add snippet to database")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *SnippetService) GetUserSnippets(w http.ResponseWriter, r *http.Request) {
	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'X-USERID' header")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, "invalid 'X-USERID' header format")
		return
	}

	params := r.URL.Query()
	limitstr := params.Get("limit")
	pagestr := params.Get("page")
	

	var snippets []dba.DBsnippet
	if limitstr != "" && pagestr != "" {
		limit, err := strconv.Atoi(limitstr)
		if err != nil {
			errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'limit' parameter")
			return
		}

		page, err := strconv.Atoi(pagestr)
		if err != nil {
			errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'page' parameter")
			return
		}

		if limit > 100 {
			errs.ErrorWithJson(w, http.StatusBadRequest, "max 'limit' is 100")
			return
		}

		offset := (page - 1) * limit
		snippets, err = dba.GetSnippetsByUserID(s.Db, userID, limit, offset)
		if err != nil {
			errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to fetch snippets from database")
			return
		}
	} else {
		snippets, err = dba.GetAllSnippetsByUserID(s.Db, userID)
		if err != nil {
			errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to fetch snippets from database")
			return
		}
	}

	if len(snippets) == 0 {
		errs.ErrorWithJson(w, http.StatusNotFound, "no snippets found for user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snippets); err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to encode snippets as JSON")
		return
	}
}


func (s *SnippetService) GetPublicSnippets(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	limitstr := params.Get("limit")
	pagestr := params.Get("page")


	if  limitstr == "" || pagestr == "" {
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'limit' or 'page' url parameter.")
		return
	}

	var snippets []dba.DBsnippet

	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'limit' parameter")
		return
	}

	page, err := strconv.Atoi(pagestr)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'page' parameter")
		return
	}

	if limit > 100 {
		errs.ErrorWithJson(w, http.StatusBadRequest, "max 'limit' is 100")
		return
	}

	offset := (page - 1) * limit
	snippets, err = dba.GetPublicSnippets(s.Db, limit, offset)
	if err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to fetch snippets from database")
		return
	}
	
	if len(snippets) == 0 {
		errs.ErrorWithJson(w, http.StatusNotFound, "no snippets found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snippets); err != nil {
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to encode snippets as JSON")
		return
	}
}
