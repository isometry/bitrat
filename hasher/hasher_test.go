package hasher

import (
	"sync"
	"testing"

	"github.com/isometry/bitrat/pathwalk"
)

const (
	testPath = "/usr/local"
	testHash = "blake2b-512"
)

func sinkFileInfoChan(input <-chan *FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for range input {
	}
}

func benchmarkHasherType1(options *pathwalk.Options) {
	var (
		pathWalkWaitGroup sync.WaitGroup
		hasherWaitGroup   sync.WaitGroup
		hashSinkWaitGroup sync.WaitGroup
	)

	fileChan := make(chan *pathwalk.File, 1024)
	hashChan := make(chan *FileHash, 64)

	walker := pathwalk.New(testPath, options, fileChan, &pathWalkWaitGroup)
	pathWalkWaitGroup.Add(1)
	go walker.Walk()

	for i := 0; i < options.Parallel; i++ {
		hasherWaitGroup.Add(1)
		h := New(testHash, []byte(""))
		go h.HashProcessor(fileChan, hashChan, &hasherWaitGroup)
	}

	hashSinkWaitGroup.Add(1)
	go sinkFileInfoChan(hashChan, &hashSinkWaitGroup)

	pathWalkWaitGroup.Wait()
	close(fileChan)
	hasherWaitGroup.Wait()
	close(hashChan)
	hashSinkWaitGroup.Wait()
}

func benchmarkHasherType2(options *pathwalk.Options) {
	var (
		pathWalkWaitGroup sync.WaitGroup
		hasherWaitGroup   sync.WaitGroup
		hashSinkWaitGroup sync.WaitGroup
	)

	fileChan := make(chan *pathwalk.File, 1024)
	hashChan := make(chan *FileHash, 64)

	walker := pathwalk.New2(testPath, options, fileChan, &pathWalkWaitGroup)
	pathWalkWaitGroup.Add(1)
	go walker.Walk()

	for i := 0; i < options.Parallel; i++ {
		hasherWaitGroup.Add(1)
		h := New(testHash, []byte(""))
		go h.HashProcessor(fileChan, hashChan, &hasherWaitGroup)
	}

	hashSinkWaitGroup.Add(1)
	go sinkFileInfoChan(hashChan, &hashSinkWaitGroup)

	pathWalkWaitGroup.Wait()
	close(fileChan)
	hasherWaitGroup.Wait()
	close(hashChan)
	hashSinkWaitGroup.Wait()
}

func BenchmarkHasherType1Parallel1(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 1,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}

func BenchmarkHasherType1Parallel2(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 2,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}
func BenchmarkHasherType1Parallel3(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 3,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}
func BenchmarkHasherType1Parallel4(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 4,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}

func BenchmarkHasherType1Parallel5(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 5,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}

func BenchmarkHasherType1Parallel6(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 6,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}

func BenchmarkHasherType1Parallel7(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 7,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}

func BenchmarkHasherType1Parallel8(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 8,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType1(&options)
	}
}

func BenchmarkHasherType2Parallel1(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 1,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}

func BenchmarkHasherType2Parallel2(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 2,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}
func BenchmarkHasherType2Parallel3(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 3,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}
func BenchmarkHasherType2Parallel4(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 4,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}

func BenchmarkHasherType2Parallel5(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 5,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}

func BenchmarkHasherType2Parallel6(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 6,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}

func BenchmarkHasherType2Parallel7(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 7,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}

func BenchmarkHasherType2Parallel8(b *testing.B) {
	options := pathwalk.Options{
		Recurse:  true,
		Parallel: 8,
	}

	for n := 0; n < b.N; n++ {
		benchmarkHasherType2(&options)
	}
}
