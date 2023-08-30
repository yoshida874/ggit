package main

import (
	"fmt"
	"os"
	"yoshida874/ggit/add"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: ggit <command> [<args>]")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		err := Init()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "add":
		if len(args) < 2 {
			fmt.Println("usage: ggit add <file>")
			os.Exit(1)
		}
		err := add.Add(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		fmt.Println("Unknown command:", args[0])
		os.Exit(1)
	}
}

func Init() error {
	err := os.Mkdir("./ggit", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .ggit directory: %s", err)
	}

	dirs := []string{".ggit/objects", ".ggit/refs/heads", ".ggit/refs/tags"}
    for _, d := range dirs {
        if err := os.MkdirAll(d, 0755); err != nil {
            return err
        }
    }

	headFile := ".ggit/HEAD"
	file, err := os.Create(headFile)
	if err != nil {
		return fmt.Errorf("failed to create .ggit/HEAD file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString("ref: refs/heads/master\n")
	if err != nil {
		return fmt.Errorf("failed to write to .ggit/HEAD file: %s", err)
	}

	fmt.Println("Initialized empty ggit repository in .ggit/")
	return nil
}
