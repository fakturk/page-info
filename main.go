package main

import (
	"fmt"
	"github.com/fakturk/page-info/page"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.SkipClean(true)
	router.HandleFunc("/", helloFunc).Methods("GET")
	router.HandleFunc("/url/{url}", page.FindUrlInfo).Methods("GET")
	router.HandleFunc("/", page.FindUrlInfoWithPost).Methods("POST")
	http.ListenAndServe(":8080", router)

}

func helloFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}
