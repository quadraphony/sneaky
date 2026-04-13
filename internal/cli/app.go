package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/net/proxy"
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
	case "probe":
		return a.runProbe(args[1:])
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
		AdapterID:  sneaky.AdapterID(metadata.AdapterID),
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
	fmt.Fprintln(a.Stdout, "commands: validate, probe, start, stop, status, version")
}

func inspectConfig(path string) (config.Metadata, error) {
	input, err := config.LoadFile(path)
	if err != nil {
		return config.Metadata{}, err
	}

	return config.DetectAndValidate(input)
}

func (a *App) runProbe(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(a.Stderr, "usage: sneakycli probe <config-path> <url>")
		return 2
	}

	configPath := args[0]
	targetURL := args[1]

	metadata, err := inspectConfig(configPath)
	if err != nil {
		fmt.Fprintf(a.Stderr, "probe failed: %v\n", err)
		return 1
	}

	startReq, proxyPort, cleanup, err := probeRequest(configPath, metadata)
	if err != nil {
		fmt.Fprintf(a.Stderr, "probe failed: %v\n", err)
		return 1
	}
	defer cleanup()

	manager := sneaky.New()
	if err := manager.Start(context.Background(), startReq); err != nil {
		fmt.Fprintf(a.Stderr, "probe failed: %v\n", err)
		return 1
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = manager.Stop(stopCtx)
	}()

	if err := waitForProxyReady(proxyPort, 60*time.Second); err != nil {
		fmt.Fprintf(a.Stderr, "probe failed: %v\n", err)
		return 1
	}

	statusCode, bodyBytes, err := probeHTTP(targetURL, proxyPort)
	if err != nil {
		fmt.Fprintf(a.Stderr, "probe failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.Stdout, "probe ok adapter=%s proxy_port=%d status=%d bytes=%d\n", metadata.AdapterID, proxyPort, statusCode, bodyBytes)
	return 0
}

func probeRequest(configPath string, metadata config.Metadata) (sneaky.StartRequest, int, func(), error) {
	switch metadata.AdapterID {
	case config.AdapterSingbox:
		raw, port, err := prepareSingboxProbeConfig(configPath)
		if err != nil {
			return sneaky.StartRequest{}, 0, noop, err
		}
		return sneaky.StartRequest{
			AdapterID: sneaky.AdapterID(metadata.AdapterID),
			RawConfig: raw,
		}, port, noop, nil
	case config.AdapterSSH:
		port, err := sshProbePort(configPath)
		if err != nil {
			return sneaky.StartRequest{}, 0, noop, err
		}
		return sneaky.StartRequest{
			AdapterID:  sneaky.AdapterID(metadata.AdapterID),
			ConfigPath: configPath,
		}, port, noop, nil
	default:
		return sneaky.StartRequest{}, 0, noop, fmt.Errorf("probe is not implemented for adapter %q", metadata.AdapterID)
	}
}

func probeHTTP(targetURL string, proxyPort int) (int, int, error) {
	proxyAddr := net.JoinHostPort("127.0.0.1", strconv.Itoa(proxyPort))
	socksDialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return 0, 0, fmt.Errorf("create SOCKS5 dialer: %w", err)
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		type contextDialer interface {
			DialContext(context.Context, string, string) (net.Conn, error)
		}
		if d, ok := socksDialer.(contextDialer); ok {
			return d.DialContext(ctx, network, addr)
		}
		return socksDialer.Dial(network, addr)
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			DialContext:         dialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return 0, 0, fmt.Errorf("perform HTTP request through proxy: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("read HTTP response body: %w", err)
	}

	return resp.StatusCode, len(body), nil
}

func waitForProxyReady(proxyPort int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(proxyPort))

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("proxy listener on %s was not ready within %s", address, timeout)
}

func prepareSingboxProbeConfig(path string) ([]byte, int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("read sing-box config: %w", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, 0, fmt.Errorf("decode sing-box config: %w", err)
	}

	outbounds, ok := doc["outbounds"].([]any)
	if !ok || len(outbounds) == 0 {
		return nil, 0, fmt.Errorf("sing-box probe requires at least one outbound")
	}

	firstOutbound, ok := outbounds[0].(map[string]any)
	if !ok {
		return nil, 0, fmt.Errorf("sing-box probe requires object outbounds")
	}

	tag, _ := firstOutbound["tag"].(string)
	if tag == "" {
		tag = "probe-out"
		firstOutbound["tag"] = tag
	}

	port, err := reserveLocalPort()
	if err != nil {
		return nil, 0, err
	}

	doc["inbounds"] = []map[string]any{
		{
			"type":        "socks",
			"tag":         "probe-in",
			"listen":      "127.0.0.1",
			"listen_port": port,
		},
	}

	route, _ := doc["route"].(map[string]any)
	if route == nil {
		route = map[string]any{}
	}
	if routeTag, _ := route["final"].(string); routeTag == "" {
		route["final"] = tag
	}
	doc["route"] = route

	updated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, 0, fmt.Errorf("encode probe config: %w", err)
	}
	return updated, port, nil
}

func sshProbePort(path string) (int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read ssh config: %w", err)
	}

	var doc struct {
		SSHTunnel struct {
			LocalSOCKSPort int `json:"local_socks_port"`
		} `json:"ssh_tunnel"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return 0, fmt.Errorf("decode ssh config: %w", err)
	}
	if doc.SSHTunnel.LocalSOCKSPort <= 0 {
		return 0, fmt.Errorf("ssh_tunnel.local_socks_port must be greater than zero")
	}
	return doc.SSHTunnel.LocalSOCKSPort, nil
}

func reserveLocalPort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("reserve local port: %w", err)
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("resolve reserved port")
	}
	return addr.Port, nil
}

func noop() {}
