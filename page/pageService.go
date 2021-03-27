package page

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func responseBuilder(current *string, title, info string) {
	fmt.Println(">>responseBuilder()")
	*current += "\n" + title + ": " + info

}
func verifyURL(url string) string {
	return "http://" + url
}

// Detect HTML Version From a http URL
func DetectHTMLTypeFromURL(url string) (string, error) {
	fmt.Println(">>DetectHTMLTypeFromURL()")
	var htmlVersion string

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in get url: ", err)
		return htmlVersion, err
	}
	//fmt.Println("resp: ",resp)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return htmlVersion, errors.New("Error Retrieving Document")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return htmlVersion, err
	}

	html, err := doc.Html()
	if err != nil {
		return htmlVersion, err
	}

	htmlVersion = checkDoctype(html)
	fmt.Println("HTML version: ", htmlVersion)

	return htmlVersion, nil

}

var doctypes = make(map[string]string)

func init() {
	doctypes["HTML 4.01 Strict"] = `"-//W3C//DTD HTML 4.01//EN"`
	doctypes["HTML 4.01 Transitional"] = `"-//W3C//DTD HTML 4.01 Transitional//EN"`
	doctypes["HTML 4.01 Frameset"] = `"-//W3C//DTD HTML 4.01 Frameset//EN"`
	doctypes["XHTML 1.0 Strict"] = `"-//W3C//DTD XHTML 1.0 Strict//EN"`
	doctypes["XHTML 1.0 Transitional"] = `"-//W3C//DTD XHTML 1.0 Transitional//EN"`
	doctypes["XHTML 1.0 Frameset"] = `"-//W3C//DTD XHTML 1.0 Frameset//EN"`
	doctypes["XHTML 1.1"] = `"-//W3C//DTD XHTML 1.1//EN"`
	doctypes["HTML 5"] = `<!DOCTYPE html>`
}
func checkDoctype(html string) string {
	var version = "UNKNOWN"

	for doctype, matcher := range doctypes {
		match := strings.Contains(html, matcher)

		if match == true {
			version = doctype
		}
	}

	return version
}
