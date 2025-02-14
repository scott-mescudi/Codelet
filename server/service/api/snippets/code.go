package snippets

import (
	"net/http"
	"strconv"
	"strings"
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
		s.Logger.Warn().Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("Invalid Content-Type, expected application/json")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Content-Type must be 'application/json'")
		return
	}
	defer r.Body.Close()

	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		s.Logger.Warn().Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("Missing 'X-USERID' header")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'X-USERID' header")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		s.Logger.Warn().Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("invalid 'X-USERID' header format")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "invalid 'X-USERID' header format")
		return
	}

	var info = SnippetPool.Get().(*Snippet)
	defer SnippetPool.Put(info)
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		s.Logger.Warn().Int("userID", userID).Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("unable to parse request body")
		errs.ErrorWithJson(w, http.StatusUnprocessableEntity, "unable to parse request body")
		return
	}

	if info.Title == "" {
		s.Logger.Warn().Int("userID", userID).Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("Missing snippet title")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing title")
		return
	}

	if info.Language == "" {
		s.Logger.Warn().Int("userID", userID).Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("Missing snippet language")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing language")
		return
	}

	if info.Code == "" {
		s.Logger.Warn().Int("userID", userID).Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("Missing snippet code")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing code text")
		return
	}

	if len(info.Code) > 3072 {
		s.Logger.Warn().Int("userID", userID).Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("Data too large")
		errs.ErrorWithJson(w, http.StatusRequestEntityTooLarge, "code too large")
		return
	}

	if err := dba.AddSnippet(s.Db, userID, info.Language, info.Description, info.Title, info.Code, info.Private, info.Favorite, info.Tags, time.Now(), time.Now()); err != nil {
		s.Logger.Error().Int("userID", userID).Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg(err.Error())
		errs.ErrorWithJson(w, http.StatusConflict, "failed to add snippet to database")
		return
	}

	s.Logger.Info().Int("userID", userID).Str("function", "AddSnippet").Str("origin", r.RemoteAddr).Msg("Served user")
	w.WriteHeader(http.StatusCreated)
}

func (s *SnippetService) GetUserSnippets(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("Missing 'X-USERID' header")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'X-USERID' header")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("invalid 'X-USERID' header format")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "invalid 'X-USERID' header format")
		return
	}

	params := r.URL.Query()
	limitstr := params.Get("limit")
	pagestr := params.Get("page")

	var snippets []dba.DBsnippet
	if limitstr == "" || pagestr == "" {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("missing 'limit' or 'page' parametr")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Missing 'limit' or 'page' parameter")
		return
	}

	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("invalid 'limit' parameter")
		errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'limit' parameter")
		return
	}

	page, err := strconv.Atoi(pagestr)
	if err != nil {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("invalid 'page' parameter")
		errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'page' parameter")
		return
	}

	if limit <= 0 || page <= 0 {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("limit or page is smaller or equal to 0")
		errs.ErrorWithJson(w, http.StatusBadRequest, "'limit' and 'page' parameter must be greater than 0")
		return
	}

	if limit > 100 {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("max 'limit' is 100")
		errs.ErrorWithJson(w, http.StatusBadRequest, "max 'limit' is 100")
		return
	}

	offset := (page - 1) * limit
	snippets, err = dba.GetSnippetsByUserID(s.Db, userID, limit, offset)
	if err != nil {
		s.Logger.Error().Int("userID", userID).Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("failed to fetch snippets")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to fetch snippets from database")
		return
	}

	if len(snippets) == 0 {
		s.Logger.Warn().Int("userID", userID).Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("no snippets found for user")
		errs.ErrorWithJson(w, http.StatusNotFound, "no snippets found for user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snippets); err != nil {
		s.Logger.Error().Int("userID", userID).Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("failed to encode snippets")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to encode snippets as JSON")
		return
	}
}

func (s *SnippetService) GetPublicSnippets(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := r.URL.Query()
	limitstr := params.Get("limit")
	pagestr := params.Get("page")

	if limitstr == "" || pagestr == "" {
		s.Logger.Warn().Str("function", "GetPublicSnippets").Str("origin", r.RemoteAddr).Msg("missing 'limit' or 'page' url parameter")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'limit' or 'page' url parameter.")
		return
	}

	var snippets []dba.DBsnippet

	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		s.Logger.Warn().Str("function", "GetPublicSnippets").Str("origin", r.RemoteAddr).Msg("invalid 'limit' parameter")
		errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'limit' parameter")
		return
	}

	page, err := strconv.Atoi(pagestr)
	if err != nil {
		s.Logger.Warn().Str("function", "GetPublicSnippets").Str("origin", r.RemoteAddr).Msg("invalid 'page' parameter")
		errs.ErrorWithJson(w, http.StatusBadRequest, "invalid 'page' parameter")
		return
	}

	if limit <= 0 || page <= 0 {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("limit or page is smaller or equal to 0")
		errs.ErrorWithJson(w, http.StatusBadRequest, "'limit' and 'page' parameter must be greater than 0")
		return
	}

	if limit > 100 {
		s.Logger.Warn().Str("function", "GetPublicSnippets").Str("origin", r.RemoteAddr).Msg("max 'limit' is 100")
		errs.ErrorWithJson(w, http.StatusBadRequest, "max 'limit' is 100")
		return
	}

	offset := (page - 1) * limit
	snippets, err = dba.GetPublicSnippets(s.Db, limit, offset)
	if err != nil {
		s.Logger.Error().Str("function", "GetPublicSnippets").Str("origin", r.RemoteAddr).Msg("failed to fetch public snippets")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to fetch snippets from database")
		return
	}

	if len(snippets) == 0 {
		s.Logger.Warn().Str("function", "GetPublicSnippets").Str("origin", r.RemoteAddr).Msg("no public snippets found")
		errs.ErrorWithJson(w, http.StatusNotFound, "no snippets found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snippets); err != nil {
		s.Logger.Error().Str("function", "GetPublicSnippets").Str("origin", r.RemoteAddr).Msg("failed to encode snippets")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to encode snippets as JSON")
		return
	}
}

func (s *SnippetService) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	path := r.URL.Path
	idString := strings.TrimPrefix(path, "/api/v1/user/snippets/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		s.Logger.Warn().Str("function", "DeleteSnippet").Str("origin", r.RemoteAddr).Msg("failed to parse snippet id in uri")
		errs.ErrorWithJson(w, http.StatusBadRequest, "failed to parse snippet id in uri")
		return
	}

	if id <= 0 {
		s.Logger.Warn().Str("function", "DeleteSnippet").Str("origin", r.RemoteAddr).Msg("Snippet id must be a positive integer")
		errs.ErrorWithJson(w, http.StatusBadRequest, "Snippet id must be a positive integer")
		return
	}

	if err := dba.DeleteSnippet(s.Db, id); err != nil {
		s.Logger.Error().Int("snippetID", id).Str("function", "DeleteSnippet").Str("origin", r.RemoteAddr).Msg("failed to delete snippet")
		errs.ErrorWithJson(w, http.StatusBadRequest, "failed to delete snippet")
		return
	}

	s.Logger.Info().Int("snippetID", id).Str("function", "DeleteSnippet").Str("origin", r.RemoteAddr).Msg("Successfully deleted snippet")
}

func (s *SnippetService) GetSmallUserSnippets(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("Missing 'X-USERID' header")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'X-USERID' header")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("invalid 'X-USERID' header format")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "invalid 'X-USERID' header format")
		return
	}

	snippets, err := dba.GetSmallUserSnippets(s.Db, userID)
	if err != nil {
		s.Logger.Error().Int("userID", userID).Str("function", "GetSmallSnippets").Str("origin", r.RemoteAddr).Err(err).Msg("failed to fetch user snippets")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to fetch snippets from database")
		return
	}

	if len(snippets) == 0 {
		s.Logger.Warn().Int("userID", userID).Str("function", "GetSmallSnippets").Str("origin", r.RemoteAddr).Msg("no snippets found for user")
		errs.ErrorWithJson(w, http.StatusNotFound, "no snippets found for user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snippets); err != nil {
		s.Logger.Error().Int("userID", userID).Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("failed to encode snippets")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to encode snippets as JSON")
		return
	}
}

func (s *SnippetService) GetUserSnippetByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	useridStr := r.Header.Get("X-USERID")
	if useridStr == "" {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("Missing 'X-USERID' header")
		errs.ErrorWithJson(w, http.StatusBadRequest, "missing 'X-USERID' header")
		return
	}
	path := r.URL.Path
	idString := strings.TrimPrefix(path, "/api/v1/user/snippets/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		s.Logger.Warn().Str("function", "DeleteSnippet").Str("origin", r.RemoteAddr).Msg("failed to parse snippet id in uri")
		errs.ErrorWithJson(w, http.StatusBadRequest, "failed to parse snippet id in uri")
		return
	}

	userID, err := strconv.Atoi(useridStr)
	if err != nil {
		s.Logger.Warn().Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("invalid 'X-USERID' header format")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "invalid 'X-USERID' header format")
		return
	}

	snippet, err := dba.GetSnippetByIDAndUserID(s.Db, userID, id)
	if err != nil {
		s.Logger.Error().Int("userID", userID).Str("function", "GetSmallSnippets").Str("origin", r.RemoteAddr).Err(err).Msg("failed to fetch user snippets")
		errs.ErrorWithJson(w, http.StatusNotFound, "failed to fetch snippets from database")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snippet); err != nil {
		s.Logger.Error().Int("userID", userID).Str("function", "GetUserSnippets").Str("origin", r.RemoteAddr).Msg("failed to encode snippets")
		errs.ErrorWithJson(w, http.StatusInternalServerError, "failed to encode snippets as JSON")
		return
	}
}
