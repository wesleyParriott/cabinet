package main

import (
	"bytes"
	"os"
	"text/template"
)

const TemplateFolder = "./tmpls/"

func ParseMainTemplate(whichdir string, paths []string) (string, error) {
	type mainInfo struct {
		Title    string
		Whichdir string
		Paths    []string
	}

	main := mainInfo{
		"Cabinet",
		whichdir,
		paths,
	}

	mainFileContents, err := os.ReadFile(TemplateFolder + "main.html")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("main").Parse(string(mainFileContents))
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, main)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
