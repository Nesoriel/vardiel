package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunVersionReportsVardiel(t *testing.T) {
	var stdout bytes.Buffer
	if err := run(context.Background(), []string{"version"}, &stdout, io.Discard); err != nil {
		t.Fatalf("run version: %v", err)
	}
	var result struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode version output: %v", err)
	}
	if result.Name != "vardiel" || result.Version != version {
		t.Fatalf("unexpected version output: %#v", result)
	}
}

func TestCLIUsageNamesVardiel(t *testing.T) {
	for _, print := range []func(io.Writer){printUsage, printToolUsage, printAgentUsage} {
		var output bytes.Buffer
		print(&output)
		if !strings.HasPrefix(output.String(), "usage: vardiel ") {
			t.Fatalf("unexpected CLI usage: %q", output.String())
		}
	}
}

func TestDeploymentManifestUsesVardielIdentity(t *testing.T) {
	payload, err := os.ReadFile("../../deploy/kubernetes/vardiel-readonly-rbac.yaml")
	if err != nil {
		t.Fatalf("read deployment manifest: %v", err)
	}
	manifest := string(payload)
	for _, expected := range []string{"name: vardiel\n", "namespace: vardiel\n", "name: vardiel-readonly-diagnostics\n"} {
		if !strings.Contains(manifest, expected) {
			t.Fatalf("deployment manifest missing %q", expected)
		}
	}
}

func TestRunAgentRequiresPrompt(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := run(context.Background(), []string{"agent", "run"}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "prompt") {
		t.Fatalf("expected missing prompt error, got %v", err)
	}
}

func TestRunAgentRejectsMissingArkConfiguration(t *testing.T) {
	for _, key := range []string{
		"ARK_API_KEY",
		"ARK_ACCESS_KEY",
		"ARK_SECRET_KEY",
		"ARK_MODEL_ID",
		"ARK_BASE_URL",
		"ARK_REGION",
		"ARK_TIMEOUT",
		"ARK_RETRY_TIMES",
		"ARK_THINKING",
		"ARK_MAX_TOOL_CALLS",
		"ARK_PARALLEL_TOOL_CALLS",
	} {
		t.Setenv(key, "")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := run(context.Background(), []string{"agent", "run", "inspect", "example.com"}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "ARK_MODEL_ID") {
		t.Fatalf("expected missing Ark configuration error, got %v", err)
	}
}
