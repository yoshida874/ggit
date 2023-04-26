package main

import (
	"fmt"
	"os"
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
			return
		}
	default:
		fmt.Println("Unknown command:", args[0])
		os.Exit(1)
	}
}

func Init() error {
	err := os.Mkdir(".ggit", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .ggit directory: %s", err)
	}

	objectsDir := ".ggit/objects"
	err = os.Mkdir(objectsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .ggit/objects directory: %s", err)
	}

	indexDir := ".ggit/.index"
	err = os.Mkdir(indexDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .ggit/.index directory: %s", err)
	}

	refsDir := ".ggit/refs"
	err = os.Mkdir(refsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .ggit/refs directory: %s", err)
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
