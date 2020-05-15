package commands

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

// // CommandLine contains all the information passed to the commands on the command line.
// type CommandLine interface {
// 	ShowHelp()

// 	ShowVersion()

// 	Application() *cli.App

// 	Args() cli.Args

// 	IsSet(name string) bool

// 	Bool(name string) bool

// 	Int(name string) int

// 	String(name string) string

// 	StringSlice(name string) []string

// 	GlobalString(name string) string

// 	FlagNames() (names []string)

// 	Generic(name string) interface{}
// }

var Commands = []cli.Command{
	{
		Flags:       GetImagesListFlags,
		Name:        "getimg",
		Usage:       "get huaweicloud images list",
		Description: fmt.Sprintf("Run '%s get --help' to include the get flags for list images in the help text.", os.Args[0]),
		Action:      cmdGetImageListAction,
	},
}
