package main

import (
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/joshmedeski/sesh/seshcli"
	"github.com/urfave/cli/v2"
)

type Man struct {
	Date time.Time
	App  cli.App
}

func main() {
	version := "dev"
	templateFile := "man.tmpl"
	manPageName := "sesh.1"

	man := &Man{
		Date: time.Now(),
		App:  seshcli.App(version),
	}

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	template, err := template.New(templateFile).Funcs(funcMap).ParseFiles(templateFile)
	if err != nil {
		log.Fatal("can't parse file")
	}

	outputFile, err := os.Create(manPageName)
	if err != nil {
		log.Fatal("error creating file:", err)
	}
	defer outputFile.Close()

	err = template.Execute(outputFile, man)
	if err != nil {
		log.Fatal("error generating man page")
	}

}
