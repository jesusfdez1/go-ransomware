package config

import "fmt"

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
