package main

import (
	"fmt"
	"testing"
)

func TestScan(t *testing.T) {
	testHTML := `
1<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
2<a href='http://google.com/' class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
3<a href="http://google.com" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
4<a href="http://google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
5<a href="https://google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
6<a href="http://www.www.google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
7<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
8<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
9<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
10<a href="http://google.com/?sdfsd" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
11<a href="google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
12<a href="google.com" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
13<a href="http://google." class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
14<a href="http://google/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
15<a href="#start-of-content" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
16<a href="http:/google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
17<a href="htp://google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
18<a href="http//google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
19<a class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
20<a href="http://#google.com/" class="p-3 color-bg-info-inverse color-text-white show-on-focus js-skip-to-content">Skip to content</a>
`
	res := Scan(testHTML)
	for i, r := range res {
		fmt.Printf("[%v:%v] %v\n", r.Row, r.Col, r.Message)
		if i+11 != r.Row {
			t.Error("Test fail", i, " : ", r)
		}
	}
}
