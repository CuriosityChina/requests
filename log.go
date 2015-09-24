package requests

import (
	"io"
	"io/ioutil"
	"log"
)

var (
	globalLogger *log.Logger = log.New(ioutil.Discard, "[debug] ", log.LstdFlags|log.Lshortfile)
	globalDebug  bool
)

func SetLogger(logOutput io.Writer) {
	globalLogger = log.New(logOutput, "[debug] ", log.LstdFlags|log.Lshortfile)
}

func SetDebug(debug bool) {
	globalDebug = debug
}

func debugln(args ...interface{}) {
	if globalDebug && globalLogger != nil {
		globalLogger.Println(args)
	}
}

func debugf(message string, args ...interface{}) {
	if globalDebug && globalLogger != nil {
		globalLogger.Printf(message+"\n", args...)
	}
}
