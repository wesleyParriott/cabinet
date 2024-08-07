package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/wesleyParriott/wlog"
)

var PASSCODE string
var Logger wlog.Wlog

func init() {
	initFlags()
	Logger = wlog.CreateWlogWithParams(os.Stdout, wlog.DEBUG)
	err := setPasscode()
	if err != nil {
		Logger.Fatal("during init couldn't set passcode because: %s", err)
	}
}

func main() {
	version := "v0.2.1"

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
	Logger.Info("Listening and Serving on %s ...", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		Logger.Fatal("%s", err)
	}
}
