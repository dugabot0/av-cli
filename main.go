package main

import (
	"github.com/yourusername/av-cli/cmd"
	"github.com/yourusername/av-cli/cmd/dmm"
	"github.com/yourusername/av-cli/cmd/duga"
)

func main() {
	cmd.AddCommand(dmm.NewCmd())
	cmd.AddCommand(duga.NewCmd())
	cmd.Execute()
}
