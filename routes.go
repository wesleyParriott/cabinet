package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const IMTHEWORST = "dcc970833371548d5c08360d9c35bcebc1afde0a923d13e994b4f9122043233306f0dbf1ce1227de37b9921385fd8370bb75bd47ba1934a190d278f44032285b"
const CABINETLOCATION = "/usr/local/share/Cabinet/"

// EntryNotAllowed will write a 200 with the text "entry not allowed"
// 200 is used here because laziness
// FIXME use 403
func EntryNotAllowed(response http.ResponseWriter, request *http.Request) {
	_ = request

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

	Logger.Debug("%s", request.RequestURI)

	contents, err := ParsePasscodeTemplate(request.RequestURI)
	if err != nil {
		Logger.Fatal("%s", err.Error())
		// FIXME: 500
	}

	response.Write([]byte(contents))

	return nil
}

// FrontDoor is middleware for routes that handles auth and some
// input handling
func FrontDoor(response http.ResponseWriter, request *http.Request) {
	Logger.Info("Knocking at the front door: %s", request.RequestURI)

	if request.RequestURI == "/favicon.ico" {
		Logger.Info("serving favicon")
		http.ServeFile(response, request, "/usr/local/share/CabinetData/favicon.ico")
		return
	}

	cookie, err := getPasscodeCookie(request)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			err = whatsThePasscode(response, request)
			if err != nil {
				Logger.Fatal("%s", err.Error())
				// FIXME: 500
			}
			return
		default:
			Logger.Error("when trying to get cookie: %s", err)
			// FIXME: 500
		}

		EntryNotAllowed(response, request)
		return
	}
	Logger.Debug("%+v", cookie)

	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		Logger.Fatal("when getting query values: %s", err.Error())
		// FIXME: 500
	}

	Logger.Info("%v", queryValues)

	if strings.ToLower(request.URL.Path) == "/slopmeup" {
		SlopMeUp(response, request)
		return
	}

	whichdir, okay := queryValues["whichdir"]
	if !okay {
		Logger.Info("no whichdir parameter. Redirecting them to the index.")
		Index(response, request)
		return
	}

	if strings.Contains(whichdir[0], "..") {
		Logger.Info("which dir contains ..! can't go backwards :(")
		Index(response, request)
		return
	}

	Logger.Info("letting them in")
	List(response, request, whichdir[0])
}

// Index is a route that serves the index page
// the index page is located in tmpls/index.html
func Index(response http.ResponseWriter, request *http.Request) {

	if request.RequestURI != "/" {
		http.Redirect(response, request, "/", 301)
	}

	index, err := ParseIndexTemplate()
	if err != nil {
		Logger.Fatal("when parsing index template: %s", err.Error())
	}

	response.Write([]byte(index))
}

// List is a route that will serve a parsed page that has all the links to the assumed media files
// in a directory. The directory is determined by the "whichdir" parameter
func List(response http.ResponseWriter, request *http.Request, whichdir string) {
	Logger.Info("Listing files in %s", whichdir)

	_ = request

	fileNames, err := listDir("/usr/local/share/Cabinet/" + whichdir)
	if err != nil {
		Logger.Error("err when listing files: %s", err)
	}

	content, err := ParseListTemplate(whichdir, fileNames)
	if err != nil {
		Logger.Fatal("error when parsing template: %s", err)
		// FIXME: 500
	}

	response.Write([]byte(content))
}

// SlopMeUp is a route that will serve a file based on the soup query parameter
func SlopMeUp(response http.ResponseWriter, request *http.Request) {
	Logger.Info("slopping up some hot soup")
	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		// FIXME: 500 instead of hard fatal
		Logger.Fatal("spilled the soup because: %s", err)
	}
	// FIXME: move to frontdoor?
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
