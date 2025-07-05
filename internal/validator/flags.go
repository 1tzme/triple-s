package validator

import (
	"flag"
	"fmt"
)

func InitFlags() (string, string, bool) {
	port, dir, help := "", "", false

	flag.StringVar(&port, "port", "8080", "Port number")
	flag.StringVar(&dir, "dir", "./data", "Path to directory")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	return port, dir, help
}

func PrintUsage() {
	fmt.Println(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)
}
