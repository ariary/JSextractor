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

	"github.com/jroimartin/gocui"
)

//suppose we have a correct HTML

func main() {
	// FLAG HANDLING
	gatherSrc := flag.Bool("gather-src", false, "Gather javascript code from script tag with src attribute. You must set domain if you enabled it (-d flag)")
	domain := flag.String("d", "", "Domain hosting the HTML page (eg https://example.net")
	skipSrc := flag.Bool("ds", false, "Do not search for javaScript code in src attributes with <script> tag")
	skipEvent := flag.Bool("de", false, "Do not search for javaScript code in event attributes")
	skipTag := flag.Bool("dt", false, "Do not search for javaScript code in <script> tag")
	tui := flag.Bool("tui", false, "terminal User Interface mode. Browse code in a full screen UI")
	flag.Parse()

	cfg := config.Config{Url: *domain, GatherSrc: *gatherSrc, SkipSrc: *skipSrc, SkipEvent: *skipEvent, SkipTag: *skipTag}

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
	if !*tui {
		for i := 0; i < len(scripts); i++ {
			extract.PrintScript(scripts[i])
		}
	} else {
		ui.Data = make(map[int]string)
		ui.Data[1] = "this item 1"
		ui.Data[2] = "this item 2"
		ui.Data[3] = "this item 3"
		ui.Data[4] = "this item 4"
		ui.Data[5] = "this item 5"
		ui.Data[6] = "this item 6"
		ui.Data[7] = "this item 7"

		ui.Scripts = scripts
		ui.UpdateUiVars()

		ui.Url = "https://www.deezer.com/fr/artist/8352118"

		g, err := gocui.NewGui(gocui.OutputNormal)
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
