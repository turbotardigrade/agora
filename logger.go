// Adapted from
// https://www.goinggo.net/2013/11/using-log-package-in-go.html

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogger() {
	// Initialize Logger
	if opts.Silent {
		// If --silent is set log to a file instead
		logPath := LogDirPath + "/log.txt"
		err := os.MkdirAll(LogDirPath, 0755)
		if err != nil {
			log.Fatalln("Failed to create log directory", err)
		}

		err = CreateFileIfNotExists(logPath)
		if err != nil {
			log.Fatalln("Failed to create file", err)
		}

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", err)
		}

		file.WriteString("----------------------------------------\n")
		file.WriteString("   Log of " + time.Now().Format("Jan _2 15:04") + "\n")
		file.WriteString("----------------------------------------\n")
		LoggerInit(file, file, file, file)
	} else {
		LoggerInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}
}

func LoggerInit(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Lshortfile)
}
