package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const testURL = "https://www.adme.ru/svoboda-sdelaj-sam/34-poleznye-ssylki-kotorye-pomogut-vam-zaschitit-svoi-prava-i-vybratsya-iz-slozhnoj-situacii-1923515/"

func GetNode(_url string) (result *html.Node) {
	resp, err := http.Get(_url)
	if err != nil {
		fmt.Println("Parse url fail")
		return nil
	}
	defer resp.Body.Close()
	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Get html fail")
		return nil
	}
	result, err = html.Parse(strings.NewReader(string(doc)))
	return result
}

func GetParentDiv(_node *html.Node) (result *html.Node) {
	for result := _node; (result.Type != html.ElementNode || result.Data != "div") && result != nil; result = result.PrevSibling {

	}
	return result
}

func FindAd(_node *html.Node, _host string) {
	if _node.Type == html.ElementNode && _node.Data == "a" {
		for i := 0; i < len(_node.Attr); i++ {
			if _node.Attr[i].Key == "href" && !strings.Contains(_node.Attr[i].Val, _host) && len(_node.Attr[i].Val) > 2 /*&& strings.Contains(_node.Attr[i].Val, "http://")*/ {
				fmt.Printf("%v\n", _node.Attr[i].Val)
			}
		}
	}
	for node := _node.FirstChild; node != nil; node = node.NextSibling {
		FindAd(node, _host)
	}
}

func main() {
	purl, _ := url.Parse(testURL)
	host := purl.Host
	root := GetNode(testURL)
	if root != nil {
		FindAd(root, host)
	} else {
		fmt.Println("GetNode fail")
	}
}
