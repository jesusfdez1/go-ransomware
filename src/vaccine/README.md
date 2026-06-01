# Vaccine Module

Recovery tool designed to revert the encryption applied by the ransomware using AES-256-GCM decryption. It safely and securely restores original files, verifying their integrity.

## Module Architecture

```
vaccine/
├── main.go                 # Entry point
├── config/
│   ├── config.go          # Decryption configuration
│   └── helpers.go         # Auxiliary functions
└── core/
    └── decryption.go      # AES-256-GCM decryption engine
```

## Vaccine Capabilities

- **File Decryption**: Reverts the AES-256-GCM encryption applied by the ransomware. It searches for `.file` extensions in configured directories, extracts necessary metadata (original extension, nonce), and restores the original content while verifying integrity via the GCM authentication tag.
- **Automatic Search**: Recursively explores all configured directories to locate encrypted files, ignoring system directories and files that do not match the expected format.
- **Integrity Verification**: Uses the AES-GCM authentication tag to verify the integrity of each decrypted file, detecting any corruption or manipulation that may have occurred.
- **Automatic Cleanup**: Deletes the encrypted `.file` archives only after successfully verifying the decryption process, preventing data loss from errors during recovery.

## Configuration

Edit `config/config.go` using the same values as the ransomware:

```go
const (
	EnableDebugMode = false                               // Enable debug messages
	EncryptionKey   = "my-secret-encryption-key-123"      // Decryption key
)

var TargetDirectories = []string{
	"C:\\Users\\",
}
```

**Important**: The encryption key must match exactly the one used by the ransomware. If the ransomware used a reverse shell with the `SET_KEY` command, you must use that specific key here.

## Usage
### Running the Vaccine

```powershell
cd vaccine
.\ejercicio2.exe
```

## Error Handling

- **Incorrect Key**: The GCM authentication tag fails. The `.file` remains intact without corruption.
- **Corrupted File**: The GCM verification detects modifications. It aborts the decryption for that file and continues with the rest.
- **Insufficient Permissions**: Reports the specific error and continues with accessible files.
- **Disk Space**: Verifies available space before writing. Aborts if insufficient.
- **Invalid Metadata**: Skips files with an invalid structure without attempting decryption.
