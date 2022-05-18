package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}
