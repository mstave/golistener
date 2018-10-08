package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func logHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got a request for", r.URL.Path)
		f(w, r)
	}
}

func local(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "headers %s \n", r.Header["User-Agent"][0]); err != nil {
		log.Println(err)
	}
}

func remote(w http.ResponseWriter, _ *http.Request) {
	remoteURL := "http://127.0.0.1:8082/local"
	response, err := http.Get(remoteURL)
	if err != nil {
		if _, pErr := fmt.Fprintf(w, "Error contacting remote at %s, err: %v \n", remoteURL, err); pErr != nil {
			log.Println(pErr)
		}
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		if _, err := w.Write(body); err != nil {
			log.Println(err)
		}
	}
}

func serve() {
	port := "8081"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	r := mux.NewRouter()
	r.HandleFunc("/local", logHandler(local))
	r.HandleFunc("/remote", logHandler(remote))
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
func main() {
	serve()
}
