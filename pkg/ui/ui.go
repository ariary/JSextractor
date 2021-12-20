package ui

import (
	"JSextractor/pkg/config"
	"JSextractor/pkg/extract"
	"JSextractor/pkg/utils"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var Index, MaxIndex int
var First bool

var Cfg *config.Config

var Scripts []extract.Script

const (
	scriptView  = "scripts"
	contentView = "content"
	urlView     = "url"
	helpView    = "help"
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

var helpWindowToggle = false

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
			v.Highlight = false //disable highlight in script view
		}
		uv.Highlight = true
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

//Fetch new url, extract Script structures from result and present them in TUI
func Fetch(g *gocui.Gui, v *gocui.View) (err error) {

	//JSE restart for a new input
	//fetch url
	Cfg.Url, err = v.Line(0)
	if err != nil {
		return err
	}
	body, err := utils.Fetch(Cfg.Url)
	if err != nil {
		return err
	}
	// Extract scripts
	begins := utils.GetBeginLinesIndex([]byte(body))
	Scripts = extract.Extract(Cfg, *bytes.NewBuffer([]byte(body)), begins)

	//update ui var first, index etc
	UpdateUiVars()
	First = true //if fetch success

	//Update views
	v.Highlight = false

	sv, err := g.View(scriptView)
	sv.Clear()
	DrawScriptView(g, sv)
	cx, _ := sv.Cursor()
	sv.SetCursor(cx, 0) //put cursor at the top
	sv.Highlight = true

	cv, err := g.View(contentView)
	cv.Clear()
	DrawContentView(g, cv)

	return err
}

//GatherSrc gather javascript from src attr
func GatherSrc(g *gocui.Gui, v *gocui.View) (err error) {
	s := Scripts[Index]
	var code string
	if s.Source == extract.FromSrc {
		domain := strings.Join(strings.SplitAfter(Cfg.Url, "/")[:3], "")
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
		var code string
		if Scripts[i].Source == extract.FromSrc {
			domain := strings.Join(strings.SplitAfter(Cfg.Url, "/")[:3], "")
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

//quit: quit the app (TUI)
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//Keybindings define the key bindings of the TUI
func Keybindings(g *gocui.Gui) error {

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

//Layout organize the different views
func Layout(g *gocui.Gui) error {
	var views = []string{scriptView, contentView, urlView}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := viewPositions[view].getCoordinates(maxX, maxY)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
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
				DrawUrlView(g, v)
				v.SetCursor(1, 0)
			}
		}
	}

	if First {
		if _, err := g.SetCurrentView(scriptView); err != nil {
			return err
		}
		First = false
	}

	return nil
}
