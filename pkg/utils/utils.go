package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	Info = Teal
	Warn = Yellow
	Evil = Red
	Good = Green
	Code = Blue
)

var (
	Black         = Color("\033[1;30m%s\033[0m")
	Red           = Color("\033[1;31m%s\033[0m")
	Green         = Color("\033[1;32m%s\033[0m")
	Yellow        = Color("\033[1;33m%s\033[0m")
	Purple        = Color("\033[1;34m%s\033[0m")
	Magenta       = Color("\033[1;35m%s\033[0m")
	Teal          = Color("\033[1;36m%s\033[0m")
	White         = Color("\033[1;37m%s\033[0m")
	Blue          = Color("\033[1;96m%s\033[0m")
	Underlined    = Color("\033[4m%s\033[24m")
	Bold          = Color("\033[1m%s\033[0m")
	Italic        = Color("\033[3m%s\033[0m")
	RedForeground = Color("\033[1;41m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

//Return true if e is in s
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//Return the byte position of each line begin
func GetBeginLinesIndex(text []byte) (result []int) {
	lines := bytes.Split(text, []byte("\n"))
	offset := 0
	for i := 0; i < len(lines); i++ {
		offset += len(lines[i])
		result = append(result, offset)
	}

	return result
}

//Return body of url after performing GET request
func Fetch(url string) (body string, err error) {
	resp, err := http.Get(url)

	//We Read the response body on the line below.
	bodyB, err := ioutil.ReadAll(resp.Body)

	//Convert the body to type string
	body = string(bodyB)
	return body, err
}
