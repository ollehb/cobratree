package cobratree

import (
	"github.com/spf13/cobra"
)

type CommandInfo struct {
	Val      string
	Children []CommandInfo
}

type InfoFunc func(*cobra.Command) string

func NameFunc(command *cobra.Command) string {
	return command.Name()
}

// Traverse the command tree collecting command names as reported by cobra.
func ParseCommandTree(command *cobra.Command) CommandInfo {
	return ParseCustom(command, NameFunc)
}

// Traverse the command tree collecting info using the supplied info func.
func ParseCustom(command *cobra.Command, infoFunc InfoFunc) CommandInfo {
	return visitCommand(command, infoFunc)
}

func visitCommand(cmd *cobra.Command, infoFunc InfoFunc) CommandInfo {
	children := cmd.Commands()
	comm := CommandInfo{
		Val: infoFunc(cmd),
	}
	if children != nil {
		childCommands := make([]CommandInfo, 0, len(children))
		for _, child := range children {
			childCommands = append(childCommands, visitCommand(child, infoFunc))
		}
		comm.Children = childCommands
	}
	return comm
}
