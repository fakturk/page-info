package page

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var amount uint64

func responseBuilder(current *string, title, info string) {
	fmt.Println(">>responseBuilder()")
	*current += "\n" + title + ": " + info

}

//TODO: correct this
func verifyURL(myUrl string) string {
	u, _ := url.Parse(myUrl)
	//fmt.Println(myUrl,u.Host+u.Path)
	if u.Scheme != "" {
		return myUrl
	}
	return "http://" + myUrl
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
	allLinks := extractLinks(doc)
	baseUrl := getBaseUrl(url)

	//for _, link := range allLinks {
	//	//fmt.Println(link)
	//	responseBuilder(&response, "link", link)
	//}
	fmt.Println("base url: ", baseUrl)
	internalLinks, externalLinks := separateLinks(baseUrl, allLinks)
	responseBuilder(&response, "Amount of internal links", strconv.Itoa(len(internalLinks)))
	responseBuilder(&response, "Amount of external links", strconv.Itoa(len(externalLinks)))
	var wg sync.WaitGroup
	wg.Add(2)
	go getInaccessibleLinks(internalLinks, &wg)
	//wg.Add(1)
	go getInaccessibleLinks(externalLinks, &wg)
	wg.Wait()
	//wg.Wait()
	fmt.Println("total amount:", atomic.LoadUint64(&amount))
	amountString := fmt.Sprint(atomic.LoadUint64(&amount))
	responseBuilder(&response, "Amount of inaccessible links", amountString)

	amount = 0
	loginFormStatus := checkLoginForm(doc)
	if loginFormStatus {
		responseBuilder(&response, "Login Form", "exists")
	} else {
		responseBuilder(&response, "Login Form", "not found")
	}
	//for _, link := range internalLinks {
	//	//fmt.Println(link)
	//	responseBuilder(&response, "internal link", link)
	//
	//}
	//externalLinks:=getExternalLinks(allLinks,internalLinks)
	//for _, link := range externalLinks {
	//	//fmt.Println(link)
	//	responseBuilder(&response, "external link", link)
	//
	//}
	return response
}

func checkLoginForm(doc *goquery.Document) bool {
	fmt.Println(">>checkLoginForm()")
	loginFormExist := false
	doc.Find("input[type='password']").Each(func(i int, selection *goquery.Selection) {
		//fmt.Println(selection.)
		fmt.Println("input field:", selection.Get(0))
		loginFormExist = true

	})
	fmt.Println(loginFormExist)
	fmt.Println("<<checkLoginForm()")
	return loginFormExist
}

func getBaseUrl(myurl string) string {
	fmt.Println(">>getBaseUrl()")
	u, err := url.Parse(myurl)
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Host)
	return u.Host
}
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

func getInaccessibleLinks(urls []string, wg *sync.WaitGroup) {
	fmt.Println(">>getInaccessibleLinks()")
	defer wg.Done()
	//var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)

		go checkUrl(url, wg)
		//wg.Done()
	}
	//wg.Wait()
	fmt.Println("<<getInaccessibleLinks()")
}
func checkUrl(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	url = verifyURL(url)
	//fmt.Println("checkURL:",url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in get url: ", err)
		atomic.AddUint64(&amount, 1)
		fmt.Println("amount :", atomic.LoadUint64(&amount))
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode, url)
		atomic.AddUint64(&amount, 1)
		fmt.Println("amount :", atomic.LoadUint64(&amount))
		return
	}
}
