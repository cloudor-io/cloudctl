package utils

import (
	"os"

	"github.com/fatih/color"
)

func CheckErr(msg interface{}) {
	if msg != nil {
		color.New(color.FgRed).Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}
