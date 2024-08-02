package main

import (
	"errors"
	"net/http"
	"net/url"
)

const IMTHEWORST = "dcc970833371548d5c08360d9c35bcebc1afde0a923d13e994b4f9122043233306f0dbf1ce1227de37b9921385fd8370bb75bd47ba1934a190d278f44032285b"
const CABINETLOCATION = "/usr/local/share/Cabinet/"

func EntryNotAllowed(response http.ResponseWriter, request *http.Request) {
	_ = request

	// TODO: 403
	response.Write([]byte("entry not allowed"))
}

func getPasscodeCookie(r *http.Request) (http.Cookie, error) {
	cookie, err := r.Cookie("passcode")
	if err != nil {
		return http.Cookie{}, err
	}
	return *cookie, nil
}

func whatsThePasscode(response http.ResponseWriter, request *http.Request) error {
	// TODO: passcode template
	_ = request

	response.Write([]byte("oi what's the passcode"))

	return nil
}

func FrontDoor(response http.ResponseWriter, request *http.Request) {
	Logger.Info("Knocking at the front door")

	cookie, err := getPasscodeCookie(request)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			err = whatsThePasscode(response, request)
			if err != nil {
				Logger.Fatal("%s", err.Error())
			}
			return
		default:
			Logger.Error("when trying to get cookie: %s", err)
			// TODO: 500
		}

		EntryNotAllowed(response, request)
		return
	}
	Logger.Debug("%+v", cookie)

	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		Logger.Fatal("bummer: %s", err)
	}
	Logger.Info("%v", queryValues)
	secrets, okay := queryValues["password"]
	if !okay {
		Logger.Info("no password parameter. Entry not allowed")
		EntryNotAllowed(response, request)
		return
	}
	whichdir, okay := queryValues["whichdir"]
	if !okay {
		Logger.Info("no whichdir parameter. Entry not allowed")
		EntryNotAllowed(response, request)
		return
	}
	Logger.Info("found secret password")
	if secrets[0] == IMTHEWORST {
		Logger.Info("letting them in")
		List(response, request, whichdir[0])
	} else {
		Logger.Info("denying entry")
		EntryNotAllowed(response, request)
		return
	}
}

func List(response http.ResponseWriter, request *http.Request, whichdir string) {
	_ = request

	fileNames, err := listDir("/usr/local/share/Cabinet/" + whichdir)
	if err != nil {
		Logger.Error("err when listing files: %s", err)
	}

	content, err := ParseMainTemplate(whichdir, fileNames)
	if err != nil {
		Logger.Fatal("error when parsing template: %s", err)
	}

	response.Write([]byte(content))
}

func SlopMeUp(response http.ResponseWriter, request *http.Request) {
	Logger.Info("slopping up some hot soup")
	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		Logger.Fatal("spilled the soup because: %s", err)
	}
	Logger.Info("%v", queryValues)
	soups, okay := queryValues["soup"]
	if !okay {
		Logger.Info("no soup parameter. Entry not allowed")
		EntryNotAllowed(response, request)
		return
	}
	soup := soups[0]
	Logger.Debug("soup " + soup)
	path := CABINETLOCATION + soup
	Logger.Debug("path  " + path)
	http.ServeFile(response, request, path)
}
