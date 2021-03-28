package page

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func responseBuilder(current *string, title, info string) {
	fmt.Println(">>responseBuilder()")
	*current += "\n" + title + ": " + info

}

//TODO: correct this
func verifyURL(url string) string {
	return "http://" + url
}
func GetDocumentFromURL(url string) (*goquery.Document, error) {
	fmt.Println(">>GetDocumentFromURL()")
	var doc *goquery.Document
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in get url: ", err)
		return doc, err
	}
	//fmt.Println("resp: ",resp)
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

func getInfo(url string) string {
	response := ""

	url = verifyURL(url)
	fmt.Println("url: " + url)
	responseBuilder(&response, "URL", url)
	doc, _ := GetDocumentFromURL(url)

	htmlVersion, err := DetectHTMLTypeFromDoc(doc)
	if err != nil {
		htmlVersion = "UNDEFINED due to error:" + err.Error()
	}
	fmt.Println(htmlVersion)
	responseBuilder(&response, "HTML Version", htmlVersion)
	title := getTitle(doc)
	responseBuilder(&response, "Title", title)
	responseBuilder(&response, "Heading Counts", "")
	headings := getHeadingCount(doc)
	keys := make([]string, 0)
	for k, _ := range headings {
		keys = append(keys, k)
		fmt.Println(k)
	}
	sort.Strings(keys)

	for _, heading := range keys {
		responseBuilder(&response, heading, strconv.Itoa(headings[heading]))
	}

	return response
}

// Detect HTML Version From a http doc
func DetectHTMLTypeFromDoc(doc *goquery.Document) (string, error) {
	fmt.Println(">>DetectHTMLTypeFromDoc()")
	var htmlVersion string

	html, err := doc.Html()
	if err != nil {
		return htmlVersion, err
	}

	htmlVersion = checkDoctype(html)
	fmt.Println("HTML version: ", htmlVersion)

	return htmlVersion, nil

}

func getTitle(doc *goquery.Document) string {
	title := doc.Find("title").Text()
	fmt.Println("title: ", title)
	return title
}

func getHeadingCount(doc *goquery.Document) map[string]int {
	headings := make(map[string]int)
	headings["h1"] = 0
	headings["h2"] = 0
	headings["h2"] = 0
	headings["h3"] = 0
	headings["h4"] = 0
	headings["h5"] = 0
	headings["h6"] = 0

	for heading, count := range headings {
		doc.Find(heading).Each(func(i int, selection *goquery.Selection) {
			count++
			//fmt.Println(selection.Contents())
		})
		headings[heading] = count
		fmt.Println("total", heading, ":", count)
	}

	return headings

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
