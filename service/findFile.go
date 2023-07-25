package service

import (
	"fmt"
	"runtime"
)

type Config struct {
	Visible     bool
	Target      string
	Destination string
}

var token chan struct{}

func Find(c *Config) {
	fmt.Println("find file service", c.Visible)
	fmt.Println("token", token)
	num := runtime.NumCPU()
	ch := make(chan struct{}, num)
	fmt.Println("ch", ch)

	readDir()
}

func readDir(ch chan) {
loop:	
	for {
		select {
			case token <- ch:
			break loop;	
		}
	}
}
