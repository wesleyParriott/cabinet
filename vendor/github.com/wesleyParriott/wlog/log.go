package wlog

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Wlog struct {
	Writer io.Writer
	Level  int
}

const (
	DEBUG int = iota + 1
	INFO
	ERROR
	FATAL
)

const (
	normalText  = "\033[0m"
	boldText    = "\033[1m"
	blackText   = "\033[30;1m"
	redText     = "\033[31;1m"
	greenText   = "\033[32;1m"
	yellowText  = "\033[33;1m"
	blueText    = "\033[34;1m"
	magentaText = "\033[35;1m"
	cyanText    = "\033[36;1m"
	whiteText   = "\033[37;1m"
)

func CreateWlogWithParams(writer io.Writer, level int) Wlog {

	var wlog Wlog
	wlog.Writer = writer
	wlog.Level = level
	return wlog

}

func (w Wlog) Debug(message string, values ...interface{}) {
	if w.Level <= DEBUG {
		currentTime := time.Now().UTC().String()
		message = fmt.Sprintf("%s[DEBUG %s]%s %s\n", blueText, currentTime, normalText, message)
		message = fmt.Sprintf(message, values...)
		w.Writer.Write([]byte(message))
	}
}
func (w Wlog) Info(message string, values ...interface{}) {
	currentTime := time.Now().UTC().String()
	if w.Level <= INFO {
		message := fmt.Sprintf("%s[INFO %s]%s %s\n", greenText, currentTime, normalText, message)
		fmt.Printf(message, values...)
	}
}
func (w Wlog) Error(message string, values ...interface{}) {
	currentTime := time.Now().UTC().String()
	if w.Level <= ERROR {
		message := fmt.Sprintf("%s[ERROR %s]%s %s\n", redText, currentTime, normalText, message)
		fmt.Printf(message, values...)
	}
}
func (w Wlog) Fatal(message string, values ...interface{}) {
	currentTime := time.Now().UTC().String()
	if w.Level <= FATAL {
		message := fmt.Sprintf("%s[FATAL %s]%s %s\n", redText, currentTime, normalText, message)
		fmt.Printf(message, values...)
		os.Exit(1)
	}
}
