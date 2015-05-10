package main

import (
	"github.com/Sirupsen/logrus"
	. "github.com/gabrielhartmann/tftp/tftp"
)

func main() {
	if err := StartNewReqSession(); err != nil {
		logrus.Fatalf("%v", err)
	}
}
