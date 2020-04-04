package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"log/syslog"
	"path/filepath"
)

type FilesToSyslogForwarder struct {
	logWriter    *syslog.Writer
	watchedFiles map[string]bool
}

func (f FilesToSyslogForwarder) FindAndForwardFiles(pattern string) {
	matches, err := filepath.Glob(pattern)
	failOnErr(err)

	for _, match := range matches {
		if f.watchedFiles[match] {
			continue
		}

		f.registerForwarder(match)
	}
}

func (f FilesToSyslogForwarder) registerForwarder(match string) {
	go func(match string, writer *syslog.Writer) {
		t, err := tail.TailFile(match, tail.Config{Follow: true, ReOpen: true})
		failOnErr(err)
		for line := range t.Lines {
			writer.Info(
				fmt.Sprintf(
					`@cee:{"message":"%s", "customer-id": "%s"}`,
					line.Text,
					"53b574d5-92a1-44e4-98ef-c2f891607d0a",
				),
			)
		}
	}(match, f.logWriter)

	f.watchedFiles[match] = true
	log.Printf(`Started forwarding file contents of  %s`, match)
}
