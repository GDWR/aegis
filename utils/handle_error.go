package utils

import (
	"log"
	"os"
	"runtime"
)

func HandleError(err error) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", filename, line, err)
		os.Exit(-1)
	}
}
