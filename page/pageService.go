package page

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var amount uint64

// collect the information about given web page and returns the result
func getInfo(url string) string {
	response := ""

	url = verifyURL(url)

	responseBuilder(&response, "URL", url)
	doc, _ := GetDocumentFromURL(url)

	//find html version of web page
	htmlVersion, err := DetectHTMLTypeFromDoc(doc)
	if err != nil {
		htmlVersion = "UNDEFINED due to error:" + err.Error()
	}

	responseBuilder(&response, "HTML Version", htmlVersion)

	// find title of web page
	title := getTitle(doc)
	responseBuilder(&response, "Title", title)

	//find each heading counts on the web page
	// we holds each count on a map but maps are unordered structures
	// because of that we sort keys of map (h1 to h6) and write to the response with order
	responseBuilder(&response, "Heading Counts", "")
	headings := getHeadingCount(doc)
	keys := make([]string, 0)
	for k, _ := range headings {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, heading := range keys {
		responseBuilder(&response, heading, strconv.Itoa(headings[heading]))
	}

	// find all links on a web page and separate those links to internal and external links
	allLinks := extractLinks(doc)
	baseUrl := getBaseUrl(url)

	internalLinks, externalLinks := separateLinks(baseUrl, allLinks)
	responseBuilder(&response, "Amount of internal links", strconv.Itoa(len(internalLinks)))
	responseBuilder(&response, "Amount of external links", strconv.Itoa(len(externalLinks)))

	// find inaccessible links
	// we use goroutines for checking each url, instead of waiting them one by one we call the urls concurrently
	// for counting on concurrent goroutines a sync.WaitGroup variable defined and add worker for each goroutine
	// we increase the inaccessible links amounts with atomic counter
	// and waitgroup wait all gouroutines to finish their job and writes result to the response
	var wg sync.WaitGroup
	wg.Add(2)
	go getInaccessibleLinks(internalLinks, &wg)

	go getInaccessibleLinks(externalLinks, &wg)
	wg.Wait()

	amountString := fmt.Sprint(atomic.LoadUint64(&amount))
	responseBuilder(&response, "Amount of inaccessible links", amountString)
	amount = 0

	// check if the login form exists
	loginFormStatus := checkLoginForm(doc)
	if loginFormStatus {
		responseBuilder(&response, "Login Form", "exists")
	} else {
		responseBuilder(&response, "Login Form", "not found")
	}

	return response
}

// Detect HTML Version From a http doc
func DetectHTMLTypeFromDoc(doc *goquery.Document) (string, error) {

	var htmlVersion string

	html, err := doc.Html()
	if err != nil {
		return htmlVersion, err
	}

	htmlVersion = checkDoctype(html)

	return htmlVersion, nil

}

// find doc type
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

//find title of the web page
func getTitle(doc *goquery.Document) string {
	title := doc.Find("title").Text()

	return title
}

//find each heading count
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

		})
		headings[heading] = count

	}

	return headings

}

// find all links on the web page by finding 'a' with 'href' attribute
func extractLinks(doc *goquery.Document) []string {
	foundUrls := []string{}
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			res, _ := s.Attr("href")
			foundUrls = append(foundUrls, res)
		})
		return foundUrls
	}
	return foundUrls
}

// seperate links to internal and external links
// if link starts with given urls host name (base url) or only '/' it is an internal link
// otherwise it is an external link
func separateLinks(baseURL string, hrefs []string) ([]string, []string) {
	internalUrls := []string{}
	externalUrls := []string{}

	for _, href := range hrefs {

		if strings.HasPrefix(href, baseURL) {
			internalUrls = append(internalUrls, href)
		} else if strings.HasPrefix(href, "/") {
			resolvedURL := fmt.Sprintf("%s%s", baseURL, href)
			internalUrls = append(internalUrls, resolvedURL)
		} else if href != "" {
			externalUrls = append(externalUrls, href)
		}
	}

	return internalUrls, externalUrls
}

// find inaccessible links
// we add a wait group worker before checking each url concurrently
// when it finish to check it give done signal to the wait group
func getInaccessibleLinks(urls []string, wg *sync.WaitGroup) {

	defer wg.Done()

	for _, url := range urls {
		wg.Add(1)

		go checkUrl(url, wg)

	}

}

// checks the given url is accessible or not
// if not accessible increase amount atomicly because we check urls in concurrent manner
// if response code is not 200 we assume that url is inaccessible
func checkUrl(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	url = verifyURL(url)

	resp, err := http.Get(url)
	if err != nil {

		atomic.AddUint64(&amount, 1)

		return
	}
	if resp.StatusCode != 200 {

		atomic.AddUint64(&amount, 1)

		return
	}
}

// checks if the login form exist or not by checking the input field with password type
func checkLoginForm(doc *goquery.Document) bool {

	loginFormExist := false
	doc.Find("input[type='password']").Each(func(i int, selection *goquery.Selection) {

		loginFormExist = true

	})

	return loginFormExist
}
