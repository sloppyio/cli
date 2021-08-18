package command

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sloppyio/cli/pkg/api"
	"github.com/sloppyio/cli/pkg/terminal"
	"github.com/sloppyio/cli/ui"
)

const (
	origin = "https://localhost/"
)

type ConsoleCommand struct {
	Client *api.Client
	UI     ui.UI

	Stdin  io.Reader
	Stdout io.WriteCloser
	Stderr io.WriteCloser
}

func (c *ConsoleCommand) Help() string {
	helpText := `
Usage: sloppy console [OPTIONS] (PROJECT/SERVICE/APP) (COMMAND)

  Attach to the console session of an application.

  If no command is specified an interactive session will be attempted,
  this uses the -i -t flags and requires a working TTY.

Options:

  -i  attach stdin to the container
  -t  allocate a pseudo-tty
  -e  sets the escape character
`
	return strings.TrimSpace(helpText)
}

func (c *ConsoleCommand) Synopsis() string {
	return "Launch the console of an application"
}

func (c *ConsoleCommand) Run(args []string) int {
	var stdinOpt, ttyOpt bool
	var escapeChar, appPath string

	cmdFlags := newFlagSet("console", flag.ContinueOnError)
	cmdFlags.BoolVar(&stdinOpt, "i", true, "")
	cmdFlags.BoolVar(&ttyOpt, "t", terminal.IsTTY(), "")
	cmdFlags.StringVar(&escapeChar, "e", "~", "")

	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		c.UI.Output("See 'sloppy change --help'.")
		return 1
	}

	args = cmdFlags.Args()

	if len(args) == 0 {
		return c.UI.ErrorNotEnoughArgs("console", "", 1)
	}

	appPath = args[0]

	if !(strings.Count(strings.Trim(appPath, "/"), "/") == 2) {
		return c.UI.ErrorInvalidAppPath(args[0])
	}

	if ttyOpt && !stdinOpt {
		c.UI.Error("-i must be enabled if running with tty")
		return 1
	}

	if !stdinOpt {
		c.Stdin = bytes.NewReader(nil)
	}

	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}

	if c.Stdout == nil {
		c.Stdout = os.Stdout
	}

	if c.Stderr == nil {
		c.Stderr = os.Stderr
	}

	code, err := c.consoleImpl(appPath, args[1:], ttyOpt, escapeChar, c.Stdin, c.Stdout, c.Stderr)
	if err != nil {
		return 1
	}

	return code
}

func (c *ConsoleCommand) consoleImpl(app string, command []string, tty bool, escapeChar string, stdin io.Reader, stdout, stderr io.WriteCloser) (int, error) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	if tty {
		if stdin == nil {
			return 1, fmt.Errorf("stdin is required with TTY")
		}

		stdinRestore, err := terminal.SetRawInput(stdin)
		if err != nil {
			return 1, err
		}
		defer stdinRestore()

		stdoutRestore, err := terminal.SetRawOutput(stdout)
		if err != nil {
			return 1, err
		}
		defer stdoutRestore()

		if escapeChar != "" {
			stdin = terminal.NewReader(stdin, escapeChar[0], func(b byte) bool {
				switch b {
				case '.':
					stdoutRestore()
					stdinRestore()

					stderr.Write([]byte("\nClosed!\n"))
					cancelFn()
					return true
				default:
					return false
				}
			})
		}
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range signalCh {
			cancelFn()
		}
	}()

	exec := &consoleExec{
		client: c.Client,
		app:    app,

		tty:     tty,
		command: command,

		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	return exec.run(ctx)
}

type consoleExec struct {
	client *api.Client
	app    string

	tty     bool
	command []string

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *consoleExec) run(ctx context.Context) (int, error) {
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	ws, err := c.initConnection()
	if err != nil {
		return 1, err
	}
	defer ws.Close()

	sendErrCh := c.setupSend(ctx, ws)
	recvErrCh := c.setupReceive(ctx, ws)

	for {
		select {
		case <-ctx.Done():
			return 1, err
		case sendErr := <-sendErrCh:
			return 1, sendErr
		case recvErr := <-recvErrCh:
			return 1, recvErr
		}
	}
}

func (c *consoleExec) initConnection() (*websocket.Conn, error) {
	url := fmt.Sprintf("%sconsole?token=%s&app=%s",
		strings.Replace(c.client.GetBaseURL(), "https://", "wss://", 1),
		strings.TrimPrefix(c.client.GetHeader("Authorization")[0], "Bearer "),
		c.app,
	)
	headers := http.Header{
		"Origin": []string{origin},
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *consoleExec) setupSend(ctx context.Context, conn *websocket.Conn) <-chan error {
	var sendLock sync.Mutex

	errCh := make(chan error, 4)
	send := func(v []byte) {
		sendLock.Lock()
		defer sendLock.Unlock()

		conn.WriteMessage(websocket.TextMessage, v)
	}

	// process stdin
	go func() {
		bytesIn := make([]byte, 2048)

		for {
			if ctx.Err() != nil {
				return
			}

			n, err := c.stdin.Read(bytesIn)

			if n != 0 {
				send(bytesIn[:n])
			}

			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// send a heartbeat every 10 seconds
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Second):
				send(nil)
			}
		}
	}()

	return errCh
}

func (c *consoleExec) setupReceive(ctx context.Context, conn *websocket.Conn) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		for ctx.Err() == nil {
			_, d, err := conn.ReadMessage()
			// check if the error is due to a closed connection
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				errCh <- fmt.Errorf("websocket closed before receiving exit code: %w", err)
				return
			} else if err != nil {
				errCh <- err
				return
			}

			if len(d) != 0 {
				c.stdout.Write(d)
			}
		}
	}()

	return errCh
}
