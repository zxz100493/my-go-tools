package service

/* package service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Config struct {
	Visible    bool
	TargetFile string
	DstDir     string
}

type FileSearchService struct {
	cfg      *Config
	sema     chan struct{}
	filePath chan string
	fileNum  chan int
}

func NewFileSearchService(cfg *Config) *FileSearchService {
	// Limit the number of concurrently running goroutines
	maxConcurrentGoroutines := runtime.NumCPU() * 3
	return &FileSearchService{
		cfg:      cfg,
		sema:     make(chan struct{}, maxConcurrentGoroutines),
		filePath: make(chan string),
		fileNum:  make(chan int),
	}
}

func (s *FileSearchService) Find() {
	fmt.Println("find file service", s.cfg.Visible)
	fmt.Println(runtime.NumCPU())

	var tick <-chan time.Time
	if s.cfg.Visible {
		tick = time.Tick(1000 * time.Millisecond)
	}

	var num int
	n := sync.WaitGroup{}
	n.Add(1)
	// s.readDir(s.cfg.TargetFile, s.cfg.DstDir, &n)
	go func() {
		defer n.Done()
		s.readDir(s.cfg.TargetFile, s.cfg.DstDir, &n)
	}()
	go func() {
		n.Wait()
		fmt.Println("close")
		close(s.filePath)
	}()

loop:
	for {
		select {
		case size, ok := <-s.fileNum:
			if !ok {
				fmt.Println("not ok")
				break loop
			}
			num += size
		case <-tick:
			fmt.Println("\n num is:", num)
		case path, ok := <-s.filePath:
			if !ok {
				// s.filePath channel is closed
				break loop // exit the loop when the channel is closed
			}
			fmt.Printf("file path is:%s\n", path)
			fmt.Println("break")
			break loop
		}
	}
	fmt.Println("watcher end")
}

func (s *FileSearchService) readDir(file string, dir string, n *sync.WaitGroup) string {
	defer n.Done()

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v", dir, err)
		return "nil"
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println("dir", len(s.sema))
			n.Add(1)
			s.sema <- struct{}{}
			subdir := filepath.Join(dir, entry.Name())
			go func(subdir string) {
				s.readDir(file, subdir, n)
				<-s.sema
			}(subdir)
		} else {
			info, err := entry.Info()
			if err != nil {
				fmt.Printf("error getting file info: %v", err)
				continue
			}
			if !info.Mode().IsRegular() {
				continue // Ignore non-regular files (e.g., symbolic links)
			}
			if !hasReadPermission(info) {
				continue // Ignore files without read permission
			}

			fmt.Println("send ", entry.Name())
			s.fileNum <- 1

			if entry.Name() == file {
				fmt.Println("---------Yes it is-----------")
				s.filePath <- filepath.Join(dir, entry.Name())
				return filepath.Join(dir, entry.Name())
			}
		}
	}

	return "nil"
}

func hasReadPermission(fi fs.FileInfo) bool {
	mode := fi.Mode()
	return mode.IsRegular() || mode&os.ModeSymlink == os.ModeSymlink || mode&0400 == 0400
}
*/
