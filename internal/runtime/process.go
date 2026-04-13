package runtime

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type ProcessHandle struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	cleanup func()
	done    chan error
	once    sync.Once
	state   State
	waitErr error
}

func StartProcess(cmd *exec.Cmd, cleanup func()) (*ProcessHandle, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	handle := &ProcessHandle{
		cmd:     cmd,
		cleanup: cleanup,
		done:    make(chan error, 1),
		state:   StateStarting,
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		err := cmd.Wait()
		handle.mu.Lock()
		handle.waitErr = err
		handle.state = StateStopped
		handle.mu.Unlock()
		handle.finalize()
		handle.done <- err
	}()

	select {
	case err := <-handle.done:
		if err != nil {
			detail := bytes.TrimSpace(append(stdout.Bytes(), stderr.Bytes()...))
			if len(detail) > 0 {
				return nil, fmt.Errorf("process exited during startup: %s", detail)
			}
			return nil, fmt.Errorf("process exited during startup: %w", err)
		}
		return nil, fmt.Errorf("process exited during startup without error")
	case <-time.After(250 * time.Millisecond):
		handle.mu.Lock()
		handle.state = StateRunning
		handle.mu.Unlock()
		return handle, nil
	}
}

func (h *ProcessHandle) Stop(ctx context.Context) error {
	h.mu.Lock()
	if h.state == StateStopped {
		err := h.waitErr
		h.mu.Unlock()
		return err
	}
	h.state = StateStopping
	process := h.cmd.Process
	h.mu.Unlock()

	if process == nil {
		return fmt.Errorf("process is not available")
	}

	if err := process.Signal(syscall.SIGTERM); err != nil && err != os.ErrProcessDone {
		return fmt.Errorf("signal process: %w", err)
	}

	select {
	case err := <-h.done:
		if err != nil && !IsExpectedExit(err) {
			return fmt.Errorf("wait for process shutdown: %w", err)
		}
		return nil
	case <-ctx.Done():
		_ = process.Kill()
		<-h.done
		return ctx.Err()
	}
}

func (h *ProcessHandle) State() State {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.state
}

func (h *ProcessHandle) finalize() {
	h.once.Do(func() {
		if h.cleanup != nil {
			h.cleanup()
		}
	})
}

func IsExpectedExit(err error) bool {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	return exitErr.ExitCode() == -1 || exitErr.ExitCode() == 143
}
