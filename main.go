package main

import (
	"net/http"
	"os"

	"github.com/wesleyParriott/wlog"
)

var PASSCODE string
var Logger wlog.Wlog

func init() {
	initFlags()

	var logLevel int

	if LogLevel == "debug" {
		// the default of the flag
		// so if they don't set anything it should go here
		logLevel = wlog.DEBUG
	} else if LogLevel == "info" {
		logLevel = wlog.INFO
	} else if LogLevel == "error" {
		logLevel = wlog.ERROR
	} else if LogLevel == "fatal" {
		logLevel = wlog.FATAL
	} else {
		println("Sorry didn't recognize " + LogLevel + " as a log level. Please use these:")
		println("\tdebug, info, error, or fatal")
		os.Exit(0)
	}

	Logger = wlog.CreateWlogWithParams(os.Stdout, logLevel)
	err := setPasscode()
	if err != nil {
		Logger.Fatal("during init couldn't set passcode because: %s", err)
	}
}

func main() {
	version := "v0.2.1"

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
