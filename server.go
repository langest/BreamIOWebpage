package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	PORT   = 8080
	STATIC = "static"
	LANG   = "lang"
)

type Language struct {
	Abbr, Country string
	Quality       float64
}

type LanguageSlice []Language

type StaticHandler struct {
	Path string
}

type IndexHandler struct {
	Files map[string]*StaticHandler
}

func newStaticHandler(path string) *StaticHandler {
	return &StaticHandler{path}
}

func (sh StaticHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, sh.Path)
	log.Printf("Served request for %s", filepath.Base(sh.Path))
}

// Accept-Language:sv,en-US;q=0.8,en;q=0.6
func parseLanguage(header string) (codes LanguageSlice) {
	langs := strings.Split(header, ",")
	codes = make([]Language, len(langs))
	for i, lang := range langs {
		code := Language{
			Country: "",
			Quality: 1,
		}

		pair := strings.Split(lang, ";")
		if len(pair) == 2 {
			fmt.Sscanf(pair[1], "q=%f", &code.Quality)
		}

		pair = strings.Split(pair[0], "-")
		if len(pair) == 2 {
			code.Country = pair[1]
		}

		code.Abbr = pair[0]
		codes[i] = code
	}
	return
}

func (langs LanguageSlice) Len() int {
	return len(langs)
}

func (langs LanguageSlice) Swap(i, j int) {
	langs[i], langs[j] = langs[j], langs[i]
}

func (langs LanguageSlice) Less(i, j int) bool {
	return langs[i].Quality > langs[j].Quality
}

func newIndexHandler(cwd string) (ih *IndexHandler) {
	ih = &IndexHandler{
		Files: make(map[string]*StaticHandler),
	}
	base := filepath.Join(cwd, LANG)
	ls, err := ioutil.ReadDir(base)
	die("ReadDir:", err)
	for _, dir := range ls {
		if dir.IsDir() {
			path := filepath.Join(base, dir.Name(), "index.html")
			ih.Files[dir.Name()] = newStaticHandler(path)
		}
	}
	return
}

func (ih IndexHandler) MatchLanguage(header string) (*StaticHandler, Language) {
	langs := parseLanguage(header)
	sort.Stable(langs)
	for _, lang := range langs {
		if handler, ok := ih.Files[lang.Abbr]; ok {
			return handler, lang
		}
	}
	return ih.Files["en"], Language{}
}

func (ih IndexHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handler, lang := ih.MatchLanguage(req.Header.Get("Accept-Language"))
	handler.ServeHTTP(rw, req)
	log.Printf("Best language match: %v", lang)
}

func die(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func main() {
	cwd, err := os.Getwd()
	die("Getwd:", err)
	ls, err := ioutil.ReadDir(filepath.Join(cwd, STATIC))
	die("ReadDir:", err)

	http.Handle("/", newIndexHandler(cwd))
	for _, file := range ls {
		if !file.IsDir() {
			sh := newStaticHandler(filepath.Join(cwd, STATIC, file.Name()))
			http.Handle(fmt.Sprintf("/%s", file.Name()), sh)
			log.Printf("Registered handler for %s", file.Name())
		}
	}

	log.Printf("Listen on %d", PORT)
	die("ListenAndServe:", http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
