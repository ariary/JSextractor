package ui

import (
	"JSextractor/pkg/config"
	"JSextractor/pkg/extract"
	"JSextractor/pkg/utils"
	"bytes"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/awesome-gocui/gocui"
)

var Index, MaxIndex int
var First bool

var Cfg *config.Config

var Scripts []extract.Script

var helpWindowToggle = false

const helpMessage = `
jse - Help
----------------------------------------------
ArrowDown		═ 	Move a line down
ArrowUp		═ 	Move a line up
ctrl + u		═ 	Change url
ctrl + g		═ 	Gather script from source attribute
ctrl + a		= 	Gather script for all source attributes
ctrl + c		═ 	Exit
ctrl + h		═ 	Toggle help message
`

var helpUrlWindowToggle = false

const helpUrlMessage = `
jse - Help
----------------------------------------------
Tab	═ Change request method (GET or cURL command line)
ctrl + z	═ Don't perform the requets, go back to scripts
Enter	═ Perform the request and parse it to gather scripts
ctrl + c	═ Exit
ctrl + h	═ Toggle help message
`

var Method, contentOtherMethod string

var cURL bool

const (
	scriptView  = "scripts"
	contentView = "content"
	urlView     = "url"
	methodView  = "method"
	helpView    = "help"
	helpUrlView = "helpUrl"
)

type position struct {
	prc    float32
	margin int
}

func (p position) getCoordinate(max int) int {
	// value = prc * MAX + abs
	return int(p.prc*float32(max)) - p.margin
}

type viewPosition struct {
	x0, y0, x1, y1 position
}

func (vp viewPosition) getCoordinates(maxX, maxY int) (int, int, int, int) {
	var x0 = vp.x0.getCoordinate(maxX)
	var y0 = vp.y0.getCoordinate(maxY)
	var x1 = vp.x1.getCoordinate(maxX)
	var y1 = vp.y1.getCoordinate(maxY)
	return x0, y0, x1, y1
}

var viewPositions = map[string]viewPosition{
	scriptView: {
		position{0.0, 0},
		position{0.0, 0},
		position{0.2, 2},
		position{0.9, 2},
	},
	contentView: {
		position{0.2, 0},
		position{0.0, 0},
		position{1.0, 2},
		position{0.9, 2},
	},
	urlView: {
		position{0.0, 0},
		position{0.89, 0},
		position{1.0, 2},
		position{1.0, 4},
	},
	methodView: {
		position{0.95, 0},
		position{0.89, 0},
		position{1.0, 2},
		position{1.0, 4},
	},
}

func UpdateUiVars() {
	Index = 0
	MaxIndex = len(Scripts) - 1
	First = true
}

//SetUrlView change the current view for the url one
func SetUrlView(g *gocui.Gui, v *gocui.View) error {
	if v.Name() != urlView {
		uv, err := g.SetCurrentView(urlView)
		if v.Name() == scriptView {
			//v.Highlight = false //disable highlight in script view
			v.SelBgColor = gocui.ColorCyan
		}
		uv.Highlight = true
		uv.FrameColor = gocui.ColorRed
		return err
	}
	return nil
}

//SetScriptView change the current view for the url one
func SetScriptView(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == urlView {
		sv, err := g.SetCurrentView(scriptView)
		v.Highlight = false //disable highlight in url view
		v.FrameColor = gocui.ColorWhite
		sv.SelBgColor = gocui.ColorGreen
		return err
	}
	return nil
}
func cursorMovement(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		nIndex := Index + d
		if nIndex < 0 || nIndex > MaxIndex {
			return nil
		}
		Index = nIndex
		if v != nil {
			cx, cy := v.Cursor()
			if err := v.SetCursor(cx, cy+d); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+d); err != nil {
					return err
				}
			}
		}
		dv, err := g.View(contentView)
		if err != nil {
			log.Fatal("failed to get contentView", err)
		}

		dv.Clear()
		DrawContentView(g, dv)
		return nil
	}
}

//Fetch new url, extract Script structures from result and present them in TUI = Change current view
func Fetch(g *gocui.Gui, v *gocui.View) (err error) {
	//JSE restart for a new input
	var body string
	Cfg.Url, err = v.Line(0)
	if err != nil {
		return err
	}

	if cURL {
		//cURL cmd
		body, err = utils.Curl(Cfg.Url)
	} else {
		//GET
		//fetch url
		body, err = utils.Fetch(Cfg.Url)
		if err != nil {
			return err
		}
	}
	// Extract scripts
	begins := utils.GetBeginLinesIndex([]byte(body))
	Scripts = extract.Extract(Cfg, *bytes.NewBuffer([]byte(body)), begins)

	//update ui var first, index etc
	UpdateUiVars()
	First = true //if fetch success

	//Update views
	v.Highlight = false
	v.FrameColor = gocui.ColorWhite

	sv, err := g.View(scriptView)
	sv.Clear()
	DrawScriptView(g, sv)
	cx, _ := sv.Cursor()
	sv.SetCursor(cx, 0) //put cursor at the top
	sv.Highlight = true
	sv.SelBgColor = gocui.ColorGreen

	cv, err := g.View(contentView)
	cv.Clear()
	DrawContentView(g, cv)

	return err
}

//GatherSrc gather javascript from src attr
func GatherSrc(g *gocui.Gui, v *gocui.View) (err error) {
	s := Scripts[Index]
	var code, domain string
	if s.Source == extract.FromSrc {
		if Cfg.Url != "" {
			domain = strings.Join(strings.SplitAfter(Cfg.Url, "/")[:3], "")
		}
		path := s.Content
		if path != "" {
			code, err = extract.GatherJS(path, domain)
			if err != nil {
				s.Content = s.Content + " (failed to retrieve code by fetching src)"
			} else {
				s.Content = code
				s.Source = extract.FromSrcGathered
			}
		}
	}
	Scripts[Index] = s

	//update script view
	cs, err := g.View(scriptView)
	cs.Clear()
	DrawScriptView(g, cs)

	//update content view
	cv, err := g.View(contentView)
	cv.Clear()
	fmt.Fprintln(cv, s.Content)

	return err
}

//GatherAll gather javascript from all script with src attr
func GatherAll(g *gocui.Gui, v *gocui.View) (err error) {
	for i := 0; i < len(Scripts); i++ {
		var code, domain string
		if Scripts[i].Source == extract.FromSrc {
			if Cfg.Url != "" {
				domain = strings.Join(strings.SplitAfter(Cfg.Url, "/")[:3], "")
			}
			path := Scripts[i].Content
			if path != "" {
				code, err = extract.GatherJS(path, domain)
				if err != nil {
					Scripts[i].Content += " (failed to retrieve code by fetching src)"
				} else {
					Scripts[i].Content = code
					Scripts[i].Source = extract.FromSrcGathered
				}
			}
		}
	}

	//update script view
	cs, err := g.View(scriptView)
	cs.Clear()
	DrawScriptView(g, cs)

	//update content view if current script is from src
	if Scripts[Index].Source == extract.FromSrc || Scripts[Index].Source == extract.FromSrcGathered {
		cv, err := g.View(contentView)
		if err != nil {
			return err
		}
		cv.Clear()
		fmt.Fprintln(cv, Scripts[Index].Content)
	}
	return err
}

//cursorDown: select element from the line below
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	return cursorMovement(1)(g, v)
}

//cursorUp: select element from the line above
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	return cursorMovement(-1)(g, v)
}

//ChangeMethod change the request method
func ChangeMethod(g *gocui.Gui, v *gocui.View) error {
	cURL = !cURL
	contentTmP, _ := v.Line(0) //save current content
	v.Clear()
	v.Write([]byte(contentOtherMethod)) //write other content
	v.EditGotoToEndOfLine()
	contentOtherMethod = contentTmP //edit other content
	return DrawMethodView(g, v)
}

//quit: quit the app (TUI)
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//toogleHelp show help message
func toggleHelp(g *gocui.Gui, v *gocui.View) error {
	helpWindowToggle = !helpWindowToggle
	return nil
}

//toogleHelp show help message
func toggleHelpUrl(g *gocui.Gui, v *gocui.View) error {
	helpUrlWindowToggle = !helpUrlWindowToggle
	return nil
}

//Keybindings define the key bindings of the TUI
func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(scriptView, gocui.KeyCtrlH, gocui.ModNone, toggleHelp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, SetUrlView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(scriptView, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(scriptView, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(scriptView, gocui.KeyCtrlG, gocui.ModNone, GatherSrc); err != nil {
		return err
	}
	if err := g.SetKeybinding(scriptView, gocui.KeyCtrlA, gocui.ModNone, GatherAll); err != nil {
		return err
	}
	if err := g.SetKeybinding(urlView, gocui.KeyEnter, gocui.ModNone, Fetch); err != nil {
		return err
	}
	if err := g.SetKeybinding(urlView, gocui.KeyCtrlZ, gocui.ModNone, SetScriptView); err != nil {
		return err
	}
	if err := g.SetKeybinding(urlView, gocui.KeyTab, gocui.ModNone, ChangeMethod); err != nil {
		return err
	}
	if err := g.SetKeybinding(urlView, gocui.KeyCtrlH, gocui.ModNone, toggleHelpUrl); err != nil {
		log.Panicln(err)
	}
	return nil
}

//DrawScriptView draw the view representing the list of script
func DrawScriptView(g *gocui.Gui, v *gocui.View) {
	for i := 0; i < len(Scripts); i++ {
		fmt.Fprintln(v, extract.ScriptInfoOutput(Scripts[i]))
	}
}

//DrawContentView draw the view representing the content of a script (js code)
func DrawContentView(g *gocui.Gui, v *gocui.View) {
	fmt.Fprintln(v, Scripts[Index].Content)
}

func DrawUrlView(g *gocui.Gui, v *gocui.View) error {
	uv, err := g.View(urlView)
	if err != nil {
		log.Fatal("failed to get pathView", err)
	}
	uv.Clear()
	//fmt.Fprintf(pv, "\t🌐 %s", url)
	fmt.Fprint(uv, Cfg.Url)
	return nil
}

func DrawMethodView(g *gocui.Gui, v *gocui.View) error {
	mv, err := g.View(methodView)
	if err != nil {
		log.Fatal("failed to get methodView", err)
	}
	mv.Clear()
	if cURL {
		Method = "cURL"
		mv.BgColor = gocui.ColorBlue
		mv.FgColor = gocui.ColorMagenta
	} else {
		Method = "GET"
		mv.BgColor = gocui.ColorGreen
		mv.FgColor = gocui.ColorBlack
	}

	fmt.Fprint(mv, utils.Bold(Method))
	return nil
}

//Layout organize the different views
func Layout(g *gocui.Gui) error {
	var views = []string{scriptView, contentView, urlView, methodView}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := viewPositions[view].getCoordinates(maxX, maxY)
		if v, err := g.SetView(view, x0, y0, x1, y1, 0); err != nil {
			v.SelFgColor = gocui.ColorBlack
			v.SelBgColor = gocui.ColorGreen
			if v.Name() != urlView {
				v.Title = " " + view + " "
			}
			if err != gocui.ErrUnknownView {
				return err
			}
			if v.Name() == scriptView {
				v.Highlight = true
				v.SelBgColor = gocui.ColorGreen
				v.SelFgColor = gocui.ColorBlack
				DrawScriptView(g, v)
			}
			if v.Name() == contentView {
				v.Editable = true
				v.Wrap = true
				DrawContentView(g, v)
			}
			if v.Name() == urlView {
				v.Title = " url 🌐"
				v.Editable = true
				v.Wrap = true
				v.KeybindOnEdit = true
				DrawUrlView(g, v)
				//v.SetCursor(1, 0)
			}
			if v.Name() == methodView {
				v.Frame = false
				v.Wrap = true
				DrawMethodView(g, v)
			}
		}
	}

	if First {
		if _, err := g.SetCurrentView(scriptView); err != nil {
			return err
		}
		First = false
	}

	if helpWindowToggle {
		height := strings.Count(helpMessage, "\n") + 1
		width := -1
		for _, line := range strings.Split(helpMessage, "\n") {
			width = int(math.Max(float64(width), float64(len(line)+2)))
		}
		if v, err := g.SetView(helpView, maxX/2-width/2, maxY/2-height/2, maxX/2+width/2, maxY/2+height/2, 0); err != nil {
			if err != gocui.ErrUnknownView {
				return err

			}
			fmt.Fprintln(v, helpMessage)
		}
	} else {
		g.DeleteView(helpView)
	}

	if helpUrlWindowToggle {
		height := strings.Count(helpUrlMessage, "\n") + 1
		width := -1
		for _, line := range strings.Split(helpUrlMessage, "\n") {
			width = int(math.Max(float64(width), float64(len(line)+2)))
		}
		if v, err := g.SetView(helpUrlView, maxX/2-width/2, maxY/2-height/2, maxX/2+width/2, maxY/2+height/2, 0); err != nil {
			if err != gocui.ErrUnknownView {
				return err

			}
			fmt.Fprintln(v, helpUrlMessage)
		}
	} else {
		g.DeleteView(helpUrlView)
	}

	return nil
}
