package main

import (
	"errors"
	"fmt"
	"github.com/saintfish/chardet"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// Struct for invalid URLs
type ErrUlr struct {
	Row     int
	Col     int
	Message string
}

func Scan(doc string) (urllist []ErrUlr) {
	defer func() {
		if err := recover(); err != nil {
			pc, _, _, _ := runtime.Caller(0)
			funcName := runtime.FuncForPC(pc).Name()
			log.Println(fmt.Sprintf("Panic happend in function %v: %v", funcName, err))
		}
	}()

	strs := strings.Split(doc, "\n")
	wg := sync.WaitGroup{}
	for row, str := range strs {
		wg.Add(1)
		go func(str string, row int) {
			defer wg.Done()
			regex := regexp.MustCompile(`<(img|a)\s.+?>`)
			tag := regex.FindAllString(str, -1)
			for _, t := range tag {
				col := strings.Index(str, t) + 1
				regex = regexp.MustCompile(`(href|src)=('|").{0,}?('|")`)
				attr := regex.FindAllString(str, -1)
				if len(attr) == 0 {
					urllist = append(urllist, ErrUlr{row, col, fmt.Sprintf("%s is not a valid Tag", t)})
					return
				}
				for _, a := range attr {
					pattern := `(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])`
					matched, err := regexp.Match(pattern, []byte(a[6:len(a)-1]))
					if err != nil {
						fmt.Println("regexp error: " + err.Error())
					}
					if !matched {
						urllist = append(urllist, ErrUlr{row, strings.Index(t, a) + col, fmt.Sprintf("%s is not a valid URL", a[5:])})
					}
				}
			}
		}(str, row)
	}
	wg.Wait()
	if len(urllist) > 0 {
		return urllist
	}
	return nil
}

// Returns error if the encoding does not match the one declared in the header
// Returns slice of struct with position and description of bad URL
func CheckPage(htm []byte, header http.Header) (error, []ErrUlr) {
	defer func() {
		if err := recover(); err != nil {
			pc, _, _, _ := runtime.Caller(0)
			funcName := runtime.FuncForPC(pc).Name()
			log.Println(fmt.Sprintf("Panic happend in function %v: %v", funcName, err))
		}
	}()
	//************************************************************************************************************************************************
	// Verify charset
	var resErr error                      // default error is nil
	ContentType := header["Content-Type"] // get content type from header
	charsetInHeader := ""

	// search for cahrset in content type
	for _, attr := range ContentType {
		if strings.Contains(attr, "charset") {
			charsetInHeader = strings.ToLower(strings.Split(attr, "=")[1])
		}
	}

	// analysis of html encoding
	charDetector := chardet.NewHtmlDetector()
	charsetDetected, err := charDetector.DetectBest(htm)
	if err != nil { //unknown error
		return err, nil
	}
	// compare detected charset with header and meta charsets
	if strings.ToLower(charsetDetected.Charset) != charsetInHeader {
		resErr = errors.New(fmt.Sprintf("Declared charset [%v], but detected charset is [%v]", charsetInHeader, charsetDetected.Charset))
	}
	//************************************************************************************************************************************************
	// Verify URLs
	urllist := make([]ErrUlr, 0, 100)        // list of invalid urls
	strs := strings.Split(string(htm), "\n") // split for row number

	// for each row search tags
	for row, str := range strs {
		// get all <a> and <img> tags
		regex := regexp.MustCompile(`<(img|a)\s.+?>`)
		tag := regex.FindAllString(str, -1)
		// for each tag check url
		for _, t := range tag {
			col := strings.Index(str, t) + 1 // remember col of tag, if tag have not url
			//fmt.Printf("[%v:%v] %v\n", row, col, t)
			// get attribute with URL
			regex = regexp.MustCompile(`(href|src)=('|").{0,}?('|")`)
			attr := regex.FindAllString(str, -1)
			if len(attr) == 0 { // if have not attribute, save error and continue
				//fmt.Printf("X %s is not a valid Tag\n", fmt.Sprintf("[%v:%v] %v", row, col, t))
				urllist = append(urllist, ErrUlr{row, col, fmt.Sprintf("%s is not a valid Tag", t)})
				continue
			}
			for _, a := range attr {
				//fmt.Printf("   [%v:%v] %v\n", row, strings.Index(t, a)+col, a[5:])
				// verify url in attribute
				pattern := `(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])`
				matched, err := regexp.Match(pattern, []byte(a[6:len(a)-1]))
				// any error message
				if err != nil {
					fmt.Println("regexp error: " + err.Error())
				}
				if !matched { // if URL is not valid, save error
					//fmt.Printf("X %s is not a valid URL\n", fmt.Sprintf("[%v:%v] %v", row, strings.Index(t, a)+col, a[5:]))
					urllist = append(urllist, ErrUlr{row, strings.Index(t, a) + col, fmt.Sprintf("%s is not a valid URL", a[5:])})
				} /* else {
					fmt.Printf("√ %s is a valid URL\n", fmt.Sprintf("[%v:%v] %v", row, strings.Index(t, a)+col, a[5:]))
				}*/
			}
		}
	}
	return resErr, urllist
}

func main() {
	res, err := http.Get("https://golang.org/pkg/regexp/syntax/") //http://www.world-art.ru/
	if err != nil {
		fmt.Printf("Get: %v\n\r", err.Error())
		return
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Read body: %v\n\r", err.Error())
	}
	err, urls := CheckPage(body, res.Header)
	if err != nil {
		fmt.Println(err.Error())
	}
	if urls != nil {
		for _, item := range urls {
			fmt.Println(item)
		}
	}
}
