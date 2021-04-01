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
```
``` code 
func DetectHTMLTypeFromDoc(doc *goquery.Document) (string, error) 
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
```
``` code 
func separateLinks(baseURL string, hrefs []string) ([]string, []string)
```
``` code 
func getInaccessibleLinks(urls []string, wg *sync.WaitGroup)
```
``` code 
func checkUrl(url string, wg *sync.WaitGroup) 
```
``` code 
func checkLoginForm(doc *goquery.Document) bool 
```
Page Utils have following methods;

``` code 
func GetDocumentFromURL(url string) (*goquery.Document, error) 
```
``` code 
func responseBuilder(current *string, title, info string) 
```
``` code 
func verifyURL(myUrl string) string
```
``` code 
func getBaseUrl(myurl string) string 
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
    - example : localhost:8080/?url=facebook.com/login

      ![Get URL with POST](/images/post.png)