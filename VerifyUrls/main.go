package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//res, err := http.Get("http://www.world-art.ru/")
	res, err := http.Get("https://golang.org/pkg/regexp/syntax/")
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
