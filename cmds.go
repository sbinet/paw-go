package main

import (
	"strconv"
	"strings"
	
	"github.com/sbinet/paw-go/pkg/gribble"
	"github.com/sbinet/paw-go/pkg/pawmgr"
)

// A Gribble environment is composed of Go struct types. Since the environment
// does not care about values, zero values of each command struct can be
// used in the gribble.Command slice.
//
// Note that these may also be specified as non-pointers.
var env *gribble.Environment = gribble.New(
	[]gribble.Command{
		new_cmd_ntuple_create(),
		new_cmd_ntuple_plot(),
	},
)


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

func (cmd *ntuple_plot) Run() gribble.Value {
	id := cmd.Id
	return cmd.mgr.Plot(id)
}

