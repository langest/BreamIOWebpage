package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	PORT   = 8080
	STATIC = "static"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd: ", err)
	}
	static := http.Dir(filepath.Join(pwd, STATIC))
	log.Printf("Listen on %d and serving files in %s", PORT, static)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), http.FileServer(static))
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
