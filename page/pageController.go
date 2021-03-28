package page

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func FindUrlInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>findUrlInfo()")
	var params = mux.Vars(r)
	url := params["url"]
	response := getInfo(url)

	fmt.Fprintf(w, response)
}

func FindUrlInfoWithPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>findUrlInfoWithPost()")
	// we decode our body request params
	//var url string
	//_ = json.NewDecoder(r.Body).Decode(&url)

	//var params = mux.Vars(r)
	//url := params["url"]
	url := r.FormValue("url")
	response := getInfo(url)

	fmt.Fprintf(w, response)
}
