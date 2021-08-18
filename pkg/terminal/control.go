package terminal

import (
	"os"

	"github.com/docker/docker/pkg/term"
)

// isTty returns true if both stdin and stdout are a TTY.
func IsTTY() bool {
	_, isStdinTerminal := term.GetFdInfo(os.Stdin)
	_, isStdoutTerminal := term.GetFdInfo(os.Stdout)
	return isStdinTerminal && isStdoutTerminal
}

// setRawInput sets the stream terminal in raw mode, so process captures
// Ctrl+C and other commands to forward to remote process.
// It returns a cleanup function that restores terminal to original mode.
func SetRawInput(stream interface{}) (cleanup func(), err error) {
	fd, isTerminal := term.GetFdInfo(stream)
	if !isTerminal {
		return nil, err
	}

	state, err := term.SetRawTerminal(fd)
	if err != nil {
		return nil, err
	}

	return func() { term.RestoreTerminal(fd, state) }, nil
}

// setRawOutput sets the output stream in Windows to raw mode,
// so it disables LF -> CRLF translation.
// It's basically a no-op on unix.
func SetRawOutput(stream interface{}) (cleanup func(), err error) {
	fd, isTerminal := term.GetFdInfo(stream)
	if !isTerminal {
		return nil, err
	}

	state, err := term.SetRawTerminalOutput(fd)
	if err != nil {
		return nil, err
	}

	return func() { term.RestoreTerminal(fd, state) }, nil
}
