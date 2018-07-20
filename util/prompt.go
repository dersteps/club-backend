package util

import (
	"log"

	"github.com/logrusorgru/aurora"
)

func Info(text string) {
	log.Printf("%s %s\n", aurora.BgBlue(" CREW | INFO    "), text)
}

func Warn(text string) {
	log.Printf("%s %s\n", aurora.BgBrown(" CREW | WARN    "), text)
}

func Fatal(text string) {
	log.Printf("%s %s\n", aurora.BgRed(" CREW | ERROR   "), text)
}

func Success(text string) {
	log.Printf("%s %s\n", aurora.BgGreen(" CREW | SUCCESS "), text)
}
