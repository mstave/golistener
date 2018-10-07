package main

import (
	"fmt"
	"log"
	"io/ioutil"
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
	fmt.Fprintf(w, "headers %s \n", r.Header["User-Agent"][0])
}

func remote(w http.ResponseWriter, r *http.Request) {
	remoteURL := "http://127.0.0.1:8082/local"
	response, err := http.Get(remoteURL)
	if err != nil {
	   fmt.Fprintf(w,"Error contacting remote at %s, err: %v \n", remoteURL, err)
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		w.Write(body)
	}
}

func main() {
	port := "8081"
	if (len(os.Args) == 2) {
		port = os.Args[1]
	} 
	http.HandleFunc("/local", logHandler(local))
	http.HandleFunc("/remote", logHandler(remote))
	log.Println("Listening on port", port)
    
	fmt.Println( http.ListenAndServe(":" + port, nil))
}