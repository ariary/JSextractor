package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//https://schier.co/blog/a-simple-web-scraper-in-go
//suppose we have a correct HTML

type Type int

/////////////////UTILS//////////////////
///////////////////////////////////////
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
func Fetch(url string) (body string) {
	resp, err := http.Get(url)
	if err != nil {
		os.Stderr.WriteString("Failed fetching url: " + url)
	}
	//We Read the response body on the line below.
	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Stderr.WriteString("Failed reading body of GET response from: " + url)
	}
	//Convert the body to type string
	body = string(bodyB)
	return body
}

/////////////////SCRIPT////////////////
///////////////////////////////////////
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
	info := Bold(Red(s.Line)) + " : " + Blue(s.Source.String())
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

//Retrieve JS code from url (src attribut of script tag). use https by default. If url i s a relative path fetch [domain]/[url]
func GatherJS(url string, domain string) (code string) {
	switch {
	case strings.HasPrefix(url, "//"):
		code = Fetch("https:" + url)
	case strings.HasPrefix(url, "http"): //handle also https
		code = Fetch(url)
	default: //realtive
		code = Fetch(domain + "/" + url)

	}
	return code
}

/////////////TAG HANDLING//////////////
///////////////////////////////////////
var eventJS = [...]string{
	"onactivate",
	"onafterprint",
	"onafterscriptexecute",
	"onanimationcancel",
	"onanimationend",
	"onanimationiteration",
	"onanimationstart",
	"onauxclick",
	"onbeforeactivate",
	"onbeforecopy",
	"onbeforecut",
	"onbeforedeactivate",
	"onbeforepaste",
	"onbeforeprint",
	"onbeforescriptexecute",
	"onbeforeunload",
	"onbegin",
	"onblur",
	"onbounce",
	"oncanplay",
	"oncanplaythrough",
	"onchange",
	"onclick",
	"onclose",
	"oncontextmenu",
	"oncopy",
	"oncuechange",
	"oncut",
	"ondblclick",
	"ondeactivate",
	"ondrag",
	"ondragend",
	"ondragenter",
	"ondragleave",
	"ondragover",
	"ondragstart",
	"ondrop",
	"ondurationchange",
	"onend",
	"onended",
	"onerror",
	"onfinish",
	"onfocus",
	"onfocusin",
	"onfocusout",
	"onfullscreenchange",
	"onhashchange",
	"oninput",
	"oninvalid",
	"onkeydown",
	"onkeypress",
	"onkeyup",
	"onload",
	"onloadeddata",
	"onloadedmetadata",
	"onloadend",
	"onloadstart",
	"onmessage",
	"onmousedown",
	"onmouseenter",
	"onmouseleave",
	"onmousemove",
	"onmouseout",
	"onmouseover",
	"onmouseup",
	"onmousewheel",
	"onmozfullscreenchange",
	"onpagehide",
	"onpageshow",
	"onpaste",
	"onpause",
	"onplay",
	"onplaying",
	"onpointerdown",
	"onpointerenter",
	"onpointerleave",
	"onpointermove",
	"onpointerout",
	"onpointerover",
	"onpointerrawupdate",
	"onpointerup",
	"onpopstate",
	"onprogress",
	"onreadystatechange",
	"onrepeat",
	"onreset",
	"onresize",
	"onscroll",
	"onsearch",
	"onseeked",
	"onseeking",
	"onselect",
	"onselectionchange",
	"onselectstart",
	"onshow",
	"onstart",
	"onsubmit",
	"ontimeupdate",
	"ontoggle",
	"ontouchend",
	"ontouchmove",
	"ontouchstart",
	"ontransitioncancel",
	"ontransitionend",
	"ontransitionrun",
	"ontransitionstart",
	"onunhandledrejection",
	"onunload",
	"onvolumechange",
	"onwaiting",
	"onwebkitanimationend",
	"onwebkitanimationiteration",
	"onwebkitanimationstart",
	"onwebkittransitionend",
	"onwheel",
}

func FindJSEvent(token html.Token) (contents []string) {
	for _, s := range token.Attr {
		if Contains(eventJS[:], s.Key) {
			contents = append(contents, s.Val)
		}
	}
	return contents
}

//find src attribute in tag
func FindSrc(token html.Token) string {
	for _, s := range token.Attr {
		if s.Key == "src" {
			return s.Val
		}
	}
	return ""
}

///////////////////MAIN////////////////
///////////////////////////////////////
func main() {
	// FLAG HANDLING
	gatherSrc := flag.Bool("gather-src", false, "Gather javascript code from script tag with src attribute. You must set domain if you enabled it (-d flag)")
	domain := flag.String("d", "", "Domain hosting the HTML page (eg https://example.net")

	flag.Parse()

	if *gatherSrc && (*domain == "") {
		log.Fatal("You must set domain if you enabled gathering js code from src (-gather-src) (-d flag)")
	}

	//RUN
	offset := 0 //snoop line
	line := 0

	var buf bytes.Buffer
	tee := io.TeeReader(os.Stdin, &buf) //Read stdin twice
	scripts := []Script{}
	tokenizer := html.NewTokenizer(&buf)

	page, _ := ioutil.ReadAll(tee)
	begins := GetBeginLinesIndex(page)

	var readAll bool

	for {
		tokenType := tokenizer.Next()
		offset += len(tokenizer.Raw())

		for i := 0; i < len(begins); i++ {
			if offset > begins[i] {
				line = i
			}
		}
		offset += len(tokenizer.Raw())

		switch {
		case tokenType == html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				readAll = true //break statement won't work
			} else {
				log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			}
		case tokenType == html.SelfClosingTagToken:
			token := tokenizer.Token()

			isScript := token.Data == "script"
			if isScript {
				//src finder
				src := FindSrc(token)
				if src != "" {
					var s Script
					if *gatherSrc {
						//retrieve JS from src attribute
						code := GatherJS(src, *domain)
						s = Script{Line: line, Source: FromSrc, Content: code}
					} else {
						s = Script{Line: line, Source: FromSrc, Content: src}
					}
					scripts = append(scripts, s)
				}
				break
			}

			//Find in attr
			contents := FindJSEvent(token)
			if len(contents) > 0 {
				for i := 0; i < len(contents); i++ {
					s := Script{Line: line, Source: FromEvent, Content: contents[i]}
					scripts = append(scripts, s)
				}
				break
			}
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			isScript := token.Data == "script"
			if isScript {
				//src finder
				src := FindSrc(token)
				if src != "" {
					var s Script
					if *gatherSrc {
						//retrieve JS from src attribute
						code := GatherJS(src, *domain)
						s = Script{Line: line, Source: FromSrc, Content: code}
					} else {
						s = Script{Line: line, Source: FromSrc, Content: src}
					}
					scripts = append(scripts, s)
					break
				}

				// between tag finder
				sType := GetScriptTagType(token)
				if sType == "" || strings.Contains(sType, "text/javascript") { //type="text/javascript;version=1.8" before firefox 59 was also accepted
					tokenizer.Next()
					s := Script{Line: line, Source: FromText, Content: string(tokenizer.Text())}
					scripts = append(scripts, s)
					break
				}
			}

			//Find in attr
			contents := FindJSEvent(token)
			if len(contents) > 0 {
				for i := 0; i < len(contents); i++ {
					s := Script{Line: line, Source: FromEvent, Content: contents[i]}
					scripts = append(scripts, s)
				}
			}
		}

		//Exit for loop if you have read the whole input
		if readAll {
			break
		}
	}

	//Print result
	for i := 0; i < len(scripts); i++ {
		PrintScript(scripts[i])
	}
}
