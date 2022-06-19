package main

import (
	"fmt"
	"github.com/aapelismith/kun/pkg/version"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%+v", version.Get())
	})

	log.Printf("http listen and server at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
