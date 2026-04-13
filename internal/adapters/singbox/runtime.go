package singbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"sneaky-core/internal/runtime"
)

type processHandle struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	cleanup func()
	done    chan error
	once    sync.Once
	state   runtime.State
	waitErr error
}

func startProcess(cmd *exec.Cmd, cleanup func()) (*processHandle, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	handle := &processHandle{
		cmd:     cmd,
		cleanup: cleanup,
		done:    make(chan error, 1),
		state:   runtime.StateStarting,
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		err := cmd.Wait()
		handle.mu.Lock()
		handle.waitErr = err
		handle.state = runtime.StateStopped
		handle.mu.Unlock()
		handle.finalize()
		handle.done <- err
	}()

	select {
	case err := <-handle.done:
		if err != nil {
			detail := bytes.TrimSpace(append(stdout.Bytes(), stderr.Bytes()...))
			if len(detail) > 0 {
				return nil, fmt.Errorf("sing-box exited during startup: %s", detail)
			}
			return nil, fmt.Errorf("sing-box exited during startup: %w", err)
		}
		return nil, fmt.Errorf("sing-box exited during startup without error")
	case <-time.After(250 * time.Millisecond):
		handle.mu.Lock()
		handle.state = runtime.StateRunning
		handle.mu.Unlock()
		return handle, nil
	}
}

func (h *processHandle) Stop(ctx context.Context) error {
	h.mu.Lock()
	if h.state == runtime.StateStopped {
		err := h.waitErr
		h.mu.Unlock()
		return err
	}
	h.state = runtime.StateStopping
	process := h.cmd.Process
	h.mu.Unlock()

	if process == nil {
		return fmt.Errorf("sing-box process is not available")
	}

	if err := process.Signal(syscall.SIGTERM); err != nil && err != os.ErrProcessDone {
		return fmt.Errorf("signal sing-box process: %w", err)
	}

	select {
	case err := <-h.done:
		if err != nil && !isExpectedExit(err) {
			return fmt.Errorf("wait for sing-box shutdown: %w", err)
		}
		return nil
	case <-ctx.Done():
		_ = process.Kill()
		<-h.done
		return ctx.Err()
	}
}

func (h *processHandle) State() runtime.State {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.state
}

func (h *processHandle) finalize() {
	h.once.Do(func() {
		if h.cleanup != nil {
			h.cleanup()
		}
	})
}

func isExpectedExit(err error) bool {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	return exitErr.ExitCode() == -1 || exitErr.ExitCode() == 143
}
