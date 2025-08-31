package hashattr_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/isometry/bitrat/hashattr"
	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
)

func TestHashAttr_GetSetRemove(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	attr := hashattr.New("user.bitrat.test")
	testData := []byte("test hash value")

	// Test Set
	err := attr.Set(testFile, testData)
	if err != nil {
		t.Skip("Extended attributes not supported on this filesystem")
	}

	// Test Get
	retrievedData := attr.Get(testFile)
	if retrievedData == nil {
		t.Fatal("Failed to get xattr: returned nil")
	}

	if string(retrievedData) != string(testData) {
		t.Errorf("Retrieved data doesn't match: got %q, want %q", retrievedData, testData)
	}

	// Test Remove
	err = attr.Remove(testFile)
	if err != nil {
		t.Fatalf("Failed to remove xattr: %v", err)
	}

	// Verify removal
	retrievedData = attr.Get(testFile)
	if retrievedData != nil {
		t.Error("Expected nil after removal, but got data")
	}
}

func TestHashAttr_Get_NonExistent(t *testing.T) {
	attr := hashattr.New("user.bitrat.test")

	// Test with non-existent file
	data := attr.Get("/non/existent/file")
	if data != nil {
		t.Error("Expected nil for non-existent file")
	}

	// Test with file that exists but has no xattr
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(tmpFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	data = attr.Get(tmpFile)
	if data != nil {
		t.Error("Expected nil for file without xattr")
	}
}

func TestHashAttr_Reader(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	attr := hashattr.New("user.bitrat.test")
	testHash := []byte("test hash")

	// Set xattr
	if err := attr.Set(testFile, testHash); err != nil {
		t.Skip("Extended attributes not supported on this filesystem")
	}

	// Create channels
	input := make(chan *pathwalk.File, 1)
	output := make(chan *hasher.FileHash, 1)
	var wg sync.WaitGroup

	// Start reader
	wg.Go(attr.Reader(input, output))

	// Send file
	input <- &pathwalk.File{
		Path:    testFile,
		Size:    12,
		ModTime: time.Now(),
	}
	close(input)

	// Wait and check output
	result := <-output
	wg.Wait()

	if result.File.Path != testFile {
		t.Errorf("Wrong file path: got %q, want %q", result.File.Path, testFile)
	}

	if string(result.Hash) != string(testHash) {
		t.Errorf("Wrong hash: got %q, want %q", result.Hash, testHash)
	}
}

func TestHashAttr_Writer_EmptyChannel(t *testing.T) {
	// Test that Writer handles empty channel without crashing
	attr := hashattr.New("user.bitrat.test")

	// Create empty channel
	input := make(chan *hasher.FileHash)
	var wg sync.WaitGroup

	// Capture stderr to avoid test output pollution
	oldStderr := os.Stderr
	os.Stderr, _ = os.Open(os.DevNull)
	defer func() { os.Stderr = oldStderr }()

	// Close channel immediately
	close(input)

	// Start writer with empty channel
	wg.Go(attr.Writer(input))
	wg.Wait()

	// If we get here without panic, the test passes
}

func TestHashAttr_Remover(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	attr := hashattr.New("user.bitrat.test")

	// Set an xattr first
	if err := attr.Set(testFile, []byte("test")); err != nil {
		t.Skip("Extended attributes not supported on this filesystem")
	}

	// Create channel
	input := make(chan *pathwalk.File, 1)
	var wg sync.WaitGroup

	// Capture stderr to avoid test output pollution
	oldStderr := os.Stderr
	os.Stderr, _ = os.Open(os.DevNull)
	defer func() { os.Stderr = oldStderr }()

	// Start remover
	wg.Go(attr.Remover(input))

	// Send file
	input <- &pathwalk.File{Path: testFile}
	close(input)

	wg.Wait()

	// Verify removal
	if data := attr.Get(testFile); data != nil {
		t.Error("Expected xattr to be removed")
	}
}
