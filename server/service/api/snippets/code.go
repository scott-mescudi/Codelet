package snippets

import (
	jsoniter "github.com/json-iterator/go"
	
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *SnippetService)AddSnippet(w http.ResponseWriter, r *http.Request) {

}