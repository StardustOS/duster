package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AtomicMalloc/debugger/debugger"
	"github.com/c-bata/go-prompt"
)

type Command interface {
	Text() string
	Description() string
	Run() (string, error)
	Matches(string) bool
}

type CLI struct {
	prompt      string
	msg         string
	commands    []Command
	dbg         debugger.Debugger
	suggestions []prompt.Suggest
}

func (cli *CLI) Init(domaind uint32, filename string) error {
	cli.prompt = ">"
	cli.suggestions = []prompt.Suggest{
		prompt.Suggest{Text: "break", Description: "Sets a break point at in a file (argument in the form of file.c:<line no>"},
		prompt.Suggest{Text: "step", Description: "Steps forward one line (note a breakpoint must be set before hand)"},
	}
	cli.dbg = debugger.Debugger{}
	err := cli.dbg.Init(domaind, filename)
	return err
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
			fmt.Printf("Break point set @ %s:%d", values[0], lineNo)
		} else {
			fmt.Println(err)
		}
	case "step":
		if cli.dbg.IsPaused() {
			cli.dbg.StartSingle(0, true)
			cli.dbg.Step(0)
			fmt.Println(cli.dbg.GetLineInformation())
		} else {
			fmt.Println("Error: the doamin has not been paused")
		}
	case "quit":
		os.Exit(0)
	}
}
