package singbox

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"sneaky-core/internal/adapter"
	"sneaky-core/internal/config"
	"sneaky-core/internal/runtime"
)

type Adapter struct {
	binaryPath string
}

func New(binaryPath string) *Adapter {
	return &Adapter{binaryPath: binaryPath}
}

func (a *Adapter) Identity() string {
	return config.AdapterSingbox
}

func (a *Adapter) Capabilities() adapter.Capabilities {
	return adapter.Capabilities{}
}

func (a *Adapter) ValidateConfig(req adapter.StartRequest) error {
	bin, err := a.resolveBinary()
	if err != nil {
		return err
	}

	resolvedPath, cleanup, err := a.prepareConfig(req)
	if err != nil {
		return err
	}
	defer cleanup()

	return a.checkConfig(context.Background(), bin, resolvedPath)
}

func (a *Adapter) Start(ctx context.Context, req adapter.StartRequest) (runtime.Handle, error) {
	bin, err := a.resolveBinary()
	if err != nil {
		return nil, err
	}

	resolvedPath, cleanup, err := a.prepareConfig(req)
	if err != nil {
		return nil, err
	}

	if err := a.checkConfig(ctx, bin, resolvedPath); err != nil {
		cleanup()
		return nil, err
	}

	// Try to extract an inbound port for readiness checking.
	inboundPort := a.detectInboundPort(req)

	ready := func() bool {
		if inboundPort <= 0 {
			return true
		}
		address := net.JoinHostPort("127.0.0.1", strconv.Itoa(inboundPort))
		conn, err := net.DialTimeout("tcp", address, 50*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}

	cmd := exec.CommandContext(ctx, bin, "run", "-c", resolvedPath, "--disable-color")
	handle, err := runtime.StartProcessWithReadiness(cmd, cleanup, ready)
	if err != nil {
		cleanup()
		return nil, err
	}

	return handle, nil
}

type partialInbound struct {
	ListenPort int `json:"listen_port"`
}

type partialConfig struct {
	Inbounds []partialInbound `json:"inbounds"`
}

func (a *Adapter) detectInboundPort(req adapter.StartRequest) int {
	var raw []byte
	if len(req.RawConfig) > 0 {
		raw = req.RawConfig
	} else if req.ConfigPath != "" {
		var err error
		raw, err = os.ReadFile(req.ConfigPath)
		if err != nil {
			return 0
		}
	}

	var pc partialConfig
	if err := json.Unmarshal(raw, &pc); err != nil {
		return 0
	}

	if len(pc.Inbounds) > 0 && pc.Inbounds[0].ListenPort > 0 {
		return pc.Inbounds[0].ListenPort
	}

	return 0
}
