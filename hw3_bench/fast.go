package main

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson/jlexer"
	"io"
	"os"
	"strings"
)

type userStruct struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

func easyjsonC76e1e44DecodeGithubComGolangHw3BenchCodegen(in *jlexer.Lexer, out *userStruct) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *userStruct) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC76e1e44DecodeGithubComGolangHw3BenchCodegen(&r, v)
	return r.Error()
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	var user userStruct
	var isAndroid bool
	var isMSIE bool

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	seenBrowsers := make(map[string]int8, 100)
	scaner := bufio.NewScanner(file)
	i := -1

	fmt.Fprintln(out, "found users:")
	for scaner.Scan() {
		i++
		err = user.UnmarshalJSON(scaner.Bytes())
		if err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				seenBrowsers[browser] = 0
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
				seenBrowsers[browser] = 0
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", 1)
		fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", i, user.Name, email))
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
