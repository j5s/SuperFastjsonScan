package main

import (
	"fmt"
)

func Info(str string) {
	fmt.Printf("[+] %s\n", str)
}

func Error(str string) {
	fmt.Printf("[!] %s\n", str)
}
