package main

import (
	"JSextractor/pkg/config"
	"JSextractor/pkg/extract"
	"JSextractor/pkg/ui"
	"JSextractor/pkg/utils"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ariary/quicli/pkg/quicli"
	"github.com/awesome-gocui/gocui"
)

//suppose we have a correct HTML

func main() {
	// FLAG HANDLING
	// gatherSrc := flag.Bool("gather-src", false, "Gather javascript code from script tag with src attribute. You must set url if you enabled it (-u flag)")
	// url := flag.String("u", "", "url of the html page HTML page (eg https://example.net/home.html)")
	// skipSrc := flag.Bool("ds", false, "Do not search for javaScript code in src attributes with <script> tag")
	// skipEvent := flag.Bool("de", false, "Do not search for javaScript code in event attributes")
	// skipTag := flag.Bool("dt", false, "Do not search for javaScript code in <script> tag")
	// tui := flag.Bool("tui", false, "Terminal User Interface mode. Browse code in a full screen UI")
	// flag.Parse()

	cli := quicli.Parse(quicli.Cli{
		Usage:       "jse [flags] [HTML_STDIN]",
		Description: "Fastly gather all JavaScript from HTML",
		Flags: quicli.Flags{
			{Name: "gather-src", Description: "gather javascript code from script tag with src attribute. You must set url if you enabled it (-u flag)"},
			{Name: "url", Default: "", Description: "url of the html page HTML page (e.g. https://example.net/home.html)"},
			{Name: "tui", Description: "Terminal User Interface mode. Browse code in a full screen UI"},
			{Name: "ds", Description: "do not search for javaScript code in src attributes with <script> tag", NoShortName: true},
			{Name: "de", Description: "do not search for javaScript code in event attributes", NoShortName: true},
			{Name: "dt", Description: "do not search for javaScript code in <script> tag", NoShortName: true},
		},
	})

	// Translate cli -> config
	//transform url = https://[domain]/path/to/file -> https://[domain]
	var domain string
	url := cli.GetStringFlag("url")
	gatherSrc := cli.GetBoolFlag("gather-src")
	if url != "" {
		domain = strings.Join(strings.SplitAfter(url, "/")[:3], "")
		if domain[len(domain)-1:] == "/" {
			domain = domain[:len(domain)-1]
		}
	}
	cfg := config.Config{Url: domain, GatherSrc: gatherSrc, SkipSrc: cli.GetBoolFlag("ds"), SkipEvent: cli.GetBoolFlag("de"), SkipTag: cli.GetBoolFlag("dt")}

	if gatherSrc && (url == "") {
		log.Fatal("You must set domain if you enabled gathering js code from src (-gather-src) (-u flag)")
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
	if !cli.GetBoolFlag("tui") {
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
