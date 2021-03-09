package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Position invalid URL in HTML
type Position struct {
	Row int
	Col int
}

// Struct for invalid URLs
type ErrUlr struct {
	Pos     Position
	Message string
}

func (p *ErrUlr) String() string {
	return fmt.Sprintf("[%v:%v] %v", p.Pos.Row, p.Pos.Col, p.Message)
}

// Search all positions of the URL
func GetPositions(text []string, substr string) (res []Position) {
	for i, str := range text {
		if strings.Contains(str, substr) {
			res = append(res, Position{i, strings.Index(str, substr)})
		}
	}
	return
}

// Check charset and search all invalid URLs
func CheckPage(htm []byte, header http.Header) (resErr error, resArr []ErrUlr) {
	resErr = nil                          // default error is nil
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
	//************************************************************************************************************************************************
	// create reader for search urls
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htm)))
	if err != nil {
		log.Fatal(err)
	}
	//using for search positions as [row:col]
	htmlStrings := strings.Split(string(htm), "\n")

	//using for collect bad urls
	mapErrUrl := make(map[string]ErrUlr)
	//get all tags which contains urls
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		checkUrl, ex := s.Attr("href")
		if ex {
			if !IsValidUrl(checkUrl) { // validate url
				//search all positions for current url
				for _, pos := range GetPositions(htmlStrings, checkUrl) {
					item := ErrUlr{pos, "Invalid URL: " + checkUrl}
					mapErrUrl[checkUrl] = item // save it to map
				}
			}
		} else {
			//fmt.Println(s.)
		}

	})

	charsetInMeta := ""
	// searching charset in meta tag
	doc.Find("[charset]").Each(func(i int, s *goquery.Selection) {
		charsetInMeta, _ = s.Attr("charset")
	})

	// compare detected charset with header and meta charsets
	if (strings.ToLower(charsetDetected.Charset) != charsetInHeader) || ((charsetInMeta != "") && (strings.ToLower(charsetDetected.Charset) != charsetInMeta)) {
		resErr = errors.New(fmt.Sprintf("Declared charset [Header: %v] [Meta: %v], but detected charset is [%v]", charsetInHeader, charsetInMeta, charsetDetected.Charset))
	}
	// if bad urls exists
	if len(mapErrUrl) > 0 {
		resArr = make([]ErrUlr, 0, len(mapErrUrl))
		for _, item := range mapErrUrl {
			resArr = append(resArr, item)
		}
		return resErr, resArr
	}
	return resErr, nil
}

// validate url
func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func main() {
	test := `
<html>
<head>
<title>Мировое искусство: живопись, литература, аниме, кино</title>
<LINK REL='SHORTCUT ICON' HREF='http://www.world-art.ru/favicon.ico'>
<link href='style.css' type='text/css' rel='stylesheet'>
<meta http-equiv='Content-Type' content='text/html; charset=windows-1251'>
<meta http-equiv='expires' content='Mon, 01 Jan 1990 00:00:00 GMT'> 
</head>

<body bottomMargin='0' leftMargin='0' topMargin='0' rightMargin='0' marginwidth='0' marginheight='0'>

<center>
<table bgcolor=#990000 width=1004 cellpadding=0 cellspacing=0 border=0 height=75>
<tr>
<td width=5></td>
<td Valign=top>&nbsp;&nbsp;&nbsp;&nbsp;<img src='http://www.world-art.ru/img/logo.gif' alt='World Art - сайт о кино, сериалах, литературе, аниме, играх, живописи и архитектуре.' width=213 height=59 border=0>
</td>
<form action='http://www.world-art.ru/search.php' method='get'>
<td align=right>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<font color=#ffffff><b>поиск:</b></font> 
<input class='web' name='public_search' value='' type='text' style='width:192;'><br>
<font color=#ffffff><b>в разделе:</b> 
<select name='global_sector' style='width:138px; margin-top:2px;'>
<OPTION value='all'>по всему сайту</option><OPTION value='all'>------------</option><OPTION value='animation'>аниме</option><OPTION value='architecture'>архитектура</option><OPTION value='games'>видеоигры</option><OPTION value='cinema'>кино</option><OPTION value='lyric'>литература</option><OPTION value='manga'>манга</option><OPTION value='painting'>живопись</option><OPTION value='people'>персоны</option><OPTION value='company'>компании</option>
</select>
<input type=submit value='Поиск' style='width:50; font-family: Verdana; font-size: 12px; border:1px; padding: 1px 0px 1px 0px; margin-top:1px;'>
</td>
<td width=5></td>
</form>
</tr>
</table>

<table height=1 width=1004 cellpadding=0 cellspacing=0 border=0 bgcolor=#5D0E0E>
<tr>
<td></td>
</tr>
</table>

<table height=29 width=1004 border=0 bgcolor=#781111 cellpadding=0 cellspacing=0 border=0>
<tr>
<td width=12></td>

<td width=42>&nbsp;&nbsp;<font color='ffffff'><b><a href='htp://www.world-art.ru/cinema/' class='main_page'>Кино</a>&nbsp;&nbsp;</td>
<td width=1 bgcolor=#5D0E0E></td>

<td width=40>&nbsp;&nbsp;<font color='ffffff'><b><a class='main_page'>Аниме</a>&nbsp;&nbsp;</td>
<td width=1 bgcolor=#5D0E0E></td>

<td width=40>&nbsp;&nbsp;<font color='ffffff'><b><a href='http://www.world-art.ru/games/' class='main_page'>Видеоигры</a>&nbsp;&nbsp;</td>
<td width=1 bgcolor=#5D0E0E></td>

<td width=40>&nbsp;&nbsp;<font color='ffffff'><b><a href='http://www.world-art.ru/lyric/' class='main_page'>Литература</a>&nbsp;&nbsp;</td>
<td width=1 bgcolor=#5D0E0E></td>

<td width=40>&nbsp;&nbsp;<font color='ffffff'><b><a href='http://www.world-art.ru/painting/' class='main_page'>Живопись</a>&nbsp;&nbsp;</td>
<td width=1 bgcolor=#5D0E0E></td>

<td width=40>&nbsp;&nbsp;<font color='ffffff'><b><a href='http://www.world-art.ru/architecture/' class='main_page'>Архитектура</a>&nbsp;&nbsp;</td>
<td width=1 bgcolor=#5D0E0E></td>

<td align=right><b><a href='http://www.world-art.ru/enter.php' class='main_page'>Вход в систему</a></b>&nbsp;</td>
<td width=1 bgcolor=#5D0E0E></td>
<td width=55><b>&nbsp;&nbsp;<a href='http://www.world-art.ru/regstart.php' class='main_page'>Регистрация</a></b>&nbsp;&nbsp;</td>

</tr>

</table>
</body></html>

`

	res, err := http.Get("http://world-art.ru")
	if err != nil {
		fmt.Printf("Get: %v\n\r", err.Error())
		return
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	b, _ := ioutil.ReadAll(strings.NewReader(test)) //res.Body)

	err, urls := CheckPage(b, res.Header)
	if err != nil {
		fmt.Println(err.Error())
	}
	if urls != nil {
		for _, item := range urls {
			fmt.Println(item.String())
		}
	}
}
