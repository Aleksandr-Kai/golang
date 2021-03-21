package main

import (
	"errors"
	"fmt"
	"github.com/saintfish/chardet"
	"net/http"
	"regexp"
	"strings"
)


// Struct for invalid URLs
type BadUrlMsg struct {
	Row     int
	Col     int
	Message string
}

// Accepts a string with html code as input and returns a slice of structures with invalid URLs
func findBadURLs(doc string) (err error, result []BadUrlMsg) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	result = make([]BadUrlMsg, 0, 100) // list of invalid urls
	rows := strings.Split(doc, "\n")   // split for row number
	regexTag := regexp.MustCompile(`<(img|a)\s.+?>`)
	regexAttr := regexp.MustCompile(`(href|src)=('|").{0,}?('|")`)
	patternURL := `(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])`
	// for each row search tags
	for rowNum, rowText := range rows {
		// get all <a> and <img> tags
		htmlTags := regexTag.FindAllString(rowText, -1)
		// for each tag check url
		for _, tag := range htmlTags {
			colNum := strings.Index(rowText, tag) + 1 // remember col of tag, if tag have not url
			// get attribute with URL
			attributes := regexAttr.FindAllString(rowText, -1)
			if len(attributes) == 0 { // if have not attribute, save error and continue
				result = append(result, BadUrlMsg{rowNum + 1, colNum, fmt.Sprintf("%s is not a valid Tag", tag)})
				continue
			}
			for _, attr := range attributes {
				// verify url in attribute
				matched, err := regexp.Match(patternURL, []byte(attr[6:len(attr)-1]))
				if err != nil {
					return err, nil
				}
				if !matched { // if URL is not valid, save error
					result = append(result, BadUrlMsg{rowNum + 1, strings.Index(tag, attr) + colNum, fmt.Sprintf("%s is not a valid URL", attr[5:])})
				}
			}
		}
	}
	if len(result) > 0 {
		return nil, result
	}
	return nil, nil
}

// Returns error if the encoding does not match the one declared in the header
// Returns slice of struct with position and description of bad URL
func CheckPage(data []byte, header http.Header) (charsetErr error, badURLs []BadUrlMsg) {
	defer func() {
		if err := recover(); err != nil {
			charsetErr = err.(error)
		}
	}()
	// Verify charset
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
	charsetDetected, err := charDetector.DetectBest(data)
	if err != nil { //unknown error
		charsetErr = err
	}
	// compare detected charset with header and meta charsets
	if strings.ToLower(charsetDetected.Charset) != charsetInHeader {
		charsetErr = errors.New(fmt.Sprintf("Declared charset [%v], but detected charset is [%v]", charsetInHeader, charsetDetected.Charset))
	}
	err, badUrls := findBadURLs(string(data))
	if err != nil {
		return err, nil
	}
	return charsetErr, badUrls
}
