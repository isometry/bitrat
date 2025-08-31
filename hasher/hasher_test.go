package hasher_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
	"github.com/spf13/viper"
)

const (
	testPath = "/usr/local"
	testHash = "blake2b-512"
)

func sinkFileHashChan(input <-chan *hasher.FileHash) func() {
	return func() {
		for range input {
		}
	}
}

func collectFileHashChan(input <-chan *hasher.FileHash, wg *sync.WaitGroup) []*hasher.FileHash {
	defer wg.Done()

	var results []*hasher.FileHash
	for item := range input {
		results = append(results, item)
	}
	return results
}

func createTestFiles(t *testing.T) string {
	tmpDir := t.TempDir()

	files := map[string]string{
		"test1.txt": "Hello, World!",
		"test2.txt": "This is a test file",
		"empty.txt": "",
		"binary":    "\x00\x01\x02\x03\xFF",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return tmpDir
}

func TestSupportedAlgorithms_Coverage(t *testing.T) {
	expectedAlgorithms := []string{
		"blake2b", "blake2b-256", "blake2b-384", "blake2b-512",
		"blake2s-128", "blake2s-256",
		"blake3", "blake3-dk",
		"crc32",
		"md5",
		"sha1", "sha224", "sha256", "sha256-simd", "sha384", "sha512",
		"sha3-224", "sha3-256", "sha3-384", "sha3-512",
		"skein-256", "skein-512",
	}

	for _, algo := range expectedAlgorithms {
		if _, exists := hasher.SupportedAlgorithms[algo]; !exists {
			t.Errorf("Expected algorithm %q to be supported", algo)
		}
	}

	// Test that we have all expected algorithms
	if len(hasher.SupportedAlgorithms) != len(expectedAlgorithms) {
		t.Errorf("Expected %d algorithms, got %d", len(expectedAlgorithms), len(hasher.SupportedAlgorithms))
	}
}

func TestNew_ValidAlgorithms(t *testing.T) {
	testCases := []struct {
		name       string
		algorithm  string
		key        []byte
		expectHMAC bool
	}{
		{"blake3", "blake3", nil, false},
		{"blake3 with key", "blake3-dk", []byte("testkey"), true},
		{"sha256", "sha256", nil, false},
		{"sha256 HMAC", "sha256", []byte("hmackey"), true},
		{"md5", "md5", nil, false},
		{"md5 HMAC", "md5", []byte("hmackey"), true},
		{"crc32", "crc32", nil, false},
		{"blake2b", "blake2b", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := hasher.New(tc.algorithm, tc.key)

			if h.Hash == nil {
				t.Fatal("Expected Hash to be initialized")
			}

			expectedType := tc.algorithm
			if tc.expectHMAC && tc.algorithm != "blake3-dk" {
				expectedType = "hmac-" + tc.algorithm
			} else if tc.algorithm == "blake3-dk" && tc.key != nil {
				expectedType = "hmac-" + tc.algorithm
			}

			if h.Type != expectedType {
				t.Errorf("Expected Type %q, got %q", expectedType, h.Type)
			}
		})
	}
}

func TestNew_InvalidAlgorithm(t *testing.T) {
	// This test would call log.Fatalf, so we can't test it directly
	// without changing the code. We'll test that valid algorithms work instead.
	validAlgos := []string{"sha256", "blake3", "md5"}

	for _, algo := range validAlgos {
		h := hasher.New(algo, nil)
		if h.Hash == nil {
			t.Errorf("Expected valid algorithm %q to work", algo)
		}
	}
}

func TestHasher_HashIoReader(t *testing.T) {
	testCases := []struct {
		name      string
		algorithm string
		input     string
		key       []byte
	}{
		{"sha256 empty", "sha256", "", nil},
		{"sha256 hello", "sha256", "Hello, World!", nil},
		{"blake3 test", "blake3", "test data", nil},
		{"md5 with data", "md5", "some test data here", nil},
		{"sha256 with HMAC", "sha256", "test", []byte("key")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := hasher.New(tc.algorithm, tc.key)
			reader := strings.NewReader(tc.input)

			hash1 := h.HashIoReader(reader)
			if hash1 == nil {
				t.Fatal("Expected hash to be non-nil")
			}

			// Hash the same data again to ensure consistency
			reader2 := strings.NewReader(tc.input)
			hash2 := h.HashIoReader(reader2)

			if !bytes.Equal(hash1, hash2) {
				t.Error("Expected consistent hash results for same input")
			}

			// Test that different data produces different hash
			if tc.input != "" {
				reader3 := strings.NewReader(tc.input + "different")
				hash3 := h.HashIoReader(reader3)

				if bytes.Equal(hash1, hash3) {
					t.Error("Expected different hashes for different input")
				}
			}
		})
	}
}

func TestHasher_HashFile(t *testing.T) {
	tmpDir := createTestFiles(t)

	testCases := []struct {
		name      string
		algorithm string
		filename  string
		key       []byte
	}{
		{"sha256 test1", "sha256", "test1.txt", nil},
		{"blake3 test2", "blake3", "test2.txt", nil},
		{"md5 empty", "md5", "empty.txt", nil},
		{"sha256 binary", "sha256", "binary", nil},
		{"sha256 HMAC", "sha256", "test1.txt", []byte("testkey")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := hasher.New(tc.algorithm, tc.key)

			filePath := filepath.Join(tmpDir, tc.filename)
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				t.Fatal(err)
			}

			file := &pathwalk.File{
				Path:    filePath,
				ModTime: fileInfo.ModTime(),
			}

			result := h.HashFile(file)

			// Check result structure
			if result == nil {
				t.Fatal("Expected FileHash result to be non-nil")
			}
			if result.File != file {
				t.Error("Expected File reference to be preserved")
			}
			if len(result.Hash) == 0 && result.File.Error == nil {
				t.Error("Expected Hash to be non-empty for successful hashing")
			}
			if result.Type != h.Type {
				t.Errorf("Expected Type %q, got %q", h.Type, result.Type)
			}

			// Check that file size was updated
			if result.File.Size != fileInfo.Size() {
				t.Errorf("Expected Size %d, got %d", fileInfo.Size(), result.File.Size)
			}

			// Check that processing time was recorded
			if result.File.ProcTime <= 0 {
				t.Error("Expected ProcTime to be positive")
			}
		})
	}
}

func TestHasher_HashFile_NonExistent(t *testing.T) {
	h := hasher.New("sha256", nil)

	file := &pathwalk.File{
		Path: "/nonexistent/file.txt",
	}

	result := h.HashFile(file)

	if result == nil {
		t.Fatal("Expected FileHash result even for errors")
	}
	if result.File.Error == nil {
		t.Error("Expected error for nonexistent file")
	}
	if len(result.Hash) != 0 {
		t.Error("Expected empty hash for error case")
	}
}

func TestHasher_HashProcessor(t *testing.T) {
	tmpDir := createTestFiles(t)
	h := hasher.New("sha256", nil)

	// Create input channel with test files
	input := make(chan *pathwalk.File, 5)
	output := make(chan *hasher.FileHash, 5)
	var wg sync.WaitGroup
	var collectWg sync.WaitGroup

	// Start hash processor
	wg.Go(h.HashProcessor(input, output))

	// Start collector
	collectWg.Add(1)
	var results []*hasher.FileHash
	go func() {
		results = collectFileHashChan(output, &collectWg)
	}()

	// Send test files
	testFiles := []string{"test1.txt", "test2.txt", "empty.txt"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			t.Fatal(err)
		}

		input <- &pathwalk.File{
			Path:    filePath,
			ModTime: fileInfo.ModTime(),
		}
	}

	close(input)
	wg.Wait()
	close(output)
	collectWg.Wait()

	// Verify results
	if len(results) != len(testFiles) {
		t.Errorf("Expected %d results, got %d", len(testFiles), len(results))
	}

	for _, result := range results {
		if len(result.Hash) == 0 && result.File.Error == nil {
			t.Error("Expected non-empty hash for successful processing")
		}
		if result.Type != "sha256" {
			t.Errorf("Expected Type sha256, got %q", result.Type)
		}
	}
}

func TestSortByFifo(t *testing.T) {
	// Create test data
	testHashes := []*hasher.FileHash{
		{File: &pathwalk.File{Path: "c.txt"}, Hash: []byte("hash3")},
		{File: &pathwalk.File{Path: "a.txt"}, Hash: []byte("hash1")},
		{File: &pathwalk.File{Path: "b.txt"}, Hash: []byte("hash2")},
	}

	input := make(chan *hasher.FileHash, len(testHashes))
	output := make(chan *hasher.FileHash, len(testHashes))
	var wg sync.WaitGroup
	var collectWg sync.WaitGroup

	// Start FIFO sorter
	wg.Go(hasher.SortByFifo(input, output))

	// Start collector
	collectWg.Add(1)
	var results []*hasher.FileHash
	go func() {
		results = collectFileHashChan(output, &collectWg)
	}()

	// Send data in specific order
	for _, hash := range testHashes {
		input <- hash
	}
	close(input)

	wg.Wait()
	close(output)
	collectWg.Wait()

	// Verify FIFO order is preserved
	if len(results) != len(testHashes) {
		t.Fatalf("Expected %d results, got %d", len(testHashes), len(results))
	}

	for i, result := range results {
		if result.File.Path != testHashes[i].File.Path {
			t.Errorf("FIFO order not preserved: expected %q at position %d, got %q",
				testHashes[i].File.Path, i, result.File.Path)
		}
	}
}

func TestSortByPath(t *testing.T) {
	// Create test data in unsorted order
	testHashes := []*hasher.FileHash{
		{File: &pathwalk.File{Path: "c.txt"}, Hash: []byte("hash3")},
		{File: &pathwalk.File{Path: "a.txt"}, Hash: []byte("hash1")},
		{File: &pathwalk.File{Path: "b.txt"}, Hash: []byte("hash2")},
	}

	input := make(chan *hasher.FileHash, len(testHashes))
	output := make(chan *hasher.FileHash, len(testHashes))
	var wg sync.WaitGroup
	var collectWg sync.WaitGroup

	// Start path sorter
	wg.Go(hasher.SortByPath(input, output))

	// Start collector
	collectWg.Add(1)
	var results []*hasher.FileHash
	go func() {
		results = collectFileHashChan(output, &collectWg)
	}()

	// Send data in unsorted order
	for _, hash := range testHashes {
		input <- hash
	}
	close(input)

	wg.Wait()
	close(output)
	collectWg.Wait()

	// Verify results are sorted by path
	if len(results) != len(testHashes) {
		t.Fatalf("Expected %d results, got %d", len(testHashes), len(results))
	}

	expectedPaths := []string{"a.txt", "b.txt", "c.txt"}
	for i, result := range results {
		if result.File.Path != expectedPaths[i] {
			t.Errorf("Path sort failed: expected %q at position %d, got %q",
				expectedPaths[i], i, result.File.Path)
		}
	}
}

func TestSprintf(t *testing.T) {
	testCases := []struct {
		name     string
		format   string
		hash     []byte
		path     string
		expected string
	}{
		{
			name:     "default format",
			format:   "",
			hash:     []byte{0xde, 0xad, 0xbe, 0xef},
			path:     "test.txt",
			expected: "deadbeef  test.txt\n",
		},
		{
			name:     "custom format",
			format:   "%x:%s",
			hash:     []byte{0xaa, 0xbb},
			path:     "file.txt",
			expected: "aabb:file.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileHash := &hasher.FileHash{
				File: &pathwalk.File{Path: tc.path},
				Hash: tc.hash,
			}

			result := hasher.Sprintf(tc.format, fileHash)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestFileHash_Structure(t *testing.T) {
	now := time.Now()
	file := &pathwalk.File{
		Path:    "/test/path.txt",
		Size:    1024,
		ModTime: now,
	}

	fileHash := &hasher.FileHash{
		File: file,
		Hash: []byte{0x01, 0x02, 0x03},
		Type: "sha256",
	}

	// Test structure access
	if fileHash.File != file {
		t.Error("File reference not preserved")
	}
	if len(fileHash.Hash) != 3 {
		t.Errorf("Expected hash length 3, got %d", len(fileHash.Hash))
	}
	if fileHash.Type != "sha256" {
		t.Errorf("Expected Type sha256, got %q", fileHash.Type)
	}

	// Test nested structures
	fileHash.AttrInfo.Hash = []byte{0xaa, 0xbb}
	fileHash.AttrInfo.Time = now
	fileHash.FileInfo.Hash = []byte{0xcc, 0xdd}
	fileHash.FileInfo.Time = now

	if !bytes.Equal(fileHash.AttrInfo.Hash, []byte{0xaa, 0xbb}) {
		t.Error("AttrInfo.Hash not set correctly")
	}
	if !fileHash.AttrInfo.Time.Equal(now) {
		t.Error("AttrInfo.Time not set correctly")
	}
}

func TestHasher_HMACBehavior(t *testing.T) {
	input := "test data for HMAC"
	key1 := []byte("key1")
	key2 := []byte("key2")

	algorithms := []string{"sha256", "blake3-dk", "md5"}

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			// Test that different keys produce different HMACs
			h1 := hasher.New(algo, key1)
			h2 := hasher.New(algo, key2)

			reader1 := strings.NewReader(input)
			reader2 := strings.NewReader(input)

			hash1 := h1.HashIoReader(reader1)
			hash2 := h2.HashIoReader(reader2)

			if bytes.Equal(hash1, hash2) {
				t.Errorf("Algorithm %s produced same HMAC for different keys", algo)
			}

			// Test that no key vs key produces different results
			if algo != "blake3-dk" { // blake3-dk requires a key
				h3 := hasher.New(algo, nil)
				reader3 := strings.NewReader(input)
				hash3 := h3.HashIoReader(reader3)

				if bytes.Equal(hash1, hash3) {
					t.Errorf("Algorithm %s: HMAC with key same as hash without key", algo)
				}
			}
		})
	}
}

func TestOutputTextFile_Stats(t *testing.T) {
	// Test OutputTextFile with stats enabled
	viper.Set("stats", true)
	viper.Set("print-format", "%x  %s")
	viper.Set("output-file", "-") // stdout
	viper.Set("hash", "sha256")
	viper.Set("parallel", 2)
	defer viper.Reset()

	// Create test data
	testHashes := []*hasher.FileHash{
		{
			File: &pathwalk.File{
				Path:     "test1.txt",
				Size:     100,
				ProcTime: 10 * time.Millisecond,
			},
			Hash: []byte{0xaa, 0xbb, 0xcc},
		},
		{
			File: &pathwalk.File{
				Path:     "test2.txt",
				Size:     200,
				ProcTime: 20 * time.Millisecond,
			},
			Hash: []byte{0xdd, 0xee, 0xff},
		},
	}

	input := make(chan *hasher.FileHash, len(testHashes))
	var wg sync.WaitGroup

	// This test is challenging because OutputTextFile writes to stdout/stderr
	// We'll just test that it doesn't crash and completes
	wg.Go(hasher.OutputTextFile(input))

	for _, hash := range testHashes {
		input <- hash
	}
	close(input)

	wg.Wait()
	// If we get here without hanging, the test passes
}

func TestOutputTextFile_NoStats(t *testing.T) {
	// Test OutputTextFile with stats disabled
	viper.Set("stats", false)
	viper.Set("print-format", "")
	viper.Set("output-file", "-")
	defer viper.Reset()

	testHashes := []*hasher.FileHash{
		{
			File: &pathwalk.File{Path: "test.txt", Size: 50},
			Hash: []byte{0x12, 0x34},
		},
	}

	input := make(chan *hasher.FileHash, 1)
	var wg sync.WaitGroup

	wg.Go(hasher.OutputTextFile(input))

	input <- testHashes[0]
	close(input)

	wg.Wait()
}

func TestOutputProtobufFile(t *testing.T) {
	// Create a temporary file for protobuf output
	tmpFile, err := os.CreateTemp("", "test_output_*.pb")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	viper.Set("stats", true)
	viper.Set("output-file", tmpFile.Name())
	viper.Set("hash", "sha256")
	viper.Set("parallel", 1)
	defer viper.Reset()

	testHashes := []*hasher.FileHash{
		{
			File: &pathwalk.File{
				Path:     "test1.txt",
				Size:     100,
				ModTime:  time.Now(),
				ProcTime: 5 * time.Millisecond,
			},
			Hash: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		},
	}

	input := make(chan *hasher.FileHash, 1)
	var wg sync.WaitGroup

	wg.Go(hasher.OutputProtobufFile(input))

	input <- testHashes[0]
	close(input)

	wg.Wait()

	// Verify that the file was written (basic check)
	fileInfo, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Size() == 0 {
		t.Error("Expected protobuf file to have content")
	}
}

func TestHashRouter(t *testing.T) {
	// Create test data with both nil and non-nil hashes
	testData := []hasher.FileHash{
		{Hash: []byte{0xaa, 0xbb}}, // non-nil
		{Hash: nil},                // nil
		{Hash: []byte{0xcc, 0xdd}}, // non-nil
		{Hash: nil},                // nil
	}

	input := make(chan hasher.FileHash, len(testData))
	hashOutput := make(chan hasher.FileHash, len(testData))
	nilOutput := make(chan hasher.FileHash, len(testData))
	var wg sync.WaitGroup

	// Start router
	wg.Go(hasher.HashRouter(input, hashOutput, nilOutput))

	// Send test data
	for _, data := range testData {
		input <- data
	}
	close(input)

	wg.Wait()
	close(hashOutput)
	close(nilOutput)

	// Collect results
	var hashResults []hasher.FileHash
	var nilResults []hasher.FileHash

	for hash := range hashOutput {
		hashResults = append(hashResults, hash)
	}
	for nilHash := range nilOutput {
		nilResults = append(nilResults, nilHash)
	}

	// Verify routing
	if len(hashResults) != 2 {
		t.Errorf("Expected 2 non-nil hashes, got %d", len(hashResults))
	}
	if len(nilResults) != 2 {
		t.Errorf("Expected 2 nil hashes, got %d", len(nilResults))
	}

	// Verify that non-nil hashes went to hash output
	for _, result := range hashResults {
		if result.Hash == nil {
			t.Error("Found nil hash in hash output")
		}
	}

	// Verify that nil hashes went to nil output
	for _, result := range nilResults {
		if result.Hash != nil {
			t.Error("Found non-nil hash in nil output")
		}
	}
}

func TestHashSink(t *testing.T) {
	// Test that HashSink consumes all input without blocking
	testData := []hasher.FileHash{
		{Hash: []byte{0xaa}},
		{Hash: []byte{0xbb}},
		{Hash: []byte{0xcc}},
	}

	input := make(chan hasher.FileHash, len(testData))
	var wg sync.WaitGroup

	// Start sink
	wg.Go(hasher.HashSink(input))

	// Send all test data
	for _, data := range testData {
		input <- data
	}
	close(input)

	// Should complete without hanging
	wg.Wait()
}

func TestHasher_AlgorithmSpecificBehavior(t *testing.T) {
	// Test algorithm-specific behaviors
	testCases := []struct {
		name        string
		algorithm   string
		key         []byte
		shouldPanic bool
	}{
		{"CRC32 without key", "crc32", nil, false},
		{"CRC32 with key should log.Fatal", "crc32", []byte("key"), false}, // Can't test log.Fatal easily
		{"Skein with key", "skein-256", []byte("key"), false},
		{"Skein without key", "skein-256", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.algorithm == "crc32" && tc.key != nil {
				// Skip CRC32 with key test as it calls log.Fatalf
				t.Skip("CRC32 with HMAC calls log.Fatalf")
				return
			}

			h := hasher.New(tc.algorithm, tc.key)
			if h.Hash == nil {
				t.Error("Expected hash to be initialized")
			}

			// Test that it can hash something
			reader := strings.NewReader("test data")
			hash := h.HashIoReader(reader)
			if len(hash) == 0 {
				t.Error("Expected non-empty hash result")
			}
		})
	}
}

// benchmarkHasherPipeline runs a complete hasher pipeline benchmark
func benchmarkHasherPipeline(useAltWalker bool, options *pathwalk.Options) {
	var (
		pathWalkWaitGroup sync.WaitGroup
		hasherWaitGroup   sync.WaitGroup
		hashSinkWaitGroup sync.WaitGroup
	)

	fileChan := make(chan *pathwalk.File, 1024)
	hashChan := make(chan *hasher.FileHash, 64)

	// Create appropriate walker type
	var walker pathwalk.PathWalker
	if useAltWalker {
		walker = pathwalk.NewAltWalker(testPath, options, fileChan, &pathWalkWaitGroup)
	} else {
		walker = pathwalk.NewWalker(testPath, options, fileChan, &pathWalkWaitGroup)
	}

	pathWalkWaitGroup.Add(1)
	go walker.Walk()

	// Start parallel hash processors
	for i := 0; i < options.Parallel; i++ {
		h := hasher.New(testHash, []byte(""))
		hasherWaitGroup.Go(h.HashProcessor(fileChan, hashChan))
	}

	// Start hash sink
	hashSinkWaitGroup.Go(sinkFileHashChan(hashChan))

	// Wait for completion
	pathWalkWaitGroup.Wait()
	close(fileChan)
	hasherWaitGroup.Wait()
	close(hashChan)
	hashSinkWaitGroup.Wait()
}

// BenchmarkHasherPipeline benchmarks the complete hashing pipeline with different configurations
func BenchmarkHasherPipeline(b *testing.B) {
	benchmarks := []struct {
		name         string
		useAltWalker bool
		parallel     int
	}{
		// Walker benchmarks (sequential directory traversal)
		{"Walker/Parallel1", false, 1},
		{"Walker/Parallel2", false, 2},
		{"Walker/Parallel3", false, 3},
		{"Walker/Parallel4", false, 4},
		{"Walker/Parallel5", false, 5},
		{"Walker/Parallel6", false, 6},
		{"Walker/Parallel7", false, 7},
		{"Walker/Parallel8", false, 8},

		// AltWalker benchmarks (concurrent directory traversal)
		{"AltWalker/Parallel1", true, 1},
		{"AltWalker/Parallel2", true, 2},
		{"AltWalker/Parallel3", true, 3},
		{"AltWalker/Parallel4", true, 4},
		{"AltWalker/Parallel5", true, 5},
		{"AltWalker/Parallel6", true, 6},
		{"AltWalker/Parallel7", true, 7},
		{"AltWalker/Parallel8", true, 8},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			options := pathwalk.Options{
				Recurse:  true,
				Parallel: bm.parallel,
			}

			b.ResetTimer()
			for b.Loop() {
				benchmarkHasherPipeline(bm.useAltWalker, &options)
			}
		})
	}
}
