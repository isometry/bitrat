package pathwalk

import (
	"sync"
	"testing"
)

func sinkFileChan(input <-chan *File, wg *sync.WaitGroup) {
	defer wg.Done()

	for range input {
	}
}

func BenchmarkPathwalker(b *testing.B) {
	options := Options{
		Recurse:     true,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		var (
			pathWalkWaitGroup  sync.WaitGroup
			pathPrintWaitGroup sync.WaitGroup
		)

		fileChan := make(chan *File, 64)

		walker := NewWalker("/usr/share", &options, fileChan, &pathWalkWaitGroup)
		pathWalkWaitGroup.Add(1)
		go walker.Walk()

		pathPrintWaitGroup.Add(1)
		go sinkFileChan(fileChan, &pathPrintWaitGroup)

		pathWalkWaitGroup.Wait()
		close(fileChan)
		pathPrintWaitGroup.Wait()
	}
}

func benchmarkPathwalker2(options *Options) {
	var (
		pathWalkWaitGroup  sync.WaitGroup
		pathPrintWaitGroup sync.WaitGroup
	)

	fileChan := make(chan *File, 64)

	walker := NewAltWalker("/usr/share", options, fileChan, &pathWalkWaitGroup)
	pathWalkWaitGroup.Add(1)
	go walker.Walk()

	pathPrintWaitGroup.Add(1)
	go sinkFileChan(fileChan, &pathPrintWaitGroup)

	pathWalkWaitGroup.Wait()
	close(fileChan)
	pathPrintWaitGroup.Wait()
}

func BenchmarkPathwalker2J1(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    1,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}

func BenchmarkPathwalker2J2(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    2,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}
func BenchmarkPathwalker2J3(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    3,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}
func BenchmarkPathwalker2J4(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    4,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}

func BenchmarkPathwalker2J5(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    5,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}

func BenchmarkPathwalker2J6(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    6,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}

func BenchmarkPathwalker2J7(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    7,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}

func BenchmarkPathwalker2J8(b *testing.B) {
	options := Options{
		Recurse:     true,
		Parallel:    8,
		HiddenDirs:  true,
		HiddenFiles: true,
		IncludeGit:  true,
	}

	for n := 0; n < b.N; n++ {
		benchmarkPathwalker2(&options)
	}
}
