package main

import (
	"bytes"
	"os"
	"text/template"
)

const TemplateFolder = "./tmpls/"

func ParseIndexTemplate() (string, error) {
	type Info struct {
		Directories []string
	}

	_, dirs, err := listDir(CABINETLOCATION)

	info := Info{
		dirs,
	}

	indexFileContents, err := os.ReadFile(TemplateFolder + "index.html")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("index").Parse(string(indexFileContents))
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, info)
	if err != nil {
		return "", err
	}

	return parseMainTemplate("Index", buff.String())

}
func ParseListTemplate(whichdir string, fileNames []string, dirNames []string) (string, error) {
	type Info struct {
		Whichdir    string
		Files       []string
		Directories []string
	}

	info := Info{
		whichdir,
		fileNames,
		dirNames,
	}

	listFileContents, err := os.ReadFile(TemplateFolder + "list.html")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("list").Parse(string(listFileContents))
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, info)
	if err != nil {
		return "", err
	}

	return parseMainTemplate(whichdir, buff.String())
}

func ParsePasscodeTemplate(route string) (string, error) {

	Logger.Debug("route passed to ParsePasscodeTemplate: %s", route)

	type Info struct {
		Route string
	}

	info := Info{
		route,
	}

	passcodeFileContents, err := os.ReadFile(TemplateFolder + "passcode.html")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("list").Parse(string(passcodeFileContents))
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, info)
	if err != nil {
		return "", err
	}

	return parseMainTemplate("Passcode", buff.String())
}

func parseMainTemplate(title string, maincontent string) (string, error) {
	type mainInfo struct {
		Title       string
		MainContent string
	}

	main := mainInfo{
		"Cabinet | " + title,
		maincontent,
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
