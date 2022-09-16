package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printMessage(t *testing.T) {
	msg = "Hello"

	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printMessage()
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, msg) {
		t.Errorf("Expected to find %q in console output: %s", msg, output)
	}
}

func Test_updateMessage(t *testing.T) {
	msg = "Hello"
	newMsg := "Goodbye"

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go updateMessage(newMsg, wg)
	wg.Wait()

	if strings.Compare(msg, newMsg) != 0 {
		t.Errorf("Expected %q to be equal to %q", msg, newMsg)
	}
}

func Test_main(t *testing.T) {
	initialMsg := "Mars"
	updatedMessages := []string{
		"Hello, universe!",
		"Hello, cosmos!",
		"Hello, world!",
	}
	msg = initialMsg

	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if strings.Contains(output, initialMsg) {
		t.Errorf("Expected not to find %q in console output: %s", initialMsg, output)
	}

	for _, message := range updatedMessages {
		if !strings.Contains(output, message) {
			t.Errorf("Expected to find %q in console output: %s", message, output)
		}
	}
}
