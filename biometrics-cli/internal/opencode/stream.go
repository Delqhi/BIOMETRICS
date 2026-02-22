package opencode

import (
	"bufio"
	"io"
	"os/exec"
)

// streamOutput reads stdout and stderr from cmd and sends lines to a channel.
// It launches an async goroutine that uses bufio.Scanner to read lines from
// both stdout and stderr, sending each line to the returned channel.
func (e *Executor) streamOutput(cmd *exec.Cmd) <-chan string {
	outChan := make(chan string)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		e.logger.Error("Failed to create stdout pipe", "error", err)
		close(outChan)
		return outChan
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		e.logger.Error("Failed to create stderr pipe", "error", err)
		close(outChan)
		return outChan
	}

	// Launch async reading goroutine
	go func() {
		defer close(outChan)

		// Combine stdout and stderr into single reader
		multiReader := io.MultiReader(stdoutPipe, stderrPipe)
		scanner := bufio.NewScanner(multiReader)

		// Read lines and send to channel
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case outChan <- line:
				// Line sent successfully
			default:
				// Channel full or closed, stop sending
				return
			}
		}

		// Handle scanner errors
		if err := scanner.Err(); err != nil {
			e.logger.Error("Scanner error while reading output", "error", err)
		}
	}()

	return outChan
}
