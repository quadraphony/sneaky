package sneaky

import "testing"

func TestInspectConfigBytes(t *testing.T) {
	metadata, err := InspectConfigBytes([]byte(`{
		"ssh_tunnel": {
			"host": "example.com",
			"user": "demo",
			"local_socks_port": 1080
		}
	}`))
	if err != nil {
		t.Fatalf("inspect config bytes: %v", err)
	}
	if metadata.AdapterID != "ssh" {
		t.Fatalf("expected ssh adapter, got %q", metadata.AdapterID)
	}
	if metadata.Format != "json" {
		t.Fatalf("expected json format, got %q", metadata.Format)
	}
}
