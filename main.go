package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

var tmplDir string
var filePattern string
var dirNodes []string
var tmplData map[string]string

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	tmplDir = os.Getenv("TMPL_DIR")
	if tmplDir == "" {
		log.Fatalf("TMPL_DIR environment variable not set")
	}

	if len(os.Args) == 1 {
		listDir(tmplDir, 0)
		os.Exit(0)
	}

	tmplData = make(map[string]string, 32)

	parseArgs()
	printTmpl()
}

func parseArgs() {
	inData := false
	args := os.Args[1:]

	for _, arg := range args {
		if !inData && arg == "--" {
			inData = true
			continue
		}

		if !inData {
			dirNodes = append(dirNodes, arg)
			continue
		}

		ss := strings.Split(arg, "=")
		if len(ss) == 1 {
			log.Fatalf("missing '=' in '%s' value assignment", arg)
		}
		tmplData[ss[0]] = ss[1]
	}

	filePattern = dirNodes[len(dirNodes)-1]
	dirNodes = dirNodes[0 : len(dirNodes)-1]
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

	fileNames := []string{}
	fileCount := 0
	for _, e := range entries {
		if strings.Contains(e.Name(), filePattern) {
			fileNames = append(fileNames, e.Name())
			fileCount++
		}
	}

	if fileCount == 0 {
		log.Fatalf("file matching pattern '%s' not found in directory %s", filePattern, dir)
	} else if fileCount > 1 {
		files := ""
		for i, f := range fileNames {
			files += fmt.Sprintf("'%s'", f)
			if i < len(fileNames)-1 {
				files += ", "
			}
		}
		log.Fatalf("%d files match pattern '%s': %s", fileCount, filePattern, files)
	}

	fileName := fileNames[0]

	filePath := path.Join(dir, fileName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("can't read template file: %v", err)
	}
	str := string(bytes)

	cursorAt := ""
	if strings.HasPrefix(str, "cursor@") {
		ss := strings.SplitN(str, "\n", 2)
		cursorAt = ss[0][7:]
		str = ss[1]
	}

	tmpl := template.New("tmpl")
	tmpl, err = tmpl.Parse(str)
	if err != nil {
		log.Fatalf("can't parse template: %v", err)
	}

	b := strings.Builder{}
	err = tmpl.Execute(&b, &tmplData)
	if err != nil {
		log.Fatalf("can't execute template: %v", err)
	}

	fmt.Printf("%s", b.String())

	if cursorAt != "" {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s", cursorAt))
	}
}
