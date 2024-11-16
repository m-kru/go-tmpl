package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var tmplDir string
var filePattern string
var dirNodes []string

func main() {
	log.SetFlags(0)

	tmplDir = os.Getenv("TMPL_DIR")
	if tmplDir == "" {
		log.Fatalf("TMPL_DIR environment variable not set")
	}

	args := os.Args[1:]
	if len(args) == 0 {
		listDir(tmplDir, 0)
		os.Exit(0)
	} else if len(args) > 1 {
		dirNodes = args[0 : len(args)-1]
	}

	filePattern = args[len(args)-1]

	printTmpl()
}

func listDir(dir string, lvl int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("can't list files in %s: %v", dir, err)
	}

	for _, e := range entries {
		for _ = range lvl {
			fmt.Printf("  ")
		}
		fmt.Printf("%s\n", e.Name())
		if e.IsDir() {
			listDir(path.Join(dir, e.Name()), lvl+1)
		}
	}
}

func printTmpl() {
	dir := path.Join(dirNodes...)
	dir = path.Join(tmplDir, dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("can't read files in %s: %v", dir, err)
	}

	fileName := ""
	for _, e := range entries {
		if strings.Contains(e.Name(), filePattern) {
			fileName = e.Name()
			break
		}
	}

	if fileName == "" {
		log.Fatalf("file matching pattern '%s' not found in directory %s", filePattern, dir)
	}

	filePath := path.Join(dir, fileName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("can't read template file: %v", err)
	}
	str := string(bytes)

	fmt.Printf("%s", str)
}
