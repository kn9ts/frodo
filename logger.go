package Frodo

import (
	"bytes"
	"github.com/fatih/color"
	"log"
	"os"
)

// Logger embeds *log.Logger thus is able to handle all logs
type Logger struct {
	FilePath, FileName string
	LogFile            *log.Logger
	buffers            *bytes.Buffer
	level              int
}

type logDetails struct {
	color color.Attribute
	level int
}

// Available Log levels for the app
var logLevel = map[string]logDetails{
	"Fatal":   {color.FgRed, 1},
	"Alert":   {color.FgRed, 2},
	"Error":   {color.FgRed, 3},
	"Success": {color.FgGreen, 4},
	"Warn":    {color.FgYellow, 5},
	"Info":    {color.FgCyan, 6},
	"Debug":   {color.FgWhite, 7},
}

// Log will be global var that handles Frodo's logging
var Log *Logger

func init() {
	Log = &Logger{
		FilePath: "./logs/",
		FileName: "frodo.log",
	}

	// Set a default log level
	Log.level = 5
	Log.FilePath = "./"

	// Initialise logging
	if Log.LogFile == nil {
		Log.Initialise()
	}
}

// Initialise checks if a logging instance exists, if not commence one
func (console *Logger) Initialise() {
	// For now collect all the buffers on to a buffer memory
	buffer := new(bytes.Buffer)
	console.LogFile = log.New(buffer, "[Frodo] ", log.Ldate|log.Ltime|log.Lshortfile)
	log.SetPrefix("[Frodo] ")
	// log.SetFlags(log.Lshortfile)
}

// Info can be used to log informative information of the application
func (console *Logger) Info(format string, args ...interface{}) {
	console.log(logLevel["Info"], color.FgCyan, false, format, args...)
}

// Debug logs debug information
func (console *Logger) Debug(format string, args ...interface{}) {
	console.log(logLevel["Debug"], color.FgWhite, false, format, args...)
}

// Success can be used to log information on successful transactions
// logs in green
func (console *Logger) Success(format string, args ...interface{}) {
	console.log(logLevel["Success"], color.FgGreen, false, format, args...)
}

// Warn can be used  to log meaningful and light errors/bugs that might want to be checked later
func (console *Logger) Warn(format string, args ...interface{}) {
	console.log(logLevel["Warn"], color.FgYellow, false, format, args...)
}

// Error can be used to log Errors
// red in color
func (console *Logger) Error(format string, args ...interface{}) {
	console.log(logLevel["Error"], color.FgRed, false, format, args...)
}

// Alert can be used to log Alerts, maybe on certain events!!
// Magenta background, white text
func (console *Logger) Alert(format string, args ...interface{}) {
	color.Set(color.BgMagenta, color.FgBlack, color.Bold)
	defer color.Unset()
	console.log(logLevel["Alert"], color.FgBlack, true, format, args...)
}

// Critical can be used to log system wide Critical information that needs to be fixed immediately
// Red background and white text
func (console *Logger) Critical(format string, args ...interface{}) {
	color.Set(color.BgRed, color.FgWhite, color.Bold)
	defer color.Unset()
	console.log(logLevel["Fatal"], color.FgRed, true, format, args...)
}

// Fatal is similar to Frodo.Log.Error but panics after logging
func (console *Logger) Fatal(format string, args ...interface{}) {
	console.Critical(format, args...)
}

// main logging handler
func (console *Logger) log(d logDetails, colorAttr color.Attribute, isBold bool, format string, args ...interface{}) {
	// Only print out, logs that are the selected level or higher
	if console.level >= d.level {
		newlog := color.Set(colorAttr)
		defer color.Unset()

		// if the request is to log bold text
		if isBold {
			newlog.Add(color.Bold)
		}

		// Check to see if we are to add background color
		// Alert and Fatal have background colors
		switch d.level {
		case 1: // for Fatal
			color.Set(color.BgRed)
			break
		case 2: // for Alerts
			color.Set(color.BgMagenta)
			break
		}

		// I want it log both into the file and on the console
		console.LogFile.Printf(format, args...)
		log.Printf(format, args...)
	}
}

// WriteToFile prompts all logs to be written to a file
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

	dir, err := os.Stat(console.FilePath)
	if err != nil {
		console.Fatal("Error: Directory %s does not exist.", console.FilePath)
	}

	if !dir.IsDir() {
		console.Fatal("Error: %s is not a directory/folder.", console.FilePath)
	}
	console.Info("Debug information will be logged at: %s/%s", dir.Name(), console.FileName)

	// 1st check to see if the path and filename provided exists
	// create a new file if none exists.
	file, err := os.OpenFile(console.FilePath+console.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		console.Error("------ Failed to open log file. Filepath: %s ------\n", console.FilePath)
		return nil, err
	}

	console.LogFile = log.New(file, "Frodo: ", log.Ldate|log.Ltime|log.Lshortfile)
	console.Info("------- File Logging activated!! -------\n")
	return console.LogFile, nil
}

// SetLevel enables the user to select a log level for the app to monitor the app
func (console *Logger) SetLevel(option string) {
	// Check to see if the log level exists
	if _, ok := logLevel[option]; ok {
		console.level = logLevel[option].level
	}
}

// Dump simply does just that. Dumps all the logging that has been collected by the buffer
func (console *Logger) Dump() {
	console.LogFile.Print(console.buffers)
}
