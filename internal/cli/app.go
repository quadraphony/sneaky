package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sneaky-core/internal/config"
	"sneaky-core/pkg/sneaky"
)

const version = "dev"

type App struct {
	Stdout io.Writer
	Stderr io.Writer
}

func New(stdout, stderr io.Writer) *App {
	return &App{
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (a *App) Run(args []string) int {
	if len(args) == 0 {
		a.printUsage()
		return 2
	}

	switch args[0] {
	case "validate":
		return a.runValidate(args[1:])
	case "start":
		return a.runStart(args[1:])
	case "stop":
		return a.runStop(args[1:])
	case "status":
		return a.runStatus(args[1:])
	case "version":
		fmt.Fprintln(a.Stdout, version)
		return 0
	case "help", "-h", "--help":
		a.printUsage()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "unknown command %q\n", args[0])
		a.printUsage()
		return 2
	}
}

func (a *App) runValidate(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(a.Stderr, "usage: sneakycli validate <config-path>")
		return 2
	}

	metadata, err := inspectConfig(args[0])
	if err != nil {
		fmt.Fprintf(a.Stderr, "validation failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.Stdout, "valid adapter=%s format=%s\n", metadata.AdapterID, metadata.Format)
	return 0
}

func (a *App) runStart(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(a.Stderr, "usage: sneakycli start <config-path>")
		return 2
	}

	metadata, err := inspectConfig(args[0])
	if err != nil {
		fmt.Fprintf(a.Stderr, "start failed: %v\n", err)
		return 1
	}

	manager := sneaky.New()
	if err := manager.Start(context.Background(), sneaky.StartRequest{
		AdapterID:  metadata.AdapterID,
		ConfigPath: args[0],
	}); err != nil {
		fmt.Fprintf(a.Stderr, "start failed: %v\n", err)
		return 1
	}

	state := stateFile{
		PID:        os.Getpid(),
		AdapterID:  metadata.AdapterID,
		ConfigPath: args[0],
		StartedAt:  time.Now().UTC(),
	}
	if err := writeStateFile(state); err != nil {
		fmt.Fprintf(a.Stderr, "start warning: unable to write state file: %v\n", err)
	}
	defer removeStateFile()

	fmt.Fprintf(a.Stdout, "running adapter=%s pid=%d\n", metadata.AdapterID, os.Getpid())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	<-signals

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := manager.Stop(stopCtx); err != nil {
		fmt.Fprintf(a.Stderr, "stop failed: %v\n", err)
		return 1
	}

	fmt.Fprintln(a.Stdout, "stopped")
	return 0
}

func (a *App) runStop(args []string) int {
	if len(args) != 0 {
		fmt.Fprintln(a.Stderr, "usage: sneakycli stop")
		return 2
	}

	state, err := readStateFile()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(a.Stdout, "stopped")
			return 0
		}
		fmt.Fprintf(a.Stderr, "stop failed: %v\n", err)
		return 1
	}

	process, err := os.FindProcess(state.PID)
	if err != nil {
		fmt.Fprintf(a.Stderr, "stop failed: %v\n", err)
		return 1
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			_ = removeStateFile()
			fmt.Fprintln(a.Stdout, "stopped")
			return 0
		}
		fmt.Fprintf(a.Stderr, "stop failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.Stdout, "stopping pid=%d adapter=%s\n", state.PID, state.AdapterID)
	return 0
}

func (a *App) runStatus(args []string) int {
	if len(args) != 0 {
		fmt.Fprintln(a.Stderr, "usage: sneakycli status")
		return 2
	}

	state, err := readStateFile()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(a.Stdout, "state=stopped")
			return 0
		}
		fmt.Fprintf(a.Stderr, "status failed: %v\n", err)
		return 1
	}

	process, err := os.FindProcess(state.PID)
	if err == nil {
		err = process.Signal(syscall.Signal(0))
	}
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH) {
			_ = removeStateFile()
			fmt.Fprintln(a.Stdout, "state=stopped")
			return 0
		}
	}

	fmt.Fprintf(
		a.Stdout,
		"state=running pid=%d adapter=%s config=%s started_at=%s\n",
		state.PID,
		state.AdapterID,
		state.ConfigPath,
		state.StartedAt.UTC().Format(time.RFC3339),
	)
	return 0
}

func (a *App) printUsage() {
	fmt.Fprintln(a.Stdout, "usage: sneakycli <command> [args]")
	fmt.Fprintln(a.Stdout, "commands: validate, start, stop, status, version")
}

func inspectConfig(path string) (config.Metadata, error) {
	input, err := config.LoadFile(path)
	if err != nil {
		return config.Metadata{}, err
	}

	return config.DetectAndValidate(input)
}
