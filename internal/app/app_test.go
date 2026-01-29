package app

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	a, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	if a == nil {
		t.Fatal("NewApp() returned nil")
	}

	if a.name != "home-ctrl" {
		t.Errorf("Expected name 'home-ctrl', got '%s'", a.name)
	}

	if a.version != "0.1.0" {
		t.Errorf("Expected version '0.1.0', got '%s'", a.version)
	}

	if a.config == nil {
		t.Fatal("Config should not be nil")
	}

	if a.server == nil {
		t.Fatal("Server should not be nil")
	}
}

func TestAppRun(t *testing.T) {
	a, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	// Note: We can't actually run the server in tests as it blocks
	// This test just verifies that NewApp() works correctly
	if a.server == nil {
		t.Fatal("Server should not be nil")
	}
}