package main

import (
	"errors"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func buildLogger(address string, tag string) *syslog.Writer {
	logger, err := syslog.Dial(
		"udp",
		address,
		syslog.LOG_DEBUG,
		tag,
	)
	failOnErr(err)

	return logger
}

func buildCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func failOnErr(err error) {
	if err == nil {
		return

	}
	log.Fatal(err)
	os.Exit(1)
}



func main() {
	token := os.Getenv("SYSLOG_TAG")
	address := os.Getenv("SYSLOG_ADDRESS")
	arguments := os.Args[1:]
	if len(arguments) != 1 {
		failOnErr(errors.New("no GLOB pattern argument passed for the watcher"))
	}
	if len(token) == 0 || len(address) == 0 {
		failOnErr(errors.New("no SYSLOG_ADDRESS or SYSLOG_TAG environment variable provided"))
	}

	logger := buildLogger(address, token)
	defer logger.Close()

	forwarder := FilesToSyslogForwarder{
		logWriter:    logger,
		watchedFiles: make(map[string]bool),
	}

	buildCloseHandler()

	// As some files are created on the fly by the application we'll register an interval to register new files
	for _ = range time.Tick(time.Minute) {
		forwarder.FindAndForwardFiles(arguments[0])
	}
}


