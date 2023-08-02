package service

import (
	"fmt"
	"io/fs"
	"log"
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
	maxCurrencyNum = runtime.NumCPU() / 2
	dirContainer   = make(chan string, maxCurrencyNum)
	fileContainer  = make(chan string, maxCurrencyNum)
	// filePath       = make(chan string)
	// done           = make(chan struct{})
	// fileNum        = make(chan int)
	sema = make(chan struct{}, 20)
)
var tick <-chan time.Time

var n sync.WaitGroup

func Find(cfg *Config) {
	defer trace("slow")()
	fmt.Println("find file service", cfg.Visible)
	fmt.Println(runtime.NumCPU())

	// 读取目录放入通道
	go readDir2(cfg)
	// select 监听通道 如果已经找到  就stop

	// n.Add(1)
	// go readDir(cfg.TargetFile, cfg.DstDir, &n)

	if cfg.Visible {
		tick = time.Tick(10000 * time.Millisecond)
	}

	// go func() {
	// 	n.Wait()
	// 	fmt.Println("close")
	// 	close(filePath)
	// }()
	/*
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
	   	fmt.Println("watcher end") */
	fmt.Println("watcher ...")

	watcher(cfg)
	fmt.Println("watcher over")
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", "msg")
	return func() {
		log.Printf("exit %s(%s)", msg, time.Since(start))
	}
}

func readDir(file string, dir string, n *sync.WaitGroup) string {
	defer n.Done()
	// fmt.Println("ssss")
	for _, entry := range dirents(dir) {
		// fmt.Printf("current dir is:%s\n", filepath.Join(dir, entry.Name()))
		if entry.IsDir() {
			// fmt.Printf("it is dir need read continue:%s\n", filepath.Join(dir, entry.Name()))
			n.Add(1)
			// sema <- struct{}{}
			subdir := filepath.Join(dir, entry.Name())
			/* go func(subdir string) {
				readDir(file, subdir, n)
				<-sema
			}(subdir) */
			readDir(file, subdir, n)
		} else {
			fmt.Println("file---")

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
	/* sema <- struct{}{}
	defer func() {
		<-sema
	}()
	fmt.Println("dirents", len(sema), cap(sema)) */

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

func watcher(cfg *Config) {
	var num int
loop:
	for {
		select {
		case _, ok := <-dirContainer:
			if !ok {
				fmt.Println("not ok dircontainer")
				break loop
			}

			// dispatch(subdir)
			// fmt.Println("i amd select filenum")
		case <-tick:
			fmt.Println("\n num is:", num)
		case fileName, ok := <-fileContainer:
			if !ok {
				fmt.Println("not ok filecontainer")
				fmt.Println("break")
				// filePath channel is closed
				break loop // exit the loop when the channel is closed
			}
			// fmt.Printf("file path is:%s\n", fileName)
			// fmt.Println(len(fileContainer))
			// fmt.Println(len(dirContainer))
			info, _ := os.Stat(fileName)

			if cfg.TargetFile == info.Name() {
				fmt.Printf("xxxxxx--yes it is---xxxxxx:%s\n", fileName)
				break loop
			}
		}
	}
	fmt.Println("watcher end")
}

// 分流文件和目录
func dispatch(dir string) {
	var tmp = make([]string, 0)
	for _, v := range dirents(dir) {
		// fmt.Println("receive the dir:", v)

		if v.IsDir() {
			subdir := filepath.Join(dir, v.Name())
			// fmt.Println(subdir)
			// fmt.Println("lenth dirContainer", len(dirContainer), cap(dirContainer))
			tmp = append(tmp, subdir)
		} else {
			info, err := os.Stat(filepath.Join(dir, v.Name()))
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			if !hasReadPermission(info) {
				continue
			}
			// fmt.Println("lenth fileContainer", len(fileContainer), cap(fileContainer))
			name := filepath.Join(dir, v.Name())
			// fmt.Println(name)

			fileContainer <- name
		}
	}
	if len(tmp) > 0 {
		for _, v := range tmp {
			dispatch(v)
		}
	}
}

func readDir2(cfg *Config) {
	for _, entry := range dirents(cfg.DstDir) {
		if entry.IsDir() {
			n.Add(1)
			sema <- struct{}{}
			go func(dir string, n *sync.WaitGroup) {
				dispatch(dir)
				<-sema
				n.Done()
			}(filepath.Join(cfg.DstDir, entry.Name()), &n)
			n.Wait()
		}
	}
}
