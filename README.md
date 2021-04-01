# Page-Info

## Purpose
Give information about a web page on plain text

Takes a website URL as an input and provides general information
about the contents of the page:

- [X] HTML Version
- [X]  Page Title
- [X]  Headings count by level
- [X]  Amount of internal and external links
- [X] Amount of inaccessible links
- [X] If a page contains a login form

## Architecture
We have page service, controller and utils (helper)

Page Controller have following methods;

``` code 
func FindUrlInfo(w http.ResponseWriter, r *http.Request)
```
``` code 
func FindUrlInfoWithPost(w http.ResponseWriter, r *http.Request) 
```

Page Service have following methods;

``` code 
func getInfo(url string) string 

// collect the information about given web page and returns the result
```

``` code 
func DetectHTMLTypeFromDoc(doc *goquery.Document) (string, error)

// Detect HTML Version From a http doc 
```
``` code 
func checkDoctype(html string) string 
```
``` code 
func getTitle(doc *goquery.Document) string 
```
``` code 
func getHeadingCount(doc *goquery.Document) map[string]int 
```
``` code 
func extractLinks(doc *goquery.Document) []string 

// find all links on the web page by finding 'a' with 'href' attribute

```
``` code 
func separateLinks(baseURL string, hrefs []string) ([]string, []string)

// seperate links to internal and external links
// if link starts with given urls host name (base url) or only '/' it is an internal link
// otherwise it is an external link
```
``` code 
func getInaccessibleLinks(urls []string, wg *sync.WaitGroup)

// find inaccessible links
// we add a wait group worker before checking each url concurrently
// when it finish to check it give done signal to the wait group
// we use goroutines for checking each url, instead of waiting them one by one we call the urls concurrently
// for counting on concurrent goroutines a sync.WaitGroup variable defined and add worker for each goroutine
// we increase the inaccessible links amounts with atomic counter
// and waitgroup wait all gouroutines to finish their job and writes result to the response
	
```
``` code 
func checkUrl(url string, wg *sync.WaitGroup) 

// checks the given url is accessible or not
// if not accessible increase amount atomicly because we check urls in concurrent manner
// if response code is not 200 we assume that url is inaccessible

```
``` code 
func checkLoginForm(doc *goquery.Document) bool 

// checks if the login form exist or not by checking the input field with password type
```
Page Utils have following methods;

``` code 
func GetDocumentFromURL(url string) (*goquery.Document, error) 

// Document represents an HTML document to be manipulated.
// It holds the root document node to manipulate, and can make selections on this document.
```
``` code 
func responseBuilder(current *string, title, info string) 

// add given information to the current response
```
``` code 
func verifyURL(myUrl string) string

// if url does not contain http or https http.Get does not work, because of that this function add http before url if not exists
```
``` code 
func getBaseUrl(myurl string) string 

// gets the only host part of the url for checking internal links
```
## How to Run

```shell
$ go run main.go
```

## Usage
- Get Url:
    - request type: GET
    - host: localhost:8080/url/{url}
    - works only with hostname without path, if the url has path POST method should use
    - example : http://localhost:8080/url/facebook.com
    
        ![Get URL](/images/get.png)
    

- Get Url with POST:
    - request type: POST
    - host: localhost:8080/?url={urlWithPath}
    - query parameters : url
    - example : localhost:8080/?url=facebook.com/login

      ![Get URL with POST](/images/post.png)