package tripanel

import (
	"fmt"
	"github.com/BellerophonMobile/commandtree"	
	"github.com/jroimartin/gocui"
)

var Commands *commandtree.CommandTree

func installbuiltincommands(commands *commandtree.CommandTree) error {

	Commands = commands

	err := Commands.Add(&commandtree.Command{
		Command: "clear",
		Description: "Clear application data view",
		Usage: "May stipulate to clear app and/or log views (space separated); defaults to app",
		Action: func(args []string) error {
			if len(args) <= 0 {
				clear(GUI, Views.App)
			} else {
				for _,a := range(args) {
					switch a {
					case "app":
						clear(GUI, Views.App)						
					case "log":
						clear(GUI, Views.Log)						
					default:
						return fmt.Errorf("No such view '%v'", a)
					}
				}
			}

			return nil
		},
	})
	if err != nil { return err }
	
	err = Commands.Add(&commandtree.Command{
		Command: "help",
		Description: "Display usage notes",
		Usage: "May be given command tree for details, or none for top level help",
		Action: func(args []string) error {
			if len(args) <= 0 {
				DisplayUsage()
			} else {
				usage,err := Commands.Help(args)
				if err != nil {
					if _,ok := err.(commandtree.NoSuchCommandError); ok {
						fmt.Fprintf(Views.App, "\033[31mNo such command %v, run 'help' for a list\033[0m\n", args)
						return nil
					} else {
						return err
					}
				}
				fmt.Fprintf(Views.App, "\033[34;7m USAGE \033[34m %v\033[0m\n", usage)
			}
			return nil
		},
	})
	if err != nil { return err }

	err = Commands.Add(&commandtree.Command{
		Command: "quit",
		Description: "Quit the application",
		Action: func([]string) error {
			return gocui.ErrQuit
		},
	})
	if err != nil { return err }
	
	return nil

}
