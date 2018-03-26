package language

import (
	"fmt"

	"github.com/greaka/gw2godiscord/commands"
)

func NotACommand(lang Language) string {
	switch lang {
	case English:
		return fmt.Sprintf("This is not a valid command. Type %s for a list of commands.", commands.CommandHelp)
	case German:
		return fmt.Sprintf("Das ist kein gültiger Befehl. Für eine Liste von Befehlen, siehe %s", commands.CommandHelp)
	default:
		return fmt.Sprintf("This is not a valid command. Type %s for a list of commands.", commands.CommandHelp)
	}
}
