package main

import (
	"JSextractor/pkg/config"
	"JSextractor/pkg/extract"
	"JSextractor/pkg/ui"
	"JSextractor/pkg/utils"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/awesome-gocui/gocui"
)

//suppose we have a correct HTML

func main() {
	// FLAG HANDLING
	gatherSrc := flag.Bool("gather-src", false, "Gather javascript code from script tag with src attribute. You must set url if you enabled it (-u flag)")
	url := flag.String("u", "", "url of the html page HTML page (eg https://example.net/home.html")
	skipSrc := flag.Bool("ds", false, "Do not search for javaScript code in src attributes with <script> tag")
	skipEvent := flag.Bool("de", false, "Do not search for javaScript code in event attributes")
	skipTag := flag.Bool("dt", false, "Do not search for javaScript code in <script> tag")
	tui := flag.Bool("tui", false, "terminal User Interface mode. Browse code in a full screen UI")
	flag.Parse()

	//transform url = https://[domain]/path/to/file -> https://[domain]/
	var domain string
	if *url != "" {
		domain = strings.Join(strings.SplitAfter(*url, "/")[:3], "")
	}
	cfg := config.Config{Url: domain, GatherSrc: *gatherSrc, SkipSrc: *skipSrc, SkipEvent: *skipEvent, SkipTag: *skipTag}

	if *gatherSrc && (*url == "") {
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
	if !*tui {
		for i := 0; i < len(scripts); i++ {
			extract.PrintScript(scripts[i])
		}
	} else {
		ui.Cfg = &cfg
		ui.Scripts = scripts
		ui.UpdateUiVars()

		ui.Cfg.Url = cfg.Url

		g, err := gocui.NewGui(gocui.OutputNormal, true)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()
		g.SetManagerFunc(ui.Layout)

		if err := ui.Keybindings(g); err != nil {
			log.Panicln(err)
		}

		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	}

}
