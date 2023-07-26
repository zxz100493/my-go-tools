package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Visible     bool
	Target      string
	Destination string
}

var ch = make(chan struct{}, runtime.NumCPU())
var filePath = make(chan string)

func Find(cfg *Config) {
	fmt.Println("find file service", cfg.Visible)
	go readDir(cfg.Target, cfg.Destination)
	watcher()

}

func watcher() {
loop:
	for {
		select {
		case <-filePath:
			break loop
		}
	}
}

func readDir(t string, d string) string {
	fmt.Println("ssss")
	for _, entry := range dirents(d) {
		fmt.Printf("current dir is:%s\n", filepath.Join(d, entry.Name()))

		if entry.IsDir() {
			fmt.Printf("it is dir need read continue:%s\n", filepath.Join(d, entry.Name()))

			go readDir(filepath.Join(d, entry.Name()), d)
		} else {
			fmt.Printf("entry name is:%s\n", entry.Name())
			if entry.Name() == t {
				return filepath.Join(d, entry.Name())
			}
		}
	}
	return fmt.Sprintf("can't find %s in target path %s", t, d)
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v", dir, err)
		return nil
	}
	return entries
}
