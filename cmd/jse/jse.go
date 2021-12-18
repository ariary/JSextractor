package main

import (
	"JSextractor/pkg/extract"
	"JSextractor/pkg/utils"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/html"
)

//suppose we have a correct HTML

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
	scripts := []extract.Script{}
	tokenizer := html.NewTokenizer(&buf)

	//Get lines begin index
	page, _ := ioutil.ReadAll(tee)
	begins := utils.GetBeginLinesIndex(page)

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
				s, found := extract.FindJSinSrc(token, *gatherSrc, *domain, line)
				if found {
					scripts = append(scripts, s)
				}
				break
			}

			//Find in attr
			sL, found := extract.FindJSinAttr(token, line)
			if found {
				scripts = append(scripts, sL...)
				break
			}
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			isScript := token.Data == "script"
			if isScript {
				//src finder
				s, found := extract.FindJSinSrc(token, *gatherSrc, *domain, line)
				if found {
					scripts = append(scripts, s)
					break
				}

				// between tag finder
				s, found = extract.FindJSinTag(token, line, tokenizer)
				if found {
					scripts = append(scripts, s)
					break
				}
			}

			//Find in attr
			sL, found := extract.FindJSinAttr(token, line)
			if found {
				scripts = append(scripts, sL...)
				break
			}
		}

		//Exit for loop if you have read the whole input
		if readAll {
			break
		}
	}

	//Print result
	for i := 0; i < len(scripts); i++ {
		extract.PrintScript(scripts[i])
	}
}
