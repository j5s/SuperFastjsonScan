package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	host   string
	port   int
	target string
	help   bool
)

var (
	socketChan chan bool
	reqChan    chan bool
)

func main() {
	printLogo()
	parserInput()
	socketChan = make(chan bool)
	reqChan = make(chan bool)
	go startListen(host, port)
	go scan()
	for {
		select {
		case <-socketChan:
			Info("find fastjson")
		case <-reqChan:
			Info("scan finish")
			os.Exit(-1)
		}
	}

}

func printLogo() {
	fmt.Println("   _____                       _____             " +
		"    \n  / ____|                     / ____|                \n " +
		"| (___  _   _ _ __   ___ _ _| (___   ___ __ _ _ __  \n  \\___ " +
		"\\| | | | '_ \\ / _ \\ '__\\___ \\ / __/ _` | '_ \\ \n  ____) |" +
		" |_| | |_) |  __/ |  ____) | (_| (_| | | | |\n |_____/ \\__,_| " +
		".__/ \\___|_| |_____/ \\___\\__,_|_| |_|\n              | |    " +
		"                                \n              |_|            " +
		"                        ")
	fmt.Println("demo version by 4ra1n")
}

func scan() {
	payload := getPayload(host, port)
	resp := doRequest(&Request{Url: target, Method: "POST", Body: payload})
	if resp != nil {
		reqChan <- true
	}
}

func parserInput() {
	flag.StringVar(&target, "u", "", "scan target")
	flag.StringVar(&host, "h", "127.0.0.1", "your ip address")
	flag.IntVar(&port, "p", 8888, "which port you want to open")
	flag.BoolVar(&help, "help", false, "help info")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	if target == "" {
		Error("no target")
	}
}
