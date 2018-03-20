package tripanel

import (
	"fmt"
	"github.com/BellerophonMobile/logberry"
	"io/ioutil"
	"strings"
	//	"github.com/BellerophonMobile/commandtree"
	"github.com/BellerophonMobile/logberry/gocuioutput"

	"github.com/BellerophonMobile/gocui"
)

var maximized = "app"

const minimizedsize = 7

var commandhistory = make([]string, 0, 128)
var commandpartial string
var commandindex int

var Views = struct {
	Cmd *gocui.View
	App *gocui.View
	Log *gocui.View
}{}

func cmdexecute(g *gocui.Gui, v *gocui.View) error {

	var err error

	_, cy := v.Cursor()
	line, _ := v.Line(cy)
	line = strings.TrimSpace(line)

	if line == "" {
		return nil
	}

	if len(commandhistory) <= 0 || commandhistory[len(commandhistory)-1] != line {
		commandhistory = append(commandhistory, line)
		commandindex = len(commandhistory)
	}
	commandpartial = ""

	cerr := Commands.Execute(line)
	if cerr != nil {
		if cerr == gocui.ErrQuit {
			return quit(g, v)
			//		} else if _,ok := cerr.(commandtree.NoSuchCommandError); ok {
		} else {
			fmt.Fprintf(Views.App, "\033[31;7m ERROR \033[31m %v\033[0m\n", cerr)
		}
	}

	fmt.Fprintln(Views.App)

	v.Clear()
	v.SetCursor(0, 0)

	return err

}

func cmdup(g *gocui.Gui, v *gocui.View) error {

	if len(commandhistory) <= 0 {
		return nil
	}

	if commandindex >= len(commandhistory) {
		_, cy := v.Cursor()
		commandpartial, _ = v.Line(cy)
		commandindex = len(commandhistory) - 1
	} else {
		commandindex--
		if commandindex < 0 {
			commandindex = 0
		}
	}

	v.Clear()
	fmt.Fprint(v, commandhistory[commandindex])
	v.SetCursor(len(commandhistory[commandindex]), 0)

	return nil

}

func cmddown(g *gocui.Gui, v *gocui.View) error {

	if commandindex == len(commandhistory) {
		return nil
	}

	var c string
	commandindex++
	if commandindex >= len(commandhistory) {
		c = commandpartial
		commandpartial = ""
		commandindex = len(commandhistory)
	} else {
		c = commandhistory[commandindex]
	}

	v.Clear()
	fmt.Fprint(v, c)
	v.SetCursor(len(c), 0)

	return nil
}

func movetoend(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	v.SetCursor(len(l), cy)
	return nil
}

func movetostart(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	v.SetCursor(0, cy)
	return nil
}

func switchview(label string, g *gocui.Gui, v *gocui.View) error {

	Views.App.Title = "   App   "
	Views.Log.Title = "   Log   "
	Views.Cmd.Title = "   Cmd   "

	switch label {
	case "cmd":
		g.Cursor = true
		Views.Cmd.Title = " [ Cmd ] "
	case "app":
		g.Cursor = false
		Views.App.Title = " [ App ] "
	case "log":
		g.Cursor = false
		Views.Log.Title = " [ Log ] "
	default:
		return fmt.Errorf("Unmanaged view '%v'", label)
	}

	_, err := g.SetCurrentView(label)

	return err

}

func prevview(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "cmd":
		return switchview("log", g, v)
	case "app":
		return switchview("cmd", g, v)
	case "log":
		return switchview("app", g, v)
	}
	return fmt.Errorf("Unknown current view '%v'", v.Name())
}

func nextview(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "cmd":
		return switchview("app", g, v)
	case "app":
		return switchview("log", g, v)
	case "log":
		return switchview("cmd", g, v)
	}
	return fmt.Errorf("Unknown current view '%v'", v.Name())
}

func scrolldown(g *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = false

	ox, oy := v.Origin()
	oy++

	_, sy := v.Size()
	my := (v.NumLines() - sy)
	if oy > my {
		oy = my
	}
	if oy < 0 {
		oy = 0
	}

	if err := v.SetOrigin(ox, oy); err != nil {
		return err
	}
	return nil

}

func scrollup(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	oy--
	if oy <= 0 {
		oy = 0
	}
	v.Autoscroll = false
	if err := v.SetOrigin(ox, oy); err != nil {
		return err
	}
	return nil
}

func pagedown(g *gocui.Gui, v *gocui.View) error {

	v.Autoscroll = false

	_, sy := v.Size()

	ox, oy := v.Origin()
	oy += sy

	my := (v.NumLines() - sy)
	if oy > my {
		oy = my
	}
	if oy < 0 {
		oy = 0
	}

	if err := v.SetOrigin(ox, oy); err != nil {
		return err
	}
	return nil

}

func pageup(g *gocui.Gui, v *gocui.View) error {

	v.Autoscroll = false

	_, sy := v.Size()

	ox, oy := v.Origin()
	oy -= sy

	if oy <= 0 {
		oy = 0
	}

	if err := v.SetOrigin(ox, oy); err != nil {
		return err
	}

	return nil
}

func home(g *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = false
	if err := v.SetOrigin(0, 0); err != nil {
		return err
	}
	return nil
}

func end(g *gocui.Gui, v *gocui.View) error {

	v.Autoscroll = true

	_, sy := v.Size()
	oy := v.NumLines() - sy
	if oy < 0 {
		oy = 0
	}

	if err := v.SetOrigin(0, oy); err != nil {
		return err
	}
	return nil
}

func cursordown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorup(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

var copybuffer string

func kill(g *gocui.Gui, v *gocui.View) error {

	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	if l != "" {
		copybuffer = l
	}

	v.Clear()
	v.SetCursor(0, 0)
	if err := v.SetOrigin(0, 0); err != nil {
		return err
	}

	return nil
}

func paste(g *gocui.Gui, v *gocui.View) error {
	fmt.Fprint(v, copybuffer)
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	v.SetCursor(len(l), cy)
	return nil
}

func clear(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	v.SetCursor(0, 0)
	if err := v.SetOrigin(0, 0); err != nil {
		return err
	}
	return nil
}

func togglemaximized(g *gocui.Gui, v *gocui.View) error {
	if maximized == "log" {
		maximized = "app"
	} else {
		maximized = "log"
	}
	return nil
}

func writeview(g *gocui.Gui, v *gocui.View) error {

	go func() {

		fn, err := TextPrompt("File?")
		if err != nil {
			UError(err)
			return
		}

		err = ioutil.WriteFile(fn, []byte(v.Buffer()), 0644)
		if err != nil {
			UError(err)
			return
		}

		UPrint(fmt.Sprintf("Wrote %v to %v\n", v.Name(), fn))

	}()

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	doquits()
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {

	var err error

	maxX, maxY := g.Size()

	if Views.Cmd, err = g.SetView("cmd", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		Views.Cmd.Highlight = true
		Views.Cmd.Editable = true
		Views.Cmd.Title = " [ Cmd ] "

		if _, err := g.SetCurrentView("cmd"); err != nil {
			return err
		}

	}

	var y int
	if maximized == "log" {
		y = 3 + 2 + minimizedsize
	} else {
		y = maxY - (minimizedsize + 2 + 1)
	}

	if Views.App, err = g.SetView("app", 0, 3, maxX-1, y); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		Views.App.Title = "   App   "
		Views.App.Wrap = true
		Views.App.Autoscroll = true

		DisplayUsage()
		fmt.Fprintln(Views.App)

	}

	if maximized == "log" {
		y = 3 + 2 + minimizedsize + 1
	} else {
		y = maxY - (minimizedsize + 2)
	}

	if Views.Log, err = g.SetView("log", 0, y, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		Views.Log.Title = "   Log   "
		Views.Log.Wrap = true
		Views.Log.Autoscroll = true

		tout := logberry.NewTextOutput(Views.Log, "client")
		tout.Color = true
		logberry.Std.SetOutputDriver(gocuioutput.New(g, tout))
		logberry.Main.Ready()

	}

	return nil

}

func UPrint(msg string) {
	GUI.Update(func(g *gocui.Gui) error {
		fmt.Fprintf(Views.App, "%v\n", msg)
		return nil
	})
}

func UFailure(msg string) {
	GUI.Update(func(g *gocui.Gui) error {
		fmt.Fprintf(Views.App, "\033[31;7m ERROR \033[31m %v\033[0m\n", msg)
		return nil
	})
}

func UError(err error) {
	GUI.Update(func(g *gocui.Gui) error {
		fmt.Fprintf(Views.App, "\033[31;7m ERROR \033[31m %v\033[0m\n", err)
		return nil
	})
}
