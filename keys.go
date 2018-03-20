package tripanel

import (
	"fmt"
	"bytes"
	"github.com/jroimartin/gocui"
)

var KeyUsages []*KeyUsage // Slice rather than map so they're ordered

type KeyUsage struct {
	Key string
	Views []string
	Usage string
}

func AddKeyUsage(viewname string, key interface{}, mod gocui.Modifier, usage string) {

	var found bool
	descrip := gocui.DescribeKey(key, mod)
	
	for _,k := range(KeyUsages) {
		if k.Key == descrip && k.Usage == usage {
			k.Views = append(k.Views, viewname)
			found = true
		}
	}

	if !found {
		KeyUsages = append(KeyUsages, &KeyUsage{
			Key: descrip,
			Views: []string{viewname},
			Usage: usage,
		})
	}
	
}

func KeyboardUsage() string {

	if len(KeyUsages) <= 0 {
		return "   <no keys installed>\n"
	}

	var buff = &bytes.Buffer{}	
	fmt.Fprintf(buff, "   %-12v   %-12v   %v\n", "View(s)", "Key(s)", "Command")
	fmt.Fprintf(buff, "   ------------   ------------   ------------------------------------------------\n")
	
	for _,k := range(KeyUsages) {
		var views string
		if len(k.Views) > 0 {
			views = k.Views[0]
			for _,v := range(k.Views[1:]) {
				views += ","+v
			}
		}
		
		fmt.Fprintf(buff, "   %-12v   %-12v   %v\n", views, k.Key, k.Usage)
	}

	return buff.String()

}

func keybindings() error {

	var err error
	setkey := func(viewname string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error, usage string) {
		if err != nil {
			return
		}
		err = GUI.SetKeybinding(viewname, key, mod, handler)
		if usage != "" {
			AddKeyUsage(viewname, key, mod, usage)
		}
	}
	
	// Global keys
	setkey("", gocui.KeyCtrlC, gocui.ModNone, quit, "Quit the application")

	for _,c := range([]string{"cmd", "app", "log"}) {
		setkey(c, gocui.KeyTab, gocui.ModNone, nextview, "Cycle to next view")
		setkey(c, gocui.KeyCtrlBackslash, gocui.ModNone, nextview, "Cycle to next view")
		setkey(c, gocui.KeyCtrlBackslash, gocui.ModAlt, prevview, "Cycle to previous view")
	}
		
	// Command prompt keys
	setkey("cmd", gocui.KeyEnter, gocui.ModNone, cmdexecute, "Execute command line")
	setkey("cmd", gocui.KeyCtrlV, gocui.ModNone, kill, "Cut command line")
	setkey("cmd", gocui.KeyCtrlY, gocui.ModNone, paste, "Paste command line")
	setkey("cmd", gocui.KeyCtrlA, gocui.ModNone, movetostart, "Move cursor to start of line")
	setkey("cmd", gocui.KeyCtrlE, gocui.ModNone, movetoend, "Move cursor to end of line")
	setkey("cmd", gocui.KeyArrowUp, gocui.ModNone, cmdup, "Previous command line history")
	setkey("cmd", gocui.KeyArrowDown, gocui.ModNone, cmddown, "Next command line history")
	
	// App and log data pane keys
	for _,k := range([]string{"app", "log"}) {
		setkey(k, gocui.KeyArrowUp, gocui.ModNone, scrollup, "Scroll up and disable autoscroll")
		setkey(k, gocui.KeyArrowDown, gocui.ModNone, scrolldown, "Scroll down and disable autoscroll")
		setkey(k, gocui.KeyPgup, gocui.ModNone, pageup, "Page up and disable autoscroll")
		setkey(k, gocui.KeyPgdn, gocui.ModNone, pagedown, "Page down and disable autoscroll")
		setkey(k, gocui.KeyHome, gocui.ModNone, home, "Scroll to beginning and disable autoscroll")
		setkey(k, gocui.KeyEnd, gocui.ModNone, end, "Scroll to end and restore autoscroll")
		setkey(k, 'w', gocui.ModNone, writeview, "Write view to file")		
		setkey(k, gocui.KeyDelete, gocui.ModNone, clear, "Clear view")
		setkey(k, gocui.KeySpace, gocui.ModNone, togglemaximized, "Maximize/minimize view")
	}

	setkey("textprompt", gocui.KeyEnter, gocui.ModNone, textreceived, "")
	setkey("selectionprompt", gocui.KeyEnter, gocui.ModNone, selectionreceived, "")
	setkey("selectionprompt", gocui.KeyArrowUp, gocui.ModNone, cursorup, "")
	setkey("selectionprompt", gocui.KeyArrowDown, gocui.ModNone, cursordown, "")

	return err

}
