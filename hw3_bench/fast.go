package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type userStruct struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
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
		err := json.Unmarshal(scaner.Bytes(), &user)
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

	//fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
