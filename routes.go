package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"facette.io/natsort"
)

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

			Forbidden(response)
			return
		}

		Logger.Info("checking the passcode")

		passcode := cookie.Value
		if passcode != PASSCODE {
			Logger.Info("passcode incorrect sending them away")
			Forbidden(response)
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
	} else if strings.ToLower(request.URL.Path) == "/stitch" {
		if request.Method != "POST" {
			Logger.Error("Wrong Method: %s", request.Method)
			BadRequest(response)
			return
		}
		Stitch(response, request)
		return
	} else if strings.ToLower(request.URL.Path) == "/mkdir" {
		if request.Method != "POST" {
			Logger.Error("Wrong Method: %s", request.Method)
			BadRequest(response)
			return
		}
		Mkdir(response, request)
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
		InternalError(response)
		return
	}

	content, err := ParseListTemplate(whichdir, fileNames)
	if err != nil {
		Logger.Error("error when parsing template: %s", err)
		InternalError(response)
		return
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

	Logger.Debug("%v", queryValues)
	soups, okay := queryValues["soup"]
	if !okay {
		Logger.Info("no soup parameter. Entry not allowed")
		Forbidden(response)
		return
	}
	soup := soups[0]
	Logger.Debug("soup " + soup)
	path := CABINETLOCATION + soup
	Logger.Debug("path  " + path)
	http.ServeFile(response, request, path)
}

func Mkdir(response http.ResponseWriter, request *http.Request) {
	Logger.Info("Entering Mkdir route")

	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		Logger.Error("during upload %s", err.Error())
		InternalError(response)
		return
	}

	whichdir, okay := queryValues["whichdir"]
	if !okay {
		Logger.Error("query value whichdir wasn't found")
		BadRequest(response)
		return
	}
	Logger.Debug("%s", whichdir)

	newdir, okay := queryValues["newdir"]
	if !okay {
		Logger.Error("query value newdir wasn't found")
		BadRequest(response)
		return
	}
	Logger.Debug("%s", newdir)

	url := CABINETLOCATION + whichdir[0] + "/" + newdir[0]
	Logger.Debug("%s", url)

	err = os.Mkdir(url, 0755)
	if err != nil {
		if os.IsExist(err) {
			Logger.Error("os already exists sending back a 409")
			Conflict(response)
			return
		}
		Logger.Error("when trying to make dir at %s: %s", url, err)
		InternalError(response)
		return
	}

	Created(response, whichdir[0]+"/"+newdir[0])
}

func Upload(response http.ResponseWriter, request *http.Request) {
	Logger.Info("Entering upload route")

	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		Logger.Error("during upload %s", err.Error())
		InternalError(response)
		return
	}
	part, okay := queryValues["part"]
	if !okay {
		Logger.Error("query value part wasn't found")
		BadRequest(response)
		return
	}

	filename, okay := queryValues["filename"]
	if !okay {
		Logger.Error("query value filename wasn't found")
		BadRequest(response)
		return
	}

	maxReader := http.MaxBytesReader(response, request.Body, int64(MB(32)))
	defer maxReader.Close()

	fileNameToSave := fmt.Sprintf("/tmp/%s-%s", part[0], filename[0])

	f, err := os.OpenFile(fileNameToSave, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		Logger.Error(err.Error())
		InternalError(response)
		return
	}
	defer f.Close()

	// Logger.Debug("reading request body")
	// intendedBuffSize := MB(32)
	// buff := make([]byte, intendedBuffSize)
	// readed, err := io.ReadFull(maxReader, buff)
	// if err != nil {
	// 	if !errors.Is(err, io.ErrUnexpectedEOF) {
	// 		Logger.Error("file %s part %s was maybe too big\nerror returned: %s", filename, part, err.Error())
	// 		InternalError(response)
	// 		return
	// 	}
	// }
	// Logger.Debug("read %db from request body", readed)

	// if readed < intendedBuffSize {
	// 	buff = buff[:readed]
	// }

	// Logger.Debug("writing %db to %s", len(buff), fileNameToSave)
	// written, err := f.Write(buff)
	// if err != nil {
	// 	Logger.Error("when writing %s: %s", fileNameToSave, err.Error)
	// 	InternalError(response)
	// 	return
	// }
	// Logger.Debug("wrote %db to %s", written, fileNameToSave)

	Logger.Debug("reading request body and copying to file")
	written, err := io.Copy(f, maxReader)
	if err != nil {
		Logger.Error("file %s part %s was maybe too big\nError returned: %s", filename, part, err.Error())
		InternalError(response)
		return
	}
	Logger.Debug("wrote %db to file", written)

	Okay(response, []byte("job done"))
}

func Stitch(response http.ResponseWriter, request *http.Request) {
	Logger.Info("Entering stitch route")
	queryValues, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		Logger.Error("during upload %s", err.Error())
		InternalError(response)
		return
	}

	whichdir, okay := queryValues["whichdir"]
	if !okay {
		Logger.Error("query value whichdir wasn't found")
		BadRequest(response)
		return
	}
	whichfile, okay := queryValues["whichfile"]
	if !okay {
		Logger.Error("query value whichfile wasn't found")
		BadRequest(response)
		return
	}

	Logger.Debug("which-dir: %s which-file: %s", whichdir, whichfile)

	// list everything in temp
	fileNames, err := listDir("/tmp")
	if err != nil {
		Logger.Error("when listing /tmp: %s", err.Error())
		InternalError(response)
		return
	}

	var parts []string
	// find all the file parts
	for _, fileName := range fileNames {
		if strings.Contains(fileName, whichfile[0]) {
			parts = append(parts, "/tmp/"+fileName)
		}
	}

	// sort based on the first -
	natsort.Sort(parts)
	for _, part := range parts {
		Logger.Debug(part)
	}

	newFilePath := CABINETLOCATION + whichdir[0] + "/" + whichfile[0]

	Logger.Info("stitching %s", newFilePath)

	// create new file in the whichdir
	f, err := os.OpenFile(newFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Error("when creating %s: %s", whichfile[0], err.Error())
	}
	defer f.Close()

	// foreach file read the data of the part and append to the newfile
	for _, part := range parts {
		contents, err := os.ReadFile(part)
		if err != nil {
			Logger.Error("when reading %s: %s", part, err.Error())
			InternalError(response)
			return
		}
		_, err = f.Write(contents)
		if err != nil {
			Logger.Error("when writing %s: %s", part, err.Error())
			InternalError(response)
			return
		}
		os.Remove(part)
	}

	Okay(response, []byte("job done"))
}

// func Upload(response http.ResponseWriter, request *http.Request) {
// 	Logger.Info("Entering upload route")
//
// 	destination := request.Header.Get("X-Destination")
// 	Logger.Info("File Destination: %s", destination)
//
// 	maxUploadSize := int64(1024 * 1024 * 1024 * 1024)
//
// 	err := request.ParseMultipartForm(maxUploadSize)
// 	if err != nil {
// 		Logger.Error("when trying to parse multipart form: %s", err.Error())
// 		InternalError(response)
// 		return
// 	}
//
// 	files := request.MultipartForm.File["file"]
//
// 	for _, fileHeader := range files {
// 		Logger.Debug("FILE HEADER INFO\n\tFile Name: %s\n\t File Size: %d", fileHeader.Filename, fileHeader.Size)
//
// 		if fileHeader.Size > maxUploadSize {
// 			Logger.Error("%s too big", fileHeader.Filename)
// 			EntityTooLarge(response)
// 			return
// 		}
//
// 		filePath := fmt.Sprintf("%s/%s/%s", CABINETLOCATION, destination, fileHeader.Filename)
//
// 		givenFile, err := fileHeader.Open()
// 		if err != nil {
// 			Logger.Error("when trying to open fileHeader: %s", fileHeader.Filename)
// 			InternalError(response)
// 			return
// 		}
//
// 		Logger.Info("writing file to %s", filePath)
// 		f, err := os.Create(filePath)
// 		if err != nil {
// 			Logger.Error(err.Error())
// 			InternalError(response)
// 			return
// 		}
// 		defer f.Close()
//
// 		_, err = io.Copy(f, givenFile)
// 		if err != nil {
// 			Logger.Error(err.Error())
// 			InternalError(response)
// 			return
// 		}
//
// 	}
//
// 	Okay(response, []byte("got it"))
// }
