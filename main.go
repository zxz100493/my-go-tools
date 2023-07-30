/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"my-go-tools/cmd"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	//远程获取pprof数据
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
	cmd.Execute()
	time.Sleep(60 * time.Second)
}
