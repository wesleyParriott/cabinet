package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/wesleyParriott/wlog"
)

var Logger wlog.Wlog

func init() {
	initFlags()
	Logger = wlog.CreateWlogWithParams(os.Stdout, wlog.DEBUG)
}

func main() {
	version := "v0.2.0"

	flag.Parse()
	if HelpFlag {
		PrintUsage()
	}
	if SetupFlag && BreakdownFlag {
		Logger.Fatal("can't do both setup and breakdown. Please choose either -b or -s")
	}
	if SetupFlag {
		Setup()
	}
	if BreakdownFlag {
		Breakdown()
	}
	if VersionFlag {
		println(version)
		os.Exit(0)
	}

	port := ":3000"
	http.HandleFunc("/", FrontDoor)
	http.HandleFunc("/slopmeup", SlopMeUp)
	Logger.Info("Listening and Serving on %s ...", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		Logger.Fatal("%s", err)
	}
}
