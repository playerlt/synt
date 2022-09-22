package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	"github.com/playerlt/synt/config"
	"github.com/playerlt/synt/server"
)

func main() {
	chChromeDie := make(chan struct{})
	chBackendDie := make(chan struct{})
	chSignal := listenToInterrupt()
	go server.Run()
	go startBrowser(chChromeDie, chBackendDie)
	for {
		select {
		case <-chSignal:
			chBackendDie <- struct{}{}
		case <-chChromeDie:
			os.Exit(0)
		}
	}
}

func startBrowser(chChromeDie chan struct{}, chBackendDie chan struct{}) {
	// 先写死路径，后面再照着 lorca 改
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// chromePath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	//--disable-web-security
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+config.GetPort()+"/static/index.html")
	err := cmd.Start()
	if err != nil {
		fmt.Println("错误!!!")
	}
	// fmt.Println("进程id: ", cmd.Process.Pid)
	go func() {
		<-chBackendDie
		// kill(cmd)
	}()
	go func() {
		// cmd.Wait()
		// for !cmd. {
		// 	time.Sleep(5 * time.Second)
		// }
		time.Sleep(300 * time.Second)
		// err := cmd.Process.Kill()
		// if err != nil {
		// 	fmt.Println(err)

		// }
		chChromeDie <- struct{}{}
	}()
}
func kill(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}
func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
