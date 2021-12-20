package extract

import (
	"JSextractor/pkg/config"
	"JSextractor/pkg/utils"
	"bytes"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// Extract all js from a buffer (representing HTML content). begins are the offset psition of byte starting a new line
func Extract(cfg *config.Config, buf bytes.Buffer, begins []int) (scripts []Script) {
	var readAll bool

	offset := 0 //snoop line
	line := 0

	tokenizer := html.NewTokenizer(&buf)

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
			if isScript && !cfg.SkipSrc {
				//src finder
				s, found := FindJSinSrc(token, cfg.GatherSrc, cfg.Url, line)
				if found {
					scripts = append(scripts, s)
				}
				break
			}

			//Find in attr
			if !cfg.SkipEvent {
				sL, found := FindJSinAttr(token, line)
				if found {
					scripts = append(scripts, sL...)
					break
				}
			}
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			isScript := token.Data == "script"
			if isScript {
				var s Script
				var found bool

				//src finder
				if !cfg.SkipSrc {
					s, found = FindJSinSrc(token, cfg.GatherSrc, cfg.Url, line)
					if found {
						scripts = append(scripts, s)
						break
					}
				}

				// between tag finder
				if !cfg.SkipTag {
					s, found = FindJSinTag(token, line, tokenizer)
					if found {
						scripts = append(scripts, s)
						break
					}
				}

			}

			//Find in attr
			if !cfg.SkipEvent {
				sL, found := FindJSinAttr(token, line)
				if found {
					scripts = append(scripts, sL...)
					break
				}
			}
		}

		//Exit for loop if you have read the whole input
		if readAll {
			break
		}
	}
	return scripts
}

//Search for js in src attribute and return the Script struct associated if found
func FindJSinSrc(token html.Token, gather bool, domain string, line int) (s Script, found bool) {
	src := FindSrc(token)
	if src != "" {
		if gather {
			//retrieve JS from src attribute
			code, err := GatherJS(src, domain)
			if err != nil {
				s = Script{Line: line, Source: FromSrc, Content: code}
			} else {
				s = Script{Line: line, Source: FromSrc, Content: src + " (failed to retrieve code by fetching src)"}
			}
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

//Retrieve JS code from url (src attribut of script tag). use https by default. If url is a relative path -> fetch [domain]/[url]
func GatherJS(url string, domain string) (code string, err error) {
	switch {
	case strings.HasPrefix(url, "//"):
		code, err = utils.Fetch("https:" + url)
	case strings.HasPrefix(url, "http"): //handle also https
		code, err = utils.Fetch(url)
	default: //realtive
		code, err = utils.Fetch(domain + "/" + url)
	}
	return code, err
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
