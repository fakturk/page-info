package page

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func FindUrlInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>findUrlInfo()")
	response := ""
	var params = mux.Vars(r)
	url := params["url"]
	url = verifyURL(url)
	fmt.Println("url: " + url)
	responseBuilder(&response, "URL", url)

	htmlVersion, err := DetectHTMLTypeFromURL(url)
	if err != nil {
		htmlVersion = "UNDEFINED due to error:" + err.Error()
	}
	fmt.Println(htmlVersion)
	responseBuilder(&response, "HTML Version", htmlVersion)
	fmt.Fprintf(w, response)
}
