package tripanel

import (
	"fmt"
	"github.com/BellerophonMobile/logberry"
	"github.com/BellerophonMobile/commandtree"	
	"github.com/jroimartin/gocui"
)

var GUI *gocui.Gui

func Run(commands *commandtree.CommandTree) error {
		
	defer logberry.Std.Stop()

	err := installbuiltincommands(commands)
	if err != nil { return err }	

	GUI, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil { return err }
	defer GUI.Close()
	
	GUI.Cursor = true

	GUI.SetManagerFunc(layout)

	err = keybindings()
	if err != nil { return err }

	err = GUI.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
	
}

func DisplayUsage() {
	fmt.Fprintln(Views.App, "\033[34;7m USAGE \033[34m Application overview---\n")	
	fmt.Fprintln(Views.App, "Keyboard actions:")
	fmt.Fprintln(Views.App, KeyboardUsage())
	fmt.Fprintln(Views.App, "Command prompt actions:")
	fmt.Fprintln(Views.App, Commands.Usage())
	fmt.Fprintln(Views.App, "Command prompt may have doublequoted entries, with '\\\"', '\\\\', and '\\n' escapes")
	fmt.Fprintln(Views.App, "Command prompt may have comments, beginning with a '#' outside doublequotes")
	fmt.Fprintln(Views.App, "Help command can provide details for other commands given as parameters")
	fmt.Fprint(Views.App, "\033[0m")
}
