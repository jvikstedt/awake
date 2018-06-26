package cron_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jvikstedt/awake/cron"
)

var logs = &bytes.Buffer{}

func TestAddEntry(t *testing.T) {
	logger := log.New(logs, "", log.LstdFlags)

	scheduler := cron.New(logger)
	go scheduler.Start()
	defer scheduler.Stop()

	callCh := make(chan bool, 3)
	callCount := 0

	callback := func(id cron.EntryID) {
		if id != 1 {
			t.Errorf(fmt.Sprintf("Expected id of %d but got %d", 1, id))
		}
		callCh <- true
	}

	scheduler.AddEntry(1, "@every 1s", callback)

	timeout := time.After(2500 * time.Millisecond)

Loop:
	for {
		select {
		case <-timeout:
			t.Fatalf("timeout, callback did not get called enough times")
		case <-callCh:
			callCount++
			if callCount >= 2 {
				break Loop
			}
		}
	}
}
