package extract

import (
	"JSextractor/pkg/utils"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

type Type int

const (
	FromText Type = iota
	FromSrc
	FromEvent
)

//TODO: Add from "javascript:" in attr

type Script struct {
	Source  Type
	Content string
	Line    int
}

func (t *Type) String() string {
	switch *t {
	case FromText:
		return "in <script> tag"
	case FromSrc:
		return "in src attribute"
	case FromEvent:
		return "in event handlers"
	}
	return ""
}

func PrintScript(s Script) {
	l := log.New(os.Stderr, "", 0) //write to stderr to don't have it if you redirect output
	info := utils.Bold(utils.Red(s.Line)) + " : " + utils.Blue(s.Source.String())
	l.Println(info)
	if s.Content != "" {
		output := s.Content //TODO: pass s.Content to js beautifier
		fmt.Println(output)
	}
}

//Return the value of the "type" attribute for <script> tag
func GetScriptTagType(token html.Token) string {
	for _, s := range token.Attr {
		if s.Key == "type" {
			return s.Val
		}
	}
	return ""
}
