package Frodo

import (
	"bytes"
	"github.com/fatih/color"
	"log"
	"os"
)

type Logger struct {
	FilePath, FileName string
	LogFile            *log.Logger
	buffers            *bytes.Buffer
}

var Log *Logger

func init() {
	Log = &Logger{
		FilePath: "./logs/",
		FileName: "frodo.log",
	}

	if Log.LogFile == nil {
		Log.Initialise()
	}
}

func (console *Logger) WriteToFile(fl ...interface{}) (*log.Logger, error) {

	// If arguements are provided, then expect:
	// 1st argument is FilePath
	if len(fl) > 0 {
		if fp, ok := fl[0].(string); ok {
			console.FilePath = fp
		}

		// 2nd argument is FileName
		if len(fl) > 1 {
			if fn, ok := fl[1].(string); ok {
				console.FileName = fn
			}
		}
	}

	err := os.Mkdir(console.FilePath, 0775)
	if err != nil {
		Log.Error("Error: Failed to create folder %s", console.FilePath)
	}

	// 1st check to see if the path and filename provided exists
	// create a new file if none exists.
	file, err := os.OpenFile(console.FilePath+console.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		console.Error("------ Failed to open log file. Filepath: %s ------\n", console.FilePath)
		return nil, err
	}

	console.LogFile = log.New(file, "Frodo: ", log.Ldate|log.Ltime|log.Lshortfile)
	console.Success("------- File Logging activated!! -------\n")
	return console.LogFile, nil
}

func (console *Logger) Initialise() {
	// For now collect all the buffers on to a buffer memory
	buffer := new(bytes.Buffer)
	console.LogFile = log.New(buffer, "Frodo: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (console *Logger) Info(format string, args ...interface{}) {
	console.log(color.FgCyan, false, format, args...)
}

func (console *Logger) Debug(format string, args ...interface{}) {
	console.log(color.FgWhite, false, format, args...)
}

func (console *Logger) Success(format string, args ...interface{}) {
	console.log(color.FgGreen, false, format, args...)
}

func (console *Logger) Warn(format string, args ...interface{}) {
	console.log(color.FgYellow, false, format, args...)
}

func (console *Logger) Error(format string, args ...interface{}) {
	console.log(color.FgRed, false, format, args...)
}

func (console *Logger) Alert(format string, args ...interface{}) {
	color.Set(color.BgMagenta, color.FgWhite, color.Bold)
	defer color.Unset()
	console.LogFile.Printf(format, args...)
}

func (console *Logger) Critical(format string, args ...interface{}) {
	color.Set(color.BgRed, color.FgWhite, color.Bold)
	defer color.Unset()
	console.LogFile.Printf(format, args...)
}

func (console *Logger) Fatal(format string, args ...interface{}) {
	color.Set(color.BgRed, color.FgWhite, color.Bold)
	defer color.Unset()
	console.LogFile.Fatalf(format, args...)
}

func (console *Logger) log(colorAttr color.Attribute, isBold bool, format string, args ...interface{}) {
	newlog := color.Set(colorAttr)
	defer color.Unset()
	if isBold {
		newlog.Add(color.Bold)
	}

	// I want it log both into the file and on the console
	log.Printf(format, args...)
	console.LogFile.Printf(format, args...)
}

func (console *Logger) Dump() {
	console.LogFile.Print(console.buffers)
}
