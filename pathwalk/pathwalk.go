package pathwalk

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// File object returned by Walker.Walk method
type File struct {
	Path     string
	Size     int64
	ProcTime time.Duration
	Error    error
}

// Walker object
type Walker struct {
	Root    string
	Options *Options
	Output  chan<- *File
	Wait    *sync.WaitGroup
}

// Options for Walker.New controlling operation of the Walker
type Options struct {
	Pattern     string
	glob        bool
	Recurse     bool
	Parallel    int
	HiddenDirs  bool
	HiddenFiles bool
	SkipGit     bool
}

// New returns a new Walker
func New(root string, options *Options, output chan<- *File, wg *sync.WaitGroup) Walker {
	if options.Pattern != "" {
		options.glob = true
		if _, err := filepath.Match(options.Pattern, ""); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Clean(os.Args[0]), err)
			os.Exit(1)
		}
	}

	return Walker{
		Root:    root,
		Options: options,
		Output:  output,
		Wait:    wg,
	}
}

// Walk the path
func (p *Walker) Walk() {
	defer p.Wait.Done()

	err := filepath.Walk(p.Root, p.step)
	if err != nil {
		panic(err)
	}
}

func (p *Walker) step(path string, info os.FileInfo, err error) error {
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Clean(os.Args[0]), err)
		return nil
	case info.IsDir():
		if path != p.Root &&
			((!p.Options.Recurse) ||
				(!p.Options.HiddenDirs && strings.HasPrefix(info.Name(), ".")) ||
				(p.Options.SkipGit && info.Name() == ".git")) {
			return filepath.SkipDir
		}
		return nil
	case !info.Mode().IsRegular():
		return nil
	case !p.Options.HiddenFiles && strings.HasPrefix(info.Name(), "."):
		return nil
	case p.Options.glob:
		if globMatch, _ := filepath.Match(p.Options.Pattern, info.Name()); !globMatch {
			return nil
		}
	}

	p.Output <- &File{Path: path, Size: info.Size()}

	return nil
}

// Walker2 is an alternate walker that walks the path concurrently without
// opening too many simultaneous files
type Walker2 struct {
	Root    string
	Options *Options
	Sync    chan bool
	Output  chan<- *File
	Wait    *sync.WaitGroup
}

// New2 is the constructor for the Walker2 type
func New2(root string, options *Options, output chan<- *File, wg *sync.WaitGroup) Walker2 {
	sync := make(chan bool, options.Parallel)
	return Walker2{
		Root:    root,
		Options: options,
		Sync:    sync,
		Output:  output,
		Wait:    wg,
	}
}

// Walk walks the path with goroutine per directory
func (p *Walker2) Walk() {
	defer p.Wait.Done()

	p.Wait.Add(1)
	go p.step(p.Root)
}

func (p *Walker2) step(path string) {
	defer p.Wait.Done()

	p.Sync <- true
	//fmt.Fprintf(os.Stderr, "> Entered WalkPath('%s')\n", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "# ERROR: %s: %v\n", path, err)
	}

	for _, file := range files {
		switch {
		case file.IsDir() && (path == p.Root || p.Options.Recurse):
			switch {
			case !p.Options.HiddenDirs && strings.HasPrefix(file.Name(), "."):
				continue
			case !p.Options.SkipGit && file.Name() == ".git":
				continue
			default:
				p.Wait.Add(1)
				go p.step(filepath.Join(path, file.Name()))
			}
		case !file.Mode().IsRegular():
			continue
		case !p.Options.HiddenFiles && strings.HasPrefix(file.Name(), "."):
			continue
		case file.Mode().IsRegular():
			p.Output <- &File{
				Path: filepath.Join(path, file.Name()),
				Size: file.Size(),
			}
		}
	}

	//fmt.Fprintf(os.Stderr, "< Finished WalkPath('%s')\n", path)
	<-p.Sync
}
