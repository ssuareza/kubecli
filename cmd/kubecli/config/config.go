package config

import (
	"fmt"
	"os"

	"github.com/ansel1/merry"
	"github.com/fatih/color"
)

type C struct {
	Debug bool `mapstructure:"debug"`
}

type Styles struct {
	Title   func(...interface{}) string
	Error   func(...interface{}) string
	Prompt  func(...interface{}) string
	Success func(...interface{}) string
	Info    func(...interface{}) string
}

var Style Styles
var Config C
var ErrorHandler func(error)

func init() {
	Style.Title = color.New(color.FgWhite).Add(color.Bold).SprintFunc()
	Style.Error = color.New(color.FgRed).SprintFunc()
	Style.Prompt = color.New(color.FgYellow).SprintFunc()
	Style.Success = color.New(color.FgGreen).SprintFunc()
	Style.Info = color.New(color.FgBlue).SprintFunc()
	ErrorHandler = func(err error) {
		fmt.Fprintf(os.Stderr, Style.Error(merry.Details(err))+"\n")
	}
}
