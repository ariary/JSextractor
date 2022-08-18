package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

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
	if err != nil {
		return "", err
	}

	//We Read the response body on the line below.
	bodyB, err := ioutil.ReadAll(resp.Body)

	//Convert the body to type string
	body = string(bodyB)
	return body, err
}

//Return the response of the curl command "cmd"
func Curl(cmd string) (output string, err error) {
	args := strings.Split(cmd, " ")[1:] //withdraw curl
	curl := exec.Command("curl", args...)
	out, err := curl.Output()
	output = string(out)
	return output, err
}
