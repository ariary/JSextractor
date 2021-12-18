package extract

import (
	"JSextractor/pkg/utils"
	"strings"

	"golang.org/x/net/html"
)

//Search for js in src attribute and return the Script struct associated if found
func FindJSinSrc(token html.Token, gather bool, domain string, line int) (s Script, found bool) {
	src := FindSrc(token)
	if src != "" {
		if gather {
			//retrieve JS from src attribute
			code := GatherJS(src, domain)
			s = Script{Line: line, Source: FromSrc, Content: code}
		} else {
			s = Script{Line: line, Source: FromSrc, Content: src}
		}
		found = true
	}
	return s, found
}

//Search for js in attribute event handler and return the Scripts struct associated if found
func FindJSinAttr(token html.Token, line int) (scripts []Script, found bool) {

	contents := FindJSEvent(token)
	if len(contents) > 0 {
		for i := 0; i < len(contents); i++ {
			s := Script{Line: line, Source: FromEvent, Content: contents[i]}
			scripts = append(scripts, s)
		}
		found = true
	}
	return scripts, found
}

//Search for js in attribute event handler and return the Scripts struct associated if found
func FindJSinTag(token html.Token, line int, tokenizer *html.Tokenizer) (s Script, found bool) {
	sType := GetScriptTagType(token)
	if sType == "" || strings.Contains(sType, "text/javascript") { //type="text/javascript;version=1.8" before firefox 59 was also accepted
		tokenizer.Next()
		s = Script{Line: line, Source: FromText, Content: string(tokenizer.Text())}
		found = true
	}
	return s, true
}

//Retrieve JS code from url (src attribut of script tag). use https by default. If url i s a relative path fetch [domain]/[url]
func GatherJS(url string, domain string) (code string) {
	switch {
	case strings.HasPrefix(url, "//"):
		code = utils.Fetch("https:" + url)
	case strings.HasPrefix(url, "http"): //handle also https
		code = utils.Fetch(url)
	default: //realtive
		code = utils.Fetch(domain + "/" + url)
	}
	return code
}

//Retrieve JS code in event attributes
func FindJSEvent(token html.Token) (contents []string) {
	for _, s := range token.Attr {
		if utils.Contains(eventJS[:], s.Key) {
			contents = append(contents, s.Val)
		}
	}
	return contents
}

//Find src attribute in tag
func FindSrc(token html.Token) string {
	for _, s := range token.Attr {
		if s.Key == "src" {
			return s.Val
		}
	}
	return ""
}
