package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var anagrams = map[string][]string{
	"act":   []string{"cat", "act"},
	"ehllo": []string{"hello"},
	"dlorw": []string{"world"},
	"acot":  []string{"coat", "taco"},
	"abcd":  []string{"abcd", "cadb"},
}

type erroringServer struct{}

func (s *erroringServer) ListenAndServe() error {
	return errors.New("failed")
}

func TestMain(t *testing.T) {
	create = func(addr string, handler http.Handler) listener {
		return &erroringServer{}
	}

	main() // will immediately error and exit
}

func TestCreateServer(t *testing.T) {
	create = func(addr string, handler http.Handler) listener {
		return &http.Server{Addr: addr, Handler: handler}
	}

	s, _ := createServer("anything")
	server, ok := s.(*http.Server)
	if !ok {
		t.Fatal("Expected http.Server", server, ok)
	}
	router, ok := server.Handler.(*httprouter.Router)
	if !ok {
		t.Fatal("Expected httprouter.Router")
	}
	handle, params, _ := router.Lookup("GET", "/anagrams/heLOL")
	w := httptest.NewRecorder()
	handle(w, nil, params)
	if w.Body.String() != `["hello"]` {
		t.Error("expected correct handler to call anagramFinder")
	}
}

func TestAnagramFinder(t *testing.T) {
	w := httptest.NewRecorder()
	p := httprouter.Params{httprouter.Param{Key: "letters", Value: "heLOL"}}
	anagramFinder(anagrams)(w, nil, p)
	if w.Body.String() != `["hello"]` {
		t.Error("Expected matching value")
	}
}

func TestHttpOut(t *testing.T) {
	w := httptest.NewRecorder()
	httpOut(w, "Has message", nil)
	if body := w.Body.String(); body != "Has message" {
		t.Error("Expected message", body, w.Code)
	}
	w = httptest.NewRecorder()
	httpOut(w, "Should error", errors.New("Has error"))
	if body := w.Body.String(); body != "Has error\n" || w.Code != http.StatusInternalServerError {
		t.Error("Expected error", body, w.Code)
	}
}

func TestFindWords(t *testing.T) {
	json, err := findWords(anagrams, "notaword")
	if json != "[]" || err != nil {
		t.Error("Expected not found")
	}
	json, err = findWords(anagrams, "LOLhe")
	if json != `["hello"]` || err != nil {
		t.Error("Expected found")
	}
}

func TestCreateDictionary(t *testing.T) {
	words := "hello,cat,act,world,coat,abcd,taco,cadb"
	expected := anagrams
	actual := createDictionary(words)
	if len(expected) != len(actual) {
		t.Fatal("Expected and actual length differ")
	}
	for key, value := range expected {
		checkArray(t, value, actual[key])
	}
}

func TestSortAlpha(t *testing.T) {
	if "ehllo" != sortAlpha("hello") {
		t.Error("Sort alpha is broken")
	}
}

func checkArray(t *testing.T, expected, actual []string) {
	if len(expected) != len(actual) {
		t.Error("Length doesn't match", len(expected), len(actual), "for", expected, actual)
		return
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Error("Values don't match at index", i, expected, actual)
		}
	}
}
