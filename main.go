package main

import (
	"log"
	"os"
)

func main() {
	tmplDir := os.Getenv("TMPL_DIR")
	if tmplDir == "" {
		log.Fatalf("TMPL_DIR environment variable not set")
	}
}
