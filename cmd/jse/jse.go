package main

import (
	"JSextractor/pkg/config"
	"JSextractor/pkg/extract"
	"JSextractor/pkg/utils"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

//suppose we have a correct HTML

func main() {
	// FLAG HANDLING
	gatherSrc := flag.Bool("gather-src", false, "Gather javascript code from script tag with src attribute. You must set domain if you enabled it (-d flag)")
	domain := flag.String("d", "", "Domain hosting the HTML page (eg https://example.net")

	flag.Parse()

	cfg := config.Config{Url: *domain, GatherSrc: *gatherSrc}

	if *gatherSrc && (*domain == "") {
		log.Fatal("You must set domain if you enabled gathering js code from src (-gather-src) (-d flag)")
	}

	//RUN
	var buf bytes.Buffer
	tee := io.TeeReader(os.Stdin, &buf) //Read stdin twice

	//Get lines begin index
	page, _ := ioutil.ReadAll(tee)
	begins := utils.GetBeginLinesIndex(page)

	//Extract
	scripts := extract.Extract(&cfg, buf, begins)

	//Print result
	for i := 0; i < len(scripts); i++ {
		extract.PrintScript(scripts[i])
	}
}
