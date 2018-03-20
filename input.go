package tripanel

import (
	"fmt"
	"github.com/BellerophonMobile/gocui"
)

var textpromptchan = make(chan string)

func TextPrompt(title string) (string, error) {

	var err error

	GUI.Update(func(g *gocui.Gui) error {

		maxX, maxY := GUI.Size()
		if v, err := GUI.SetView("textprompt", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Editable = true
			v.Title = " " + title + " "
			if _, err := GUI.SetCurrentView("textprompt"); err != nil {
				return err
			}
		}
		return err

	})

	line := <-textpromptchan

	return line, err

}

func textreceived(g *gocui.Gui, v *gocui.View) error {

	_, cy := v.Cursor()
	line, _ := v.Line(cy)
	textpromptchan <- line

	if err := g.DeleteView("textprompt"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("cmd"); err != nil {
		return err
	}

	return nil

}

func SelectionPrompt(title string, options []string) (string, error) {

	var err error

	GUI.Update(func(g *gocui.Gui) error {

		maxX, maxY := GUI.Size()
		if v, err := GUI.SetView("selectionprompt", maxX/2-30, maxY/2-(len(options)/2), maxX/2+30, maxY/2+2+(len(options)/2)); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Highlight = true
			v.SelBgColor = gocui.ColorCyan
			v.SelFgColor = gocui.ColorBlack

			v.Title = " " + title + " "

			for _, o := range options {
				fmt.Fprintln(v, o)
			}

			if _, err := GUI.SetCurrentView("selectionprompt"); err != nil {
				return err
			}
		}
		return err

	})

	line := <-textpromptchan

	return line, err

}

func selectionreceived(g *gocui.Gui, v *gocui.View) error {

	_, cy := v.Cursor()
	line, _ := v.Line(cy)
	textpromptchan <- line

	if err := g.DeleteView("selectionprompt"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("cmd"); err != nil {
		return err
	}

	return nil

}
