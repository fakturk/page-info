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
