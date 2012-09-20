package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/sbinet/liner"
	_ "github.com/sbinet/go-croot/pkg/croot"
	"github.com/sbinet/paw-go/pkg/gribble"
	"github.com/sbinet/paw-go/pkg/pawmgr"
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

// A Gribble environment is composed of Go struct types. Since the environment
// does not care about values, zero values of each command struct can be
// used in the gribble.Command slice.
//
// Note that these may also be specified as non-pointers.
var env *gribble.Environment = gribble.New(
	[]gribble.Command{
		new_cmd_ntuple_create(),
	},
)

var paw_completer liner.Completer = func(line string) []string {
	completions := []string{}
	for _,cmd := range env.CommandNames() {
		if strings.HasPrefix(cmd, line) {
			completions = append(completions, cmd)
		}
	}
	return completions
}

func new_cmd_ntuple_create() *ntuple_create {
	return &ntuple_create{
		//mgr: nil,
	}
}

type ntuple_create struct {
	name     string      `/ntuple/create`
	Id int `param:"1" types:"int"`
	Title string `param:"2" types:"string"`
	Vars  string `param:"3" types:"string"`

	mgr *pawmgr.NtupleMgr
}

func (cmd *ntuple_create) Run() gribble.Value {
	id := cmd.Id
	title,err := strconv.Unquote(cmd.Title)
	if err != nil {
		return err
	}
	//nvars := cmd.Nvars.(int)
	vars, err := strconv.Unquote(cmd.Vars)
	ntvars := strings.Split(vars, ",")
	return cmd.mgr.Create(id, title, ntvars)
}

func new_cmd_ntuple_plot() *ntuple_plot {
	return &ntuple_plot{
		//mgr: nil,
	}
}

type ntuple_plot struct {
	name     string      `/ntuple/plot`
	Id int `param:"1" types:"int"`
	//Func string `param:"2" types:"string"`
	//Vars  string `param:"3" types:"string"`

	mgr *pawmgr.NtupleMgr
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
