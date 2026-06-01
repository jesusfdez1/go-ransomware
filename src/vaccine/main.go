package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"vaccine/config"
	"vaccine/core"
)

func main() {
	// Confirmación en modo debug
	if config.EnableDebugMode {
		fmt.Println("=== VACUNA - DESENCRIPTACIÓN ===")
		fmt.Print("¿Continuar? (s/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(response))[0] != 's' {
			return
		}
	}

	// Iniciar proceso de descifrado
	config.DebugPrintln("Iniciando descifrado...")
	err := core.StartDecryption()
	if err != nil {
		if config.EnableDebugMode {
			fmt.Printf("Error: %v\n", err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
		return
	}
	config.DebugPrintln("Proceso completado")

	if config.EnableDebugMode {
		fmt.Println("\nPresiona Enter para cerrar...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
