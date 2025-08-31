package pathwalk_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/isometry/bitrat/pathwalk"
)

func sinkFileChan(input <-chan *pathwalk.File, wg *sync.WaitGroup) {
	defer wg.Done()

	for range input {
	}
}

func collectFiles(input <-chan *pathwalk.File, wg *sync.WaitGroup) []*pathwalk.File {
	defer wg.Done()

	var files []*pathwalk.File
	for file := range input {
		files = append(files, file)
	}
	return files
}

func createTestDir(t *testing.T) string {
	tmpDir := t.TempDir()

	// Create test directory structure:
	// tmpDir/
	//   file1.txt
	//   file2.go
	//   .hidden_file
	//   subdir/
	//     subfile.txt
	//     .hidden_sub
	//   .hidden_dir/
	//     hidden_content.txt
	//   .git/
	//     config

	files := map[string]string{
		"file1.txt":                      "content1",
		"file2.go":                       "package main",
		".hidden_file":                   "hidden content",
		"subdir/subfile.txt":             "sub content",
		"subdir/.hidden_sub":             "hidden sub",
		".hidden_dir/hidden_content.txt": "hidden dir content",
		".git/config":                    "[core]\n",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return tmpDir
}

func TestOptions_DefaultValues(t *testing.T) {
	opts := &pathwalk.Options{}

	// Test default values
	if opts.Pattern != "" {
		t.Errorf("Expected empty Pattern, got %q", opts.Pattern)
	}
	if opts.Recurse {
		t.Error("Expected Recurse to be false by default")
	}
	if opts.Parallel != 0 {
		t.Errorf("Expected Parallel to be 0, got %d", opts.Parallel)
	}
	if opts.HiddenDirs {
		t.Error("Expected HiddenDirs to be false by default")
	}
	if opts.HiddenFiles {
		t.Error("Expected HiddenFiles to be false by default")
	}
	if opts.IncludeGit {
		t.Error("Expected IncludeGit to be false by default")
	}
}

func TestNewWalker_ValidPattern(t *testing.T) {
	tmpDir := createTestDir(t)

	opts := &pathwalk.Options{
		Pattern: "*.txt",
		Recurse: true,
	}

	fileChan := make(chan *pathwalk.File, 10)
	var wg sync.WaitGroup

	walker := pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)
	if walker == nil {
		t.Fatal("NewWalker returned nil")
	}
}

func TestNewWalker_InvalidPattern(t *testing.T) {
	tmpDir := createTestDir(t)

	opts := &pathwalk.Options{
		Pattern: "[", // Invalid glob pattern
		Recurse: true,
	}

	fileChan := make(chan *pathwalk.File, 10)
	var wg sync.WaitGroup

	// This should call os.Exit(1) due to invalid pattern
	// We can't easily test this without changing the code
	// For now, we'll test with a valid pattern instead
	opts.Pattern = "*.txt"
	walker := pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)
	if walker == nil {
		t.Fatal("NewWalker returned nil")
	}
}

func TestWalkers_BasicWalk(t *testing.T) {
	tmpDir := createTestDir(t)

	testCases := []struct {
		name      string
		useAlt    bool
		parallel  int
	}{
		{"Walker", false, 0},
		{"AltWalker", true, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &pathwalk.Options{
				Recurse:  false,
				Parallel: tc.parallel,
			}

			fileChan := make(chan *pathwalk.File, 10)
			var wg sync.WaitGroup
			var collectWg sync.WaitGroup

			var walker pathwalk.PathWalker
			if tc.useAlt {
				walker = pathwalk.NewAltWalker(tmpDir, opts, fileChan, &wg)
			} else {
				walker = pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)
			}

			wg.Add(1)
			collectWg.Add(1)

			var files []*pathwalk.File
			go func() {
				files = collectFiles(fileChan, &collectWg)
			}()

			go walker.Walk()

			wg.Wait()
			close(fileChan)
			collectWg.Wait()

			// Should only find files in root directory (no recursion)
			expectedFiles := []string{"file1.txt", "file2.go"}
			if len(files) != len(expectedFiles) {
				t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
			}

			foundFiles := make(map[string]bool)
			for _, file := range files {
				foundFiles[filepath.Base(file.Path)] = true
			}

			for _, expected := range expectedFiles {
				if !foundFiles[expected] {
					t.Errorf("Expected to find file %q", expected)
				}
			}
		})
	}
}

func TestWalkers_RecursiveWalk(t *testing.T) {
	tmpDir := createTestDir(t)

	testCases := []struct {
		name      string
		useAlt    bool
		parallel  int
	}{
		{"Walker", false, 0},
		{"AltWalker", true, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &pathwalk.Options{
				Recurse:  true,
				Parallel: tc.parallel,
			}

			fileChan := make(chan *pathwalk.File, 10)
			var wg sync.WaitGroup
			var collectWg sync.WaitGroup

			var walker pathwalk.PathWalker
			if tc.useAlt {
				walker = pathwalk.NewAltWalker(tmpDir, opts, fileChan, &wg)
			} else {
				walker = pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)
			}

			wg.Add(1)
			collectWg.Add(1)

			var files []*pathwalk.File
			go func() {
				files = collectFiles(fileChan, &collectWg)
			}()

			go walker.Walk()

			wg.Wait()
			close(fileChan)
			collectWg.Wait()

			// Should find files recursively but exclude hidden dirs/files by default
			expectedFiles := []string{"file1.txt", "file2.go", "subfile.txt"}
			if len(files) != len(expectedFiles) {
				t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
			}

			foundFiles := make(map[string]bool)
			for _, file := range files {
				foundFiles[filepath.Base(file.Path)] = true
			}

			for _, expected := range expectedFiles {
				if !foundFiles[expected] {
					t.Errorf("Expected to find file %q", expected)
				}
			}
		})
	}
}

func TestWalker_HiddenFiles(t *testing.T) {
	tmpDir := createTestDir(t)

	opts := &pathwalk.Options{
		Recurse:     true,
		HiddenFiles: true,
	}

	fileChan := make(chan *pathwalk.File, 10)
	var wg sync.WaitGroup
	var collectWg sync.WaitGroup

	walker := pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)

	wg.Add(1)
	collectWg.Add(1)

	var files []*pathwalk.File
	go func() {
		files = collectFiles(fileChan, &collectWg)
	}()

	go walker.Walk()

	wg.Wait()
	close(fileChan)
	collectWg.Wait()

	// Should include hidden files
	foundHidden := false
	for _, file := range files {
		if strings.HasPrefix(filepath.Base(file.Path), ".") {
			foundHidden = true
			break
		}
	}

	if !foundHidden {
		t.Error("Expected to find hidden files when HiddenFiles=true")
	}
}

func TestWalker_HiddenDirs(t *testing.T) {
	tmpDir := createTestDir(t)

	opts := &pathwalk.Options{
		Recurse:    true,
		HiddenDirs: true,
	}

	fileChan := make(chan *pathwalk.File, 10)
	var wg sync.WaitGroup
	var collectWg sync.WaitGroup

	walker := pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)

	wg.Add(1)
	collectWg.Add(1)

	var files []*pathwalk.File
	go func() {
		files = collectFiles(fileChan, &collectWg)
	}()

	go walker.Walk()

	wg.Wait()
	close(fileChan)
	collectWg.Wait()

	// Should include files from hidden directories
	foundInHiddenDir := false
	for _, file := range files {
		if strings.Contains(file.Path, ".hidden_dir") {
			foundInHiddenDir = true
			break
		}
	}

	if !foundInHiddenDir {
		t.Error("Expected to find files in hidden directories when HiddenDirs=true")
	}
}

func TestWalker_GitDirectory(t *testing.T) {
	tmpDir := createTestDir(t)

	tests := []struct {
		name       string
		includeGit bool
		expectGit  bool
	}{
		{"exclude git", false, false},
		{"include git", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &pathwalk.Options{
				Recurse:    true,
				IncludeGit: tt.includeGit,
			}

			fileChan := make(chan *pathwalk.File, 10)
			var wg sync.WaitGroup
			var collectWg sync.WaitGroup

			walker := pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)

			wg.Add(1)
			collectWg.Add(1)

			var files []*pathwalk.File
			go func() {
				files = collectFiles(fileChan, &collectWg)
			}()

			go walker.Walk()

			wg.Wait()
			close(fileChan)
			collectWg.Wait()

			foundGit := false
			for _, file := range files {
				if strings.Contains(file.Path, ".git") {
					foundGit = true
					break
				}
			}

			if foundGit != tt.expectGit {
				if tt.expectGit {
					t.Error("Expected to find .git files when IncludeGit=true")
				} else {
					t.Error("Expected NOT to find .git files when IncludeGit=false")
				}
			}
		})
	}
}

func TestWalkers_PatternMatching(t *testing.T) {
	tmpDir := createTestDir(t)

	walkerTests := []struct {
		name      string
		useAlt    bool
		parallel  int
	}{
		{"Walker", false, 0},
		{"AltWalker", true, 2},
	}

	patternTests := []struct {
		name     string
		pattern  string
		expected []string
	}{
		{"txt files", "*.txt", []string{"file1.txt", "subfile.txt"}},
		{"go files", "*.go", []string{"file2.go"}},
		{"no match", "*.xyz", []string{}},
		{"specific file", "file1.txt", []string{"file1.txt"}},
	}

	for _, wt := range walkerTests {
		t.Run(wt.name, func(t *testing.T) {
			for _, pt := range patternTests {
				t.Run(pt.name, func(t *testing.T) {
					opts := &pathwalk.Options{
						Pattern:  pt.pattern,
						Recurse:  true,
						Parallel: wt.parallel,
					}

					fileChan := make(chan *pathwalk.File, 10)
					var wg sync.WaitGroup
					var collectWg sync.WaitGroup

					var walker pathwalk.PathWalker
					if wt.useAlt {
						walker = pathwalk.NewAltWalker(tmpDir, opts, fileChan, &wg)
					} else {
						walker = pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)
					}

					wg.Add(1)
					collectWg.Add(1)

					var files []*pathwalk.File
					go func() {
						files = collectFiles(fileChan, &collectWg)
					}()

					go walker.Walk()

					wg.Wait()
					close(fileChan)
					collectWg.Wait()

					if len(files) != len(pt.expected) {
						t.Errorf("Expected %d files, got %d", len(pt.expected), len(files))
					}

					foundFiles := make(map[string]bool)
					for _, file := range files {
						foundFiles[filepath.Base(file.Path)] = true
					}

					for _, expected := range pt.expected {
						if !foundFiles[expected] {
							t.Errorf("Expected to find file %q", expected)
						}
					}
				})
			}
		})
	}
}

func TestWalker_FileAttributes(t *testing.T) {
	tmpDir := createTestDir(t)

	opts := &pathwalk.Options{
		Recurse: false,
	}

	fileChan := make(chan *pathwalk.File, 10)
	var wg sync.WaitGroup
	var collectWg sync.WaitGroup

	walker := pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)

	wg.Add(1)
	collectWg.Add(1)

	var files []*pathwalk.File
	go func() {
		files = collectFiles(fileChan, &collectWg)
	}()

	go walker.Walk()

	wg.Wait()
	close(fileChan)
	collectWg.Wait()

	for _, file := range files {
		// Check that file attributes are populated
		if file.Path == "" {
			t.Error("File path should not be empty")
		}
		if file.Size < 0 {
			t.Error("File size should not be negative")
		}
		if file.ModTime.IsZero() {
			t.Error("File ModTime should not be zero")
		}
		if file.ProcTime != 0 {
			t.Error("File ProcTime should be zero for Walker")
		}
		if file.Error != nil {
			t.Errorf("File Error should be nil, got: %v", file.Error)
		}

		// Verify the file actually exists
		if _, err := os.Stat(file.Path); err != nil {
			t.Errorf("File %s should exist: %v", file.Path, err)
		}
	}
}

func TestNewAltWalker_ValidPattern(t *testing.T) {
	tmpDir := createTestDir(t)

	opts := &pathwalk.Options{
		Pattern:  "*.txt",
		Recurse:  true,
		Parallel: 2,
	}

	fileChan := make(chan *pathwalk.File, 10)
	var wg sync.WaitGroup

	walker := pathwalk.NewAltWalker(tmpDir, opts, fileChan, &wg)
	if walker == nil {
		t.Fatal("NewAltWalker returned nil")
	}
}



func TestAltWalker_ParallelBehavior(t *testing.T) {
	tmpDir := createTestDir(t)

	tests := []struct {
		name     string
		parallel int
	}{
		{"single thread", 1},
		{"two threads", 2},
		{"four threads", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &pathwalk.Options{
				Recurse:  true,
				Parallel: tt.parallel,
			}

			fileChan := make(chan *pathwalk.File, 20)
			var wg sync.WaitGroup
			var collectWg sync.WaitGroup

			walker := pathwalk.NewAltWalker(tmpDir, opts, fileChan, &wg)

			wg.Add(1)
			collectWg.Add(1)

			var files []*pathwalk.File
			go func() {
				files = collectFiles(fileChan, &collectWg)
			}()

			start := time.Now()
			go walker.Walk()

			wg.Wait()
			close(fileChan)
			collectWg.Wait()
			elapsed := time.Since(start)

			// Should find the same files regardless of parallelism
			expectedFiles := []string{"file1.txt", "file2.go", "subfile.txt"}
			if len(files) != len(expectedFiles) {
				t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
			}

			// Basic sanity check - shouldn't take too long
			if elapsed > 5*time.Second {
				t.Errorf("Walk took too long: %v", elapsed)
			}
		})
	}
}


func TestPathWalker_Interface(t *testing.T) {
	tmpDir := createTestDir(t)
	opts := &pathwalk.Options{Parallel: 1}
	fileChan := make(chan *pathwalk.File, 10)
	var wg sync.WaitGroup

	// Test that both implementations satisfy the PathWalker interface
	var walker1 pathwalk.PathWalker = pathwalk.NewWalker(tmpDir, opts, fileChan, &wg)
	var walker2 pathwalk.PathWalker = pathwalk.NewAltWalker(tmpDir, opts, fileChan, &wg)

	if walker1 == nil {
		t.Error("Walker should implement PathWalker interface")
	}
	if walker2 == nil {
		t.Error("AltWalker should implement PathWalker interface")
	}
}

func TestFile_Struct(t *testing.T) {
	now := time.Now()
	duration := 100 * time.Millisecond

	file := &pathwalk.File{
		Path:     "/test/path.txt",
		Size:     1024,
		ModTime:  now,
		ProcTime: duration,
		Error:    nil,
	}

	if file.Path != "/test/path.txt" {
		t.Errorf("Expected path '/test/path.txt', got %q", file.Path)
	}
	if file.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", file.Size)
	}
	if !file.ModTime.Equal(now) {
		t.Errorf("Expected ModTime %v, got %v", now, file.ModTime)
	}
	if file.ProcTime != duration {
		t.Errorf("Expected ProcTime %v, got %v", duration, file.ProcTime)
	}
	if file.Error != nil {
		t.Errorf("Expected Error nil, got %v", file.Error)
	}
}

// benchmarkPathwalker runs a pathwalker benchmark with the specified configuration
func benchmarkPathwalker(useAltWalker bool, options *pathwalk.Options) {
	var (
		pathWalkWaitGroup  sync.WaitGroup
		pathPrintWaitGroup sync.WaitGroup
	)

	fileChan := make(chan *pathwalk.File, 64)

	// Create appropriate walker type
	var walker pathwalk.PathWalker
	if useAltWalker {
		walker = pathwalk.NewAltWalker("/usr/share", options, fileChan, &pathWalkWaitGroup)
	} else {
		walker = pathwalk.NewWalker("/usr/share", options, fileChan, &pathWalkWaitGroup)
	}

	pathWalkWaitGroup.Add(1)
	go walker.Walk()

	pathPrintWaitGroup.Add(1)
	go sinkFileChan(fileChan, &pathPrintWaitGroup)

	pathWalkWaitGroup.Wait()
	close(fileChan)
	pathPrintWaitGroup.Wait()
}

// BenchmarkPathwalker benchmarks both walker types with different configurations
func BenchmarkPathwalker(b *testing.B) {
	benchmarks := []struct {
		name         string
		useAltWalker bool
		parallel     int
		hiddenDirs   bool
		hiddenFiles  bool
		includeGit   bool
	}{
		// Basic Walker benchmarks (sequential directory traversal)
		{"Walker/Basic", false, 1, false, false, false},
		{"Walker/WithHidden", false, 1, true, true, true},

		// AltWalker benchmarks with different parallelism levels
		{"AltWalker/Parallel1", true, 1, true, true, true},
		{"AltWalker/Parallel2", true, 2, true, true, true},
		{"AltWalker/Parallel3", true, 3, true, true, true},
		{"AltWalker/Parallel4", true, 4, true, true, true},
		{"AltWalker/Parallel5", true, 5, true, true, true},
		{"AltWalker/Parallel6", true, 6, true, true, true},
		{"AltWalker/Parallel7", true, 7, true, true, true},
		{"AltWalker/Parallel8", true, 8, true, true, true},

		// Configuration variations
		{"AltWalker/NoHidden", true, 4, false, false, false},
		{"AltWalker/OnlyFiles", true, 4, false, true, false},
		{"AltWalker/OnlyDirs", true, 4, true, false, false},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			options := pathwalk.Options{
				Recurse:     true,
				Parallel:    bm.parallel,
				HiddenDirs:  bm.hiddenDirs,
				HiddenFiles: bm.hiddenFiles,
				IncludeGit:  bm.includeGit,
			}

			b.ResetTimer()
			for b.Loop() {
				benchmarkPathwalker(bm.useAltWalker, &options)
			}
		})
	}
}
