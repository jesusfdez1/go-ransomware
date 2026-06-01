package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vaccine/config"
)

// GenerateKeyFromString deriva una clave AES de 32 bytes desde un string
func GenerateKeyFromString(keyString string) []byte {
	hash := sha256.Sum256([]byte(keyString))
	return hash[:]
}

// createGCM crea un cipher GCM a partir de una clave
func createGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creando cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creando GCM: %v", err)
	}

	return gcm, nil
}

// DecryptFile descifra un archivo individual
func DecryptFile(encryptedFilePath string) error {
	// Leer archivo cifrado
	data, err := os.ReadFile(encryptedFilePath)
	if err != nil {
		return fmt.Errorf("error leyendo archivo cifrado: %v", err)
	}

	// Crear GCM cipher
	key := GenerateKeyFromString(config.EncryptionKey)
	gcm, err := createGCM(key)
	if err != nil {
		return err
	}

	// Validar y extraer nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return fmt.Errorf("datos cifrados inválidos")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Descifrar datos (que incluyen extensión + contenido)
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("error descifrando: %v", err)
	}

	// Extraer nombre y extensión originales de los datos descifrados
	if len(decryptedData) < 3 {
		return fmt.Errorf("datos descifrados inválidos")
	}

	// Leer longitud del nombre (2 bytes: alto, bajo)
	nameLen := (int(decryptedData[0]) << 8) | int(decryptedData[1])

	if len(decryptedData) < 3+nameLen {
		return fmt.Errorf("datos descifrados corruptos: nombre inválido")
	}

	// Extraer nombre original (sin extensión)
	originalName := string(decryptedData[2 : 2+nameLen])

	// Leer longitud de extensión
	extLen := int(decryptedData[2+nameLen])

	if len(decryptedData) < 3+nameLen+extLen {
		return fmt.Errorf("datos descifrados corruptos: extensión inválida")
	}

	// Extraer extensión original
	originalExt := string(decryptedData[3+nameLen : 3+nameLen+extLen])

	// Extraer datos del archivo
	originalData := decryptedData[3+nameLen+extLen:]

	// Restaurar archivo con nombre y extensión originales
	dir := filepath.Dir(encryptedFilePath)
	originalFilePath := filepath.Join(dir, originalName+originalExt)

	if err := os.WriteFile(originalFilePath, originalData, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %v", err)
	}

	if err := os.Remove(encryptedFilePath); err != nil {
		return fmt.Errorf("error eliminando archivo cifrado: %v", err)
	}

	return nil
}

// findEncryptedFiles busca archivos .file en los directorios objetivo
func findEncryptedFiles() []string {
	var files []string
	for _, dir := range config.TargetDirectories {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".file") {
				files = append(files, path)
			}
			return nil
		})
	}
	return files
}

// StartDecryption inicia el proceso de descifrado
func StartDecryption() error {
	files := findEncryptedFiles()
	config.DebugPrint("Archivos encontrados: %d\n", len(files))

	if len(files) == 0 {
		config.DebugPrintln("No se encontraron archivos cifrados")
		config.DebugPrintln("Directorios configurados:")
		for _, dir := range config.TargetDirectories {
			config.DebugPrint("  - %s\n", dir)
		}
		return nil
	}

	// Intentar descifrar cada archivo
	successCount := 0
	for _, file := range files {
		if DecryptFile(file) == nil {
			successCount++
		}
	}

	config.DebugPrint("Descifrados: %d/%d\n", successCount, len(files))

	// Eliminar archivos de rescate generados por el ransomware
	deleteRansomNotes()

	return nil
}

// deleteRansomNotes elimina los archivos de rescate generados por el ransomware
func deleteRansomNotes() {
	locations := []string{
		"C:\\Users\\Public\\Desktop\\NOTE.txt",
	}

	for _, location := range locations {
		if err := os.Remove(location); err != nil {
			// No reportar error si el archivo no existe
			if !os.IsNotExist(err) {
				config.DebugPrint("No se pudo eliminar %s: %v\n", location, err)
			}
		} else {
			config.DebugPrint("Archivo eliminado: %s\n", location)
		}
	}
}
