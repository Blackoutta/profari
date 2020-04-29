package profari

import (
	"fmt"
	"io"
	"log"
	"os"
)

func NewLogger(suiteName string) (*log.Logger, *os.File, error) {
	logpath := "./logs"
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir(logpath, os.ModePerm|os.ModeDir)
	}

	f, err := os.OpenFile("./logs/"+suiteName+".log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("logger: error while opening log file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	l := log.New(mw, suiteName+" ", log.LstdFlags)
	l.SetOutput(mw)
	return l, f, nil
}
