package hasher

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/dchest/skein"
	"github.com/isometry/bitrat/pathwalk"
	"github.com/isometry/bitrat/protobuf/bitratpb"
	sha256simd "github.com/minio/sha256-simd"
	"github.com/spf13/viper"
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Hasher type
type Hasher struct {
	Type string
	Hash hash.Hash
}

// FileHash type
type FileHash struct {
	File     *pathwalk.File
	Hash     []byte
	AttrInfo struct {
		Hash []byte
		Time time.Time
	}
	FileInfo struct {
		Hash []byte
		Time time.Time
	}
	Type string
}

type fileHashByPath []*FileHash

func (a fileHashByPath) Len() int           { return len(a) }
func (a fileHashByPath) Less(i, j int) bool { return a[i].File.Path < a[j].File.Path }
func (a fileHashByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type algorithm interface{}

type algoSize struct {
	algo algorithm
	size int
}

func newBlake3Hash() hash.Hash {
	return hash.Hash(blake3.New())
}

func newBlake3DerivedKey(key []byte) (hash.Hash, error) {
	return hash.Hash(blake3.NewDeriveKey(string(key))), nil
}

// SupportedAlgorithms maps algorithm names to implementations
var SupportedAlgorithms = map[string]algorithm{
	"blake2b":     blake2b.New512,
	"blake2b-256": blake2b.New256,
	"blake2b-384": blake2b.New384,
	"blake2b-512": blake2b.New512,
	"blake2s-128": blake2s.New128,
	"blake2s-256": blake2s.New256,
	"blake3":      newBlake3Hash, // XXX: broken on Intel Atom C3000
	"blake3-dk":   newBlake3DerivedKey,
	"crc32":       crc32.NewIEEE,
	"md5":         md5.New,
	"ripemd160":   ripemd160.New,
	"sha1":        sha1.New,
	"sha224":      sha256.New224,
	"sha256":      sha256simd.New,
	"sha384":      sha512.New384,
	"sha512":      sha512.New,
	"sha3-224":    sha3.New224,
	"sha3-256":    sha3.New256,
	"sha3-384":    sha3.New384,
	"sha3-512":    sha3.New512,
	"skein-256":   algoSize{skein.NewMAC, 256 / 8},
	"skein-512":   algoSize{skein.NewMAC, 512 / 8},
}

// New returns an initialised instance of the appropriate hash function
func New(algo string, key []byte) Hasher {
	var name string
	switch string(key) {
	case "":
		name = algo
	default:
		name = fmt.Sprintf("hmac-%s", algo)
	}

	if selectedAlgorithm, ok := SupportedAlgorithms[algo]; ok {
		switch hashFn := selectedAlgorithm.(type) {
		// hashes with direct HMAC support
		case func([]byte) (hash.Hash, error):
			h, e := hashFn(key)
			if e != nil {
				log.Fatal(e)
			}
			return Hasher{
				Type: name,
				Hash: h,
			}
		// hashes without direct HMAC support
		case func() hash.Hash:
			if bytes.Equal(key, nil) {
				return Hasher{
					Type: name,
					Hash: hashFn(),
				}
			}
			return Hasher{
				Type: name,
				Hash: hmac.New(hashFn, key),
			}
		// support for hash.Hash32 hashes (i.e. CRC-32 and Adler-32)
		case func() hash.Hash32:
			if bytes.Equal(key, nil) {
				return Hasher{
					Type: name,
					Hash: hashFn(),
				}
			}
			log.Fatalf("HMAC unsupported for %s\n", algo)
		// hashes requiring size argument
		case algoSize:
			switch sHashFn := hashFn.algo.(type) {
			case func(int, []byte) hash.Hash:
				h := sHashFn(hashFn.size, key)
				return Hasher{
					Type: name,
					Hash: h,
				}
			case func(uint64, []byte) hash.Hash:
				h := sHashFn(uint64(hashFn.size), key)
				return Hasher{
					Type: name,
					Hash: h,
				}
			}
		}
	}

	log.Fatalf("unsupported hash algorithm: %v\n", name)
	return Hasher{}
}

// HashFile calculates the hash/hmac of file at path
func (hasher *Hasher) HashFile(file *pathwalk.File) *FileHash {
	defer hasher.Hash.Reset()

	fd, err := os.Open(file.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "# %s: %v\n", file.Path, err)
		file.Error = err
		return &FileHash{
			File: file,
			Type: hasher.Type,
		}
	}
	defer fd.Close()

	start := time.Now()
	size, err := io.Copy(hasher.Hash, fd)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	file.Size = size
	file.ProcTime = elapsed

	return &FileHash{
		File: file,
		Hash: hasher.Hash.Sum(nil),
		Type: hasher.Type,
	}
}

// HashIoReader calculates the hash/hmac of stdin
func (hasher *Hasher) HashIoReader(reader io.Reader) []byte {
	defer hasher.Hash.Reset()

	_, err := io.Copy(hasher.Hash, reader)
	if err != nil {
		log.Fatal(err)
	}

	return hasher.Hash.Sum(nil)
}

// HashProcessor goroutine
func (hasher *Hasher) HashProcessor(input <-chan *pathwalk.File, output chan<- *FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		output <- hasher.HashFile(item)
	}
}

// SortByFifo takes an input channel of FileHash and adds each item directly to the output in the order received
func SortByFifo(input <-chan *FileHash, output chan<- *FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		output <- item
	}
}

// SortByPath takes an input channel of FileHash and lexicographically sorts by file path
func SortByPath(input <-chan *FileHash, output chan<- *FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	var hashes []*FileHash

	for item := range input {
		hashes = append(hashes, item)
	}
	sort.Stable(fileHashByPath(hashes))

	for _, item := range hashes {
		output <- item
	}
}

// Sprintf returns a formatted string representation of FileHash
func Sprintf(format string, item *FileHash) string {
	switch format {
	case "":
		return fmt.Sprintf("%x  %s\n", item.Hash, item.File.Path)
	default:
		return fmt.Sprintf(format, item.Hash, item.File.Path)
	}
}

// OutputTextFile goroutine
func OutputTextFile(input <-chan *FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()

	printFormat := viper.GetString("print-format")
	outputFileName := viper.GetString("output-file")

	numFiles := int64(0)
	totalSize := int64(0)
	totalTime := time.Duration(0)

	var outputFile *os.File
	switch outputFileName {
	case "-":
		outputFile = os.Stdout
	default:
		var err error
		outputFile, err = os.Create(outputFileName)
		if err != nil {
			log.Fatalf("error opening output file: %v\n", err)
		}
	}
	defer outputFile.Close()

	for item := range input {
		numFiles++
		totalSize += item.File.Size
		totalTime += item.File.ProcTime
		fmt.Fprintln(outputFile, Sprintf(printFormat, item))
	}

	if viper.GetBool("stats") {
		elapsedTime := time.Since(startTime)
		p := message.NewPrinter(language.English)
		p.Fprintf(os.Stderr,
			"# hashed %v bytes from %v files in %s (%s cpu over %v routines) => %.1f MB/s\n",
			totalSize,
			numFiles,
			elapsedTime.Truncate(time.Millisecond).String(),
			totalTime.Truncate(time.Millisecond).String(),
			viper.GetInt("parallel"),
			float64(totalSize)/elapsedTime.Seconds()/1000000)
	}
}

// OutputProtobufFile goroutine
func OutputProtobufFile(input <-chan *FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()

	outputFileName := viper.GetString("output-file")

	numFiles := int64(0)
	totalSize := int64(0)
	totalTime := time.Duration(0)

	var outputFile *os.File
	switch outputFileName {
	case "-":
		outputFile = os.Stdout
	default:
		var err error
		outputFile, err = os.Create(outputFileName)
		if err != nil {
			log.Fatalf("error opening output file: %v\n", err)
		}
	}
	defer outputFile.Close()

	pathHashMap := make(map[string]*bitratpb.HashData)
	recordSet := &bitratpb.RecordSet{
		Algorithm:   viper.GetString("hash"),
		PathHashMap: pathHashMap,
	}

	for item := range input {
		numFiles++
		totalSize += item.File.Size
		totalTime += item.File.ProcTime
		recordSet.PathHashMap[item.File.Path] = &bitratpb.HashData{
			Hash:    item.Hash,
			Size:    item.File.Size,
			ModTime: timestamppb.New(item.File.ModTime),
		}
	}

	if viper.GetBool("stats") {
		elapsedTime := time.Since(startTime)
		statistics := &bitratpb.Statistics{
			NumFiles:    numFiles,
			TotalBytes:  totalSize,
			ElapsedTime: durationpb.New(elapsedTime),
			TotalTime:   durationpb.New(totalTime),
			Parallel:    int32(viper.GetInt("parallel")),
		}
		recordSet.Statistics = statistics
	}

	out, err := proto.Marshal(recordSet)
	if err != nil {
		log.Fatalln("Failed to encode result set:", err)
	}
	// TODO: catch errors here
	outputFile.Write(out)
}

// HashRouter splits an input FileHash channel into separate outputs based upon whether the item Hash is nil
func HashRouter(input <-chan FileHash, hashOutput, nilOutput chan<- FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		switch item.Hash {
		case nil:
			nilOutput <- item
		default:
			hashOutput <- item
		}
	}
}

// HashSink blackholes items from a FileHash channel
func HashSink(input <-chan FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for range input {
		// blackhole input
	}
}
