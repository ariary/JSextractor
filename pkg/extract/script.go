package extract

import (
	"fmt"
	"log"
	"os"

	"github.com/ariary/go-utils/pkg/color"
	"golang.org/x/net/html"
)

//Type represent the script source type
type Type int

const (
	FromText Type = iota
	FromSrc
	FromSrcGathered
	FromEvent
)

//Script structure containing the source , the code, and the line where it appears in source code
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
	case FromSrcGathered:
		return "in src attribute " + color.Green("✔")
	case FromEvent:
		return "in event handlers"
	}
	return ""
}

func ScriptInfoOutput(s Script) string {
	return color.Bold(color.Red(s.Line)) + " : " + color.Blue(s.Source.String())
}

func PrintScript(s Script) {
	l := log.New(os.Stderr, "", 0) //write to stderr to don't have it if you redirect output
	info := ScriptInfoOutput(s)
	l.Println(info)
	if s.Content != "" {
		output := s.Content //TODO: pass s.Content to js beautifier
		fmt.Println(output)
	}
}

//GetScriptTagType Return the value of the "type" attribute for <script> tag
func GetScriptTagType(token html.Token) string {
	for _, s := range token.Attr {
		if s.Key == "type" {
			return s.Val
		}
	}
	return ""
}
