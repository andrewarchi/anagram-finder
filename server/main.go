package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var create = func(addr string, handler http.Handler) listener {
	return &http.Server{Addr: addr, Handler: handler}
}

type listener interface {
	ListenAndServe() error
}

func main() {
	s, err := createServer(":3141")
	if err != nil {
		log.Fatal(err)
	}
	s.ListenAndServe()
}

func createServer(addr string) (listener, error) {
	words, err := ioutil.ReadFile("words.txt") // https://github.com/dwyl/english-words/blob/master/words_alpha.txt
	if err != nil {
		return nil, err
	}
	anagrams := createDictionary(string(words))
	router := httprouter.New()
	router.GET("/anagrams/:letters", anagramFinder(anagrams))
	router.NotFound = http.FileServer(http.Dir("../client"))
	return create(addr, router), nil
}

func anagramFinder(anagrams map[string][]string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		json, err := findWords(anagrams, p.ByName("letters"))
		httpOut(w, json, err)
	}
}

func httpOut(w http.ResponseWriter, message string, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, message)
}

func findWords(anagrams map[string][]string, letters string) (string, error) {
	sorted := sortAlpha(strings.ToLower(letters))
	if value, ok := anagrams[sorted]; ok {
		data, err := json.Marshal(value)
		return string(data), err
	}
	return "[]", nil
}

func createDictionary(words string) map[string][]string {
	anagrams := make(map[string][]string)
	for _, word := range strings.Split(words, ",") {
		sorted := sortAlpha(word)
		if value, ok := anagrams[sorted]; ok {
			anagrams[sorted] = append(value, word)
		} else {
			anagrams[sorted] = []string{word}
		}
	}
	return anagrams
}

func sortAlpha(word string) string {
	s := strings.Split(word, "")
	sort.Strings(s)
	return strings.Join(s, "")
}
