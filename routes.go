package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

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

	Logger.Info("asking for the passcode")
	Logger.Debug("request uri: %s", request.RequestURI)

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

	if request.RequestURI == "/passcode.jpg" {
		Logger.Info("serving passcode jpg")
		http.ServeFile(response, request, "/usr/local/share/CabinetData/passcode.jpg")
		return
	}

	if PublicFlag == false {
		Logger.Info("seeing if there's a passcode already")
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

		Logger.Info("checking the passcode")

		passcode := cookie.Value
		if passcode != PASSCODE {
			Logger.Info("passcode incorrect sending them away")
			EntryNotAllowed(response, request)
			return
		}

		Logger.Info("passcode correct sending them in")
		Logger.Debug("cookie: %+v", cookie)
	}

	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		Logger.Fatal("when getting query values: %s", err.Error())
		// FIXME: 500
	}

	Logger.Debug("query values: %v", queryValues)

	if strings.ToLower(request.URL.Path) == "/slopmeup" {
		SlopMeUp(response, request)
		return
	} else if strings.ToLower(request.URL.Path) == "/upload" {
		if request.Method != "POST" {
			Logger.Error("Wrong Method: %s", request.Method)
			BadRequest(response)
			return
		}
		Upload(response, request)
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

	Okay(response, []byte(content))
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

func Upload(response http.ResponseWriter, request *http.Request) {
	Logger.Info("Entering upload route")

	destination := request.Header.Get("X-Destination")
	Logger.Info("File Destination: %s", destination)

	maxUploadSize := int64(1024 * 1024 * 1024 * 1024)

	err := request.ParseMultipartForm(maxUploadSize)
	if err != nil {
		Logger.Error("when trying to parse multipart form: %s", err.Error())
		InternalError(response)
		return
	}

	files := request.MultipartForm.File["file"]

	for _, fileHeader := range files {
		Logger.Debug("FILE HEADER INFO\n\tFile Name: %s\n\t File Size: %d", fileHeader.Filename, fileHeader.Size)

		if fileHeader.Size > maxUploadSize {
			Logger.Error("%s too big", fileHeader.Filename)
			EntityTooLarge(response)
			return
		}

		filePath := fmt.Sprintf("%s/%s/%s", CABINETLOCATION, destination, fileHeader.Filename)

		givenFile, err := fileHeader.Open()
		if err != nil {
			Logger.Error("when trying to open fileHeader: %s", fileHeader.Filename)
			InternalError(response)
			return
		}

		Logger.Info("writing file to %s", filePath)
		f, err := os.Create(filePath)
		if err != nil {
			Logger.Error(err.Error())
			InternalError(response)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, givenFile)
		if err != nil {
			Logger.Error(err.Error())
			InternalError(response)
			return
		}

	}

	Okay(response, []byte("got it"))
}
