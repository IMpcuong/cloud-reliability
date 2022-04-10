package main

import (
	"io"
	"log"
)

var (
	Trace   *log.Logger // Logging the traceback of errors.
	Info    *log.Logger // Logging the info of errors.
	Warning *log.Logger // Logging the warning may cause errors.
	Error   *log.Logger // Logging the errors itself.
)

func GenerateLogger(
	traceLogger io.Writer,
	infoLogger io.Writer,
	warnLogger io.Writer,
	errorLogger io.Writer) {
	Trace = log.New(traceLogger, "TRACE: ", log.Ltime|log.Lshortfile)

	Info = log.New(infoLogger, "INFO: ", log.Ltime|log.Lshortfile)

	Warning = log.New(warnLogger, "WARNING: ", log.Ltime|log.Lshortfile)

	Error = log.New(errorLogger, "ERROR: ", log.Ltime|log.Lshortfile)
}
