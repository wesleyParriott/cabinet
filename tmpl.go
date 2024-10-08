package main

import (
	"bytes"
	"os"
	"text/template"
)

const TemplateFolder = "./tmpls/"

func ParseIndexTemplate() (string, error) {
	type Info struct {
		Directories            []string
		CreateDirectoryForm    string
		CreateDirFunctionality string
	}

	_, dirs, err := listDir(CABINETLOCATION)

	createDirectoryForm, err := getCreateFormDirHTML()
	if err != nil {
		Logger.Error("when getting the create directory dir html %s", err)
		return "", err
	}

	createDirectoryFunctionality, err := getCreateFormDirJS()
	if err != nil {
		Logger.Error("when getting the create directory dir js %s", err)
		return "", err
	}

	info := Info{
		dirs,
		createDirectoryForm,
		createDirectoryFunctionality,
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
		Whichdir               string
		Files                  []string
		Directories            []string
		CreateDirectoryForm    string
		CreateDirFunctionality string
	}

	createDirectoryForm, err := getCreateFormDirHTML()
	if err != nil {
		Logger.Error("when getting the create directory dir html %s", err)
		return "", err
	}

	createDirectoryFunctionality, err := getCreateFormDirJS()
	if err != nil {
		Logger.Error("when getting the create directory dir js %s", err)
		return "", err
	}

	info := Info{
		whichdir,
		fileNames,
		dirNames,
		createDirectoryForm,
		createDirectoryFunctionality,
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

func getCreateFormDirHTML() (string, error) {
	createDirForm, err := os.ReadFile(TemplateFolder + "createformdir.html")
	if err != nil {
		return "", err
	}

	return string(createDirForm), err
}

func getCreateFormDirJS() (string, error) {
	createDirForm, err := os.ReadFile(TemplateFolder + "createformdir.js")
	if err != nil {
		return "", err
	}

	return string(createDirForm), err
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
