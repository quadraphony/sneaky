package config

import (
	"errors"
	"testing"
)

func TestDetectAndValidateSingboxConfig(t *testing.T) {
	input, err := Parse([]byte(`{
		"log": {"level": "info"},
		"outbounds": [
			{"type": "direct", "tag": "direct"}
		]
	}`), "inline")
	if err != nil {
		t.Fatalf("parse input: %v", err)
	}

	metadata, err := DetectAndValidate(input)
	if err != nil {
		t.Fatalf("detect and validate: %v", err)
	}

	if metadata.Format != FormatJSON {
		t.Fatalf("expected json format, got %q", metadata.Format)
	}
	if metadata.AdapterID != AdapterSingbox {
		t.Fatalf("expected singbox adapter, got %q", metadata.AdapterID)
	}
	if len(metadata.TopLevelKeys) != 2 || metadata.TopLevelKeys[0] != "log" || metadata.TopLevelKeys[1] != "outbounds" {
		t.Fatalf("unexpected top-level keys: %#v", metadata.TopLevelKeys)
	}
}

func TestDetectRejectsAmbiguousConfig(t *testing.T) {
	input, err := Parse([]byte(`{
		"log": {"level": "info"},
		"route": {}
	}`), "inline")
	if err != nil {
		t.Fatalf("parse input: %v", err)
	}

	_, err = Detect(input)
	if err == nil {
		t.Fatal("expected ambiguous format error")
	}

	var cfgErr *ValidationError
	if !errors.As(err, &cfgErr) {
		t.Fatalf("expected validation error, got %T", err)
	}
	if cfgErr.Code != ErrCodeAmbiguousFormat {
		t.Fatalf("expected ambiguous format code, got %q", cfgErr.Code)
	}
}

func TestDetectRejectsUnknownConfig(t *testing.T) {
	input, err := Parse([]byte(`{"name":"demo"}`), "inline")
	if err != nil {
		t.Fatalf("parse input: %v", err)
	}

	_, err = Detect(input)
	if err == nil {
		t.Fatal("expected unknown format error")
	}

	var cfgErr *ValidationError
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrCodeUnknownFormat {
		t.Fatalf("expected unknown format error, got %v", err)
	}
}

func TestDetectRejectsInvalidJSON(t *testing.T) {
	input, err := Parse([]byte(`{"outbounds":[}`), "inline")
	if err != nil {
		t.Fatalf("parse input: %v", err)
	}

	_, err = Detect(input)
	if err == nil {
		t.Fatal("expected invalid json error")
	}

	var cfgErr *ValidationError
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrCodeInvalidJSON {
		t.Fatalf("expected invalid json error, got %v", err)
	}
}

func TestValidateRejectsMissingOutboundType(t *testing.T) {
	input, err := Parse([]byte(`{
		"outbounds": [
			{"tag": "direct"}
		]
	}`), "inline")
	if err != nil {
		t.Fatalf("parse input: %v", err)
	}

	metadata, err := Detect(input)
	if err != nil {
		t.Fatalf("detect: %v", err)
	}

	err = Validate(input, metadata)
	if err == nil {
		t.Fatal("expected validation failure")
	}

	var cfgErr *ValidationError
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrCodeMissingRequiredField {
		t.Fatalf("expected missing field error, got %v", err)
	}
}

func TestParseRejectsEmptyInput(t *testing.T) {
	_, err := Parse([]byte(" \n\t "), "inline")
	if err == nil {
		t.Fatal("expected empty input error")
	}

	var cfgErr *ValidationError
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrCodeEmptyInput {
		t.Fatalf("expected empty input error, got %v", err)
	}
}

func TestDetectAndValidateSSHConfig(t *testing.T) {
	input, err := Parse([]byte(`{
		"ssh_tunnel": {
			"host": "example.com",
			"user": "demo",
			"local_socks_port": 1080
		}
	}`), "inline")
	if err != nil {
		t.Fatalf("parse input: %v", err)
	}

	metadata, err := DetectAndValidate(input)
	if err != nil {
		t.Fatalf("detect and validate ssh config: %v", err)
	}
	if metadata.AdapterID != AdapterSSH {
		t.Fatalf("expected ssh adapter, got %q", metadata.AdapterID)
	}
}

func TestDetectAndValidateSSHConfigRejectsNegativeKeepaliveOptions(t *testing.T) {
	input, err := Parse([]byte(`{
		"ssh_tunnel": {
			"host": "example.com",
			"user": "demo",
			"local_socks_port": 1080,
			"server_alive_interval_seconds": -1
		}
	}`), "inline")
	if err != nil {
		t.Fatalf("parse input: %v", err)
	}

	_, err = DetectAndValidate(input)
	if err == nil {
		t.Fatal("expected validation failure")
	}

	var cfgErr *ValidationError
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrCodeMissingRequiredField {
		t.Fatalf("expected missing required field error, got %v", err)
	}
}
