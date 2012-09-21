package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sbinet/liner"
	_ "github.com/sbinet/go-croot/pkg/croot"
)

var term *liner.State = nil

func init() {
	fmt.Println(`
************************
** PAW-Go interpreter **
************************
`)
	fmt.Printf(":: available commands:\n%v\n::\n\n", env)

	term = liner.NewLiner()

	fname := path.Join(os.Getenv("HOME"), ".pawgo.history")
	f, err := os.Open(fname)
	if err != nil {
		f, err = os.Create(fname)
		if err != nil {
			fmt.Printf("**warning: could not access nor create history file [%s]\n", fname)
			return
		}
	}
	defer f.Close()
	_, err = term.ReadHistory(f)
	if err != nil {
		fmt.Printf("**warning: could not read history file [%s]\n", fname)
		return
	}

	term.SetCompleter(paw_completer)
}

func atexit() {
	fname := path.Join(os.Getenv("HOME"), ".pawgo.history")
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("**warning: could not access history file [%s]\n", fname)
		return
	}
	defer f.Close()
	_, err = term.WriteHistory(f)
	if err != nil {
		fmt.Printf("**warning: could not write history file [%s]\n", fname)
		return
	}

	err = term.Close()
	if err != nil {
		fmt.Printf("**warning: problem closing term: %v\n", err)
		return
	}
}

var paw_completer liner.Completer = func(line string) []string {
	completions := []string{}
	for _,cmd := range env.CommandNames() {
		if strings.HasPrefix(cmd, line) {
			completions = append(completions, cmd)
		}
	}
	return completions
}

func main() {

	defer atexit()

	prompt := "paw> "
	cmd := ""

	for {
		line, err := term.Prompt(prompt)
		if err != nil {
			break //os.Exit(0)
		}
		cmd = line
		if cmd != "" {
			for _, ll := range strings.Split(cmd, "\n") {
				term.AppendHistory(ll)
			}
		} else {
			continue
		}

		_, err = env.Run(cmd)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Printf("(error %T)\n", err)
			continue
		}
	}

	fmt.Printf("\n:: bye.\n")
}

// EOF
