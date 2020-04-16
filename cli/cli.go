package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
)

type Debugger interface {
	Continue(uint32) error
	SetBreakpoint(string, int, uint32) error
	RemoveBreakpoint(string, int, uint32) error
	Step(uint32) error
	GetLineInformation() string
	GetVariable(string) (string, error)
	Dereference(uint32, string) (string, error)
}

type CLI struct {
	prompt      string
	dbg         Debugger
	suggestions []prompt.Suggest
}

func (cli *CLI) Init(debugger Debugger) {
	cli.prompt = ">"
	cli.suggestions = []prompt.Suggest{
		prompt.Suggest{Text: "break", Description: "Sets a break point at in a file (argument in the form of file.c:<line no>"},
		prompt.Suggest{Text: "step", Description: "Steps forward one line (note a breakpoint must be set before hand)"},
		prompt.Suggest{Text: "continue", Description: "Continue to the next breakpoint"},
		prompt.Suggest{Text: "quit", Description: "Exit the debugger"},
		prompt.Suggest{Text: "read", Description: "Read a variable"},
		prompt.Suggest{Text: "def", Description: "Deference a variable"},
		prompt.Suggest{Text: "remove", Description: "Remove breakpoint"},
	}
	cli.dbg = debugger
}

func (cli *CLI) completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(cli.suggestions, d.GetWordBeforeCursor(), true)
}

func (cli *CLI) ReadInput() string {
	input := prompt.Input(cli.prompt, cli.completer)
	return input
}

func (cli *CLI) ProcessInput(input string) {
	values := strings.Split(input, " ")
	switch values[0] {
	default:
		fmt.Printf("Error: %s is not a recognised command", values[0])
	case "break":
		if len(values) != 2 {
			fmt.Println("Error: too many arguments for break (expected argument in the form file.c:<line no>)")
			return
		}
		args := strings.Split(values[1], ":")
		if len(args) != 2 {
			fmt.Println("Error: argument in the wrong format (expected file.c:<line no>)")
		}
		lineNo, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
		}
		err = cli.dbg.SetBreakpoint(args[0], lineNo, 0)
		if err == nil {
			fmt.Printf("Break point set @ %s:%d\n", args[0], lineNo)
		} else {
			fmt.Println(err)
		}
	case "step":
		err := cli.dbg.Step(0)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cli.dbg.GetLineInformation())
	case "read":
		val, err := cli.dbg.GetVariable(values[1])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}
	
	case "def":
		val, err := cli.dbg.Dereference(0, values[1])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}
		
	case "remove":
		if len(values) != 2 {
			fmt.Println("Error: too many arguments for break (expected argument in the form file.c:<line no>)")
			return
		}
		args := strings.Split(values[1], ":")
		if len(args) != 2 {
			fmt.Println("Error: argument in the wrong format (expected file.c:<line no>)")
		}

		lineNo, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
		}

		err = cli.dbg.RemoveBreakpoint(args[0], lineNo, 0)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Removed breakpoint at %s:%d\n", args[0], lineNo)
		}

		
	case "quit":
		os.Exit(0)
	case "continue":
		err := cli.dbg.Continue(0)
		if err != nil {
			fmt.Println(err)
		}
	}
}
