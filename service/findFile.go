package service

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

var (
	sema     = make(chan struct{}, runtime.NumCPU()/1)
	filePath = make(chan string)
	done     = make(chan struct{})
	fileNum  = make(chan int)
)

var n sync.WaitGroup

func Find(cfg *Config) {
	fmt.Println("find file service", cfg.Visible)
	fmt.Println(runtime.NumCPU())

	n.Add(1)
	go readDir(cfg.TargetFile, cfg.DstDir, &n)

	var tick <-chan time.Time

	if cfg.Visible {
		tick = time.Tick(1000 * time.Millisecond)
	}

	go func() {
		n.Wait()
		fmt.Println("close")
		close(filePath)
	}()

	var num int
loop:
	for {
		select {
		case size, ok := <-fileNum:
			if !ok {
				fmt.Println("not ok")
				break loop
			}
			// fmt.Println("i amd select filenum")
			num += size
		case <-tick:
			fmt.Println("\n num is:", num)
		case s, ok := <-filePath:
			if !ok {
				// filePath channel is closed
				break loop // exit the loop when the channel is closed
			}
			fmt.Printf("file path is:%s\n", s)
			fmt.Println("break")
			break loop
		}
	}
	fmt.Println("watcher end")
}

func readDir(file string, dir string, n *sync.WaitGroup) string {
	defer n.Done()
	// fmt.Println("ssss")
	for _, entry := range dirents(dir) {
		// fmt.Printf("current dir is:%s\n", filepath.Join(dir, entry.Name()))
		if entry.IsDir() {
			// fmt.Printf("it is dir need read continue:%s\n", filepath.Join(dir, entry.Name()))
			n.Add(1)
			sema <- struct{}{}
			subdir := filepath.Join(dir, entry.Name())
			go func(subdir string) {
				readDir(file, subdir, n)
				<-sema
			}(subdir)
		} else {
			fmt.Println("file does not have readable permission")

			info, err := os.Stat(filepath.Join(dir, entry.Name()))
			if err != nil {
				// fmt.Printf("error getting file info: %v", err)
				continue
			}
			if !info.Mode().IsRegular() {
				continue // 忽略非普通文件（如软链接）
			}
			if !hasReadPermission(info) {
				// fmt.Println("file does not have readable permission")
				continue
			}
			fmt.Println("send ", entry.Name())
			fileNum <- 1

			// fmt.Printf("entry name is:%s file name is:%s\n", entry.Name(), file)
			if entry.Name() == file {
				fmt.Println("---------Yes it is-----------")
				filePath <- filepath.Join(dir, entry.Name())
				done <- struct{}{}
				n.Done()
				return filepath.Join(dir, entry.Name())
			}
		}
	}
	// fmt.Printf("\ncan't find %s in target path %s\n", file, dir)
	return "nil"
}

func hasReadPermission(fi os.FileInfo) bool {
	mode := fi.Mode()
	return mode.IsRegular() || mode&os.ModeSymlink == os.ModeSymlink || mode&0400 == 0400
}

func dirents(dir string) []fs.DirEntry {

	/* defer func() {
		<-sema
	}()

	select {
	case sema <- struct{}{}:
	case <-done:
		fmt.Println("----done------")
		return nil
	} */

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v", dir, err)
		return nil
	}
	return entries
}
