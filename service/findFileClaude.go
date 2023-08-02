package service

/* package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type pathStack struct {
	paths []string
}

// Search searches for target file under root dir recursively.
// Returns the full path if found, empty string if not found.
func Search(rootDir, targetFile string) string {

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Use a stack to traverse directories non-recursively
	var stack pathStack
	stack.paths = []string{rootDir}

	// Worker pool limits the concurrency
	workerPool := make(chan struct{}, 10)

	// Found channel returns the result
	foundChan := make(chan string)

	go func() {
		for len(stack.paths) > 0 {
			select {
			case workerPool <- struct{}{}:
				path := stack.paths[len(stack.paths)-1]
				stack.paths = stack.paths[:len(stack.paths)-1]
				wg.Add(1)
				go func(path string) {
					defer wg.Done()
					doSearch(path, targetFile, foundChan, stack, workerPool)
				}(path)
			case <-foundChan:
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
	close(foundChan)
	return <-foundChan
}

func doSearch(root, target string, found chan string, stack *pathStack, workerPool chan struct{}) {
	defer func() {
		<-workerPool
	}()

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("read dir failed:", root, err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			stack.paths = append(stack.paths, filepath.Join(root, entry.Name()))
		} else {
			if entry.Name() == target {
				found <- filepath.Join(root, entry.Name())
				return
			}
		}
	}
}
*/
