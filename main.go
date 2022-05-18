package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	defer func() {
		_ = resp.Close
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	log.Printf("%s", data)

	_, _ = w.Write([]byte("OK\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}
