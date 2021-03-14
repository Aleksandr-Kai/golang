package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	validURLs := `
<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href='http://google.com/' class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://google.com" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="https://google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://www.www.google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
`
	invalidURLs := `
<a href="http://go#ogle.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="google.com" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://google." class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://google/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="#start-of-content" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http:/google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="htp://google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http//google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
<a href="http://#google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
`
	fmt.Println("Test valid URLs")
	res := Scan(validURLs)
	if res != nil {
		for _, r := range res {
			fmt.Printf("[%v:%v] %v\n", r.Row, r.Col, r.Message)
		}
		t.Error("Test fail")
	} else {
		fmt.Println("Test pass")
	}

	fmt.Println("Test bad URLs")
	res = Scan(invalidURLs)
	l := len(strings.Split(invalidURLs, "\n")) - 2
	if len(res) != l {
		for _, r := range res {
			fmt.Printf("[%v:%v] %v\n", r.Row, r.Col, r.Message)
		}
		t.Error(fmt.Sprintf("Test pass: [expected %v] vs [result %v]", len(res), l))
	} else {
		fmt.Println("Test pass")
	}
}
