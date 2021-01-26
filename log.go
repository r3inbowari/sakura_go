package Sakura

import (
	"fmt"
	"time"
)

type Ext struct {
}

func Info(msg string, ext ...Ext) {
	t := time.Now().Format("2006-01-02 15:04:05")
	str := "[INFO] " + t + " " + msg
	fmt.Printf("\x1b[%dm"+str+" \x1b[0m\n", 32)
}

func Warn(msg string, ext ...Ext) {
	t := time.Now().Format("2006-01-02 15:04:05")
	str := "[WARN] " + t + " " + msg
	fmt.Printf("\x1b[%dm"+str+" \x1b[0m\n", 33)
}

func Fail(msg string, ext ...Ext) {
	t := time.Now().Format("2006-01-02 15:04:05")
	str := "[FAIL] " + t + " " + msg
	fmt.Printf("\x1b[%dm"+str+" \x1b[0m\n", 31)
}

func Blue(msg string) {
	fmt.Printf("\x1b[%dm"+msg+" \x1b[0m\n", 34)
}