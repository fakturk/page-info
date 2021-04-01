package page

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

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

// Document represents an HTML document to be manipulated.
// It holds the root document node to manipulate, and can make selections on this document.
func GetDocumentFromURL(url string) (*goquery.Document, error) {

	var doc *goquery.Document
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in get url: ", err)
		return doc, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return doc, errors.New("Error Retrieving Document")
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

// add given information to the current response
func responseBuilder(current *string, title, info string) {

	*current += "\n" + title + ": " + info

}

// if url does not contain http or https http.Get does not work, because of that this function add http before url if not exists
func verifyURL(myUrl string) string {
	u, _ := url.Parse(myUrl)

	if u.Scheme != "" {
		return myUrl
	}
	return "http://" + myUrl
}

// gets the only host part of the url for checking internal links
func getBaseUrl(myurl string) string {

	u, err := url.Parse(myurl)
	if err != nil {
		panic(err)
	}

	return u.Host
}
