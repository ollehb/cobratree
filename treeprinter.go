package cobratree

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

var defaultFormat = FormatConfig{
	WithChild:    '├',
	WithOutChild: '│',
	LastChild:    '└',
	Horizontal:   '─',
	Padding:      3,
}

type formatDirective struct {
	depth       int
	crumbs      string
	isLast      bool
	hasSiblings bool
	format      *printConfig
}

// Character output scheme
type FormatConfig struct {
	// Symbol to draw for connecting children
	WithChild    rune
	// Connection for a level without a direct child
	WithOutChild rune
	// Symbol for last child to a command
	LastChild    rune
	// Horizontal connector
	Horizontal   rune
	// Amount of padding chars
	Padding      int
}

// Write command tree with default formatting.
func WriteTree(writer io.Writer, command *cobra.Command) error {
	return WriteCustom(writer, command, defaultFormat)
}

// Write command tree with a custom format.
func WriteCustom(writer io.Writer, command *cobra.Command, config FormatConfig) error {
	tree := ParseCommandTree(command)
	return WriteCustomParsed(writer, tree, config)
}

// Write a custom parsed command tree.
// Too much info might produce ugly output.
func WriteCustomParsed(writer io.Writer, info CommandInfo, config FormatConfig) error {
	printConf := newPrintConfig(config)
	return printRecursive(writer, info, formatDirective{format: printConf})
}

func printRecursive(writer io.Writer, info CommandInfo, directive formatDirective) error {
	currentCrumbs := directive.crumbs
	switch {
	case directive.hasSiblings && !directive.isLast:
		currentCrumbs += directive.format.withChild
	case directive.isLast:
		currentCrumbs += directive.format.lastChild
	}
	newCrumbs := ""
	if directive.hasSiblings {
		newCrumbs = directive.format.notCommandsLevel
	} else if directive.depth > 0 {
		newCrumbs = directive.format.padding
	}

	if _, err := fmt.Fprintf(writer, "%s%s\n", currentCrumbs, info.Val); err != nil {
		return err
	}
	for i := 0; i < len(info.Children); i++ {
		f := formatDirective{
			depth:       directive.depth + 1,
			crumbs:      directive.crumbs + newCrumbs,
			isLast:      i == len(info.Children)-1,
			format:      directive.format,
			hasSiblings: len(info.Children) > 1,
		}
		if err := printRecursive(writer, info.Children[i], f); err != nil {
			return err
		}
	}
	return nil
}

type printConfig struct {
	withChild        string
	withoutChild     string
	lastChild        string
	notCommandsLevel string
	padding          string
}

func newPrintConfig(config FormatConfig) *printConfig {
	printConf := new(printConfig)
	padding := strings.Repeat(" ", config.Padding)
	printConf.withChild = string(config.WithChild) + string(config.Horizontal)
	printConf.withoutChild = string(config.WithOutChild)
	printConf.lastChild = string(config.LastChild) + string(config.Horizontal)
	printConf.padding = padding
	printConf.notCommandsLevel = string(config.WithOutChild) + padding
	return printConf
}
