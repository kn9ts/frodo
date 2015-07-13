package Frodo

import (
    "github.com/fatih/color"
    "log"
    "os"
)

type Logger struct {
    FilePath, FileName string
    LogFile            *log.Logger
}

var Log *Logger

func init() {
    Log = &Logger{
        FilePath: "./logs/",
        FileName: "frodo.log",
    }
}

func (console *Logger) WriteToFile() (*log.Logger, error) {
    file, err := os.OpenFile(console.FilePath+console.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        clog := color.New(color.FgRed).Add(color.Bold)
        clog.Println("------ Failed to open Logger file ------\n")
        return nil, err
    }

    color.Set(color.FgGreen).Add(color.Bold)
    console.LogFile = log.New(file, "Frodo: ", log.Ldate|log.Ltime|log.Lshortfile)
    log.Printf("\n------- File Logging activated!! -------\n")
    return console.LogFile, nil
}

func (console *Logger) Info(format string, args ...interface{}) {
    console.printOut(color.FgCyan, false, format, args...)
}

func (console *Logger) Notice(format string, args ...interface{}) {
    console.printOut(color.FgGreen, false, format, args...)
}

func (console *Logger) Warn(format string, args ...interface{}) {
    console.printOut(color.FgYellow, false, format, args...)
}

func (console *Logger) Error(format string, args ...interface{}) {
    console.printOut(color.FgRed, false, format, args...)
}

func (console *Logger) Alert(format string, args ...interface{}) {
    color.Set(color.BgMagenta, color.FgWhite, color.Bold)
    log.Printf(format, args...)
}

func (console *Logger) Critical(format string, args ...interface{}) {
    color.Set(color.BgRed, color.FgWhite, color.Bold)
    log.Printf(format, args...)
}

func (console *Logger) printOut(colorAttr color.Attribute, isBold bool, format string, args ...interface{}) {
    newlog := color.Set(colorAttr)
    if isBold {
        newlog.Add(color.Bold)
    }
    log.Printf(format, args...)
}
