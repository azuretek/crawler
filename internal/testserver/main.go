package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", DocRoot)
	http.HandleFunc("/about.html", About)
	http.ListenAndServe(":80", nil)
}

// DocRoot outputs the example1 fixture
func DocRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404")
		return
	}

	example1, err := ioutil.ReadFile("./fixtures/example1.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(example1))
}

// About outputs the example2 fixture
func About(w http.ResponseWriter, r *http.Request) {
	example2, err := ioutil.ReadFile("./fixtures/example2.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(example2))
}
