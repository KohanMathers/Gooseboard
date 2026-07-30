package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: gooseboard <run|init|version> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("gooseboard v0.0.1")
	case "run":
		if len(os.Args) < 4 {
			fmt.Println("usage: gooseboard run <template_dir> <port>")
			os.Exit(1)
		}
		templateDir := os.Args[2]
		port, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("invalid port: %v\n", err)
			os.Exit(1)
		}
		if err := runCommand(templateDir, port); err != nil {
			fmt.Printf("error running gooseboard: %v\n", err)
			os.Exit(1)
		}
	case "init":
		fmt.Println("not implemented yet")
	default:
		fmt.Printf("unknown command %q\n", os.Args[1])
		os.Exit(1)
	}
}
