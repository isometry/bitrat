package pathwalk

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// PathWalker interface
type PathWalker interface {
	Walk()
}

// Options for PathWalker implementations
type Options struct {
	Pattern     string
	glob        bool
	Recurse     bool
	Parallel    int
	HiddenDirs  bool
	HiddenFiles bool
	IncludeGit  bool
}

// File object returned by Walker.Walk method
type File struct {
	Path     string
	Size     int64
	ModTime  time.Time
	ProcTime time.Duration
	Error    error
}

// Walker object
type Walker struct {
	Root    string
	Options *Options
	Output  chan<- *File
}

// NewWalker returns a new Walker
func NewWalker(root string, options *Options, output chan<- *File) PathWalker {
	if options.Pattern != "" {
		options.glob = true
		if _, err := filepath.Match(options.Pattern, ""); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Clean(os.Args[0]), err)
			os.Exit(1)
		}
	}

	return &Walker{
		Root:    root,
		Options: options,
		Output:  output,
	}
}

// Walk the path
func (p *Walker) Walk() {
	err := filepath.WalkDir(p.Root, p.step)
	if err != nil {
		panic(err)
	}
}

func (p *Walker) step(path string, d fs.DirEntry, err error) error {
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Clean(os.Args[0]), err)
		return nil
	case d.IsDir():
		if path != p.Root {
			switch {
			case !p.Options.Recurse:
				return filepath.SkipDir
			case !p.Options.IncludeGit && d.Name() == ".git":
				return filepath.SkipDir
			case p.Options.IncludeGit && d.Name() == ".git":
				break
			case !p.Options.HiddenDirs && strings.HasPrefix(d.Name(), "."):
				return filepath.SkipDir
			}
		}
		return nil
	case !d.Type().IsRegular():
		return nil
	case !p.Options.HiddenFiles && strings.HasPrefix(d.Name(), "."):
		return nil
	case p.Options.glob:
		if globMatch, _ := filepath.Match(p.Options.Pattern, d.Name()); !globMatch {
			return nil
		}
	}

	info, err := d.Info()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Clean(os.Args[0]), err)
		return nil
	}

	p.Output <- &File{Path: path, Size: info.Size(), ModTime: info.ModTime()}

	return nil
}

// AltWalker is an alternate walker that walks the path concurrently without
// opening too many simultaneous files
type AltWalker struct {
	Root    string
	Options *Options
	Sync    chan bool
	Output  chan<- *File
}

// NewAltWalker is the constructor for the AltWalker type
func NewAltWalker(root string, options *Options, output chan<- *File) PathWalker {
	if options.Pattern != "" {
		options.glob = true
		if _, err := filepath.Match(options.Pattern, ""); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Clean(os.Args[0]), err)
			os.Exit(1)
		}
	}

	sync := make(chan bool, options.Parallel)
	return &AltWalker{
		Root:    root,
		Options: options,
		Sync:    sync,
		Output:  output,
	}
}

// Walk walks the path with goroutine per directory
func (p *AltWalker) Walk() {
	var wg sync.WaitGroup
	wg.Go(func() { p.step(&wg, p.Root) })
	wg.Wait()
}

func (p *AltWalker) step(wg *sync.WaitGroup, path string) {
	p.Sync <- true

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "# ERROR: %s: %v\n", path, err)
	}

	for _, file := range files {
		switch {
		case file.IsDir():
			switch {
			case !p.Options.Recurse:
				continue
			case !p.Options.IncludeGit && file.Name() == ".git":
				continue
			case p.Options.IncludeGit && file.Name() == ".git":
				break
			case !p.Options.HiddenDirs && strings.HasPrefix(file.Name(), "."):
				continue
			}
			dirPath := filepath.Join(path, file.Name())
			wg.Go(func() { p.step(wg, dirPath) })
			continue
		case !file.Type().IsRegular():
			continue
		case !p.Options.HiddenFiles && strings.HasPrefix(file.Name(), "."):
			continue
		case p.Options.glob:
			if globMatch, _ := filepath.Match(p.Options.Pattern, file.Name()); !globMatch {
				continue
			}
		}

		info, err := file.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "# ERROR: %s: %v\n", filepath.Join(path, file.Name()), err)
			continue
		}

		p.Output <- &File{
			Path:    filepath.Join(path, file.Name()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}
	}

	<-p.Sync
}
