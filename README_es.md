# Go-Ransomware y vacuna
> [!CAUTION]
> **Aviso de seguridad**
> Por motivos de seguridad, algunas partes del código original han sido sustituidas por placeholders.

> [!NOTE]
> To read this document in English, visit this [file](readme.md)

Un proyecto completo en dos partes desarrollado en Go para el estudio de técnicas modernas de cifrado, evasión de defensas, comunicación de Comando y Control (C2) y mecanismos de recuperación.

## Arquitectura del proyecto

El sistema se organiza en dos módulos independientes que operan de manera complementaria:

```text
go-ransomware/
├── ransomware/              # Módulo de cifrado (Omitido)
│   ├── main.go             # Punto de entrada
│   ├── config/             # Configuración centralizada
│   ├── core/               # Motor de cifrado AES-GCM (Omitido)
│   └── utils/
│       ├── infostealer.go  # Recopilación de información (Omitido)
│       ├── obfuscation.go  # Técnicas de ofuscación de comandos
│       ├── persistence.go  # Mecanismos de persistencia (Omitido)
│       ├── reverse.go      # Shell inversa para C2 (Omitido)
│       └── shadowcopy.go   # Eliminación de Volume Shadow Copy (Omitido)
│
└── vaccine/                 # Módulo de recuperación
    ├── main.go             # Punto de entrada de la vacuna
    ├── config/             # Configuración de descifrado
    └── core/               # Motor de descifrado AES-GCM
```

## Características y conceptos técnicos (teóricos)

- **Cifrado avanzado (Omitido):** Implementación de AES-GCM (Galois/Counter Mode) para cifrado autenticado.
- **Evasión de defensas y ofuscación:** Técnicas para evadir protección moderna de endpoints, incluyendo ofuscación dinámica de cadenas (strings) y ejecución en ventana oculta.
- **Comunicación C2 (Omitido):** Integración de shell inversa persistente para ejecución remota de comandos e intercambio de claves.
- **Manipulación del sistema (Omitido):** Robo de información, establecimiento de persistencia a través del registro y tareas programadas, y destrucción de respaldos (VSS/WMIC).
- **Compilación automatizada:** Scripts de construcción en PowerShell con compresión UPX y stripping para reducir la huella del ejecutable.

## Demostraciones de evasión de antivirus

Los payloads fueron probadas con éxito contra soluciones AV del mercado.

![Bitdefender Bypass](vid/bitdefender.gif)

![ESET Bypass](vid/eset.gif)

## Compilación y uso

El proyecto incluye un script de PowerShell para automatizar el proceso de construcción de ambos módulos:

```powershell
.\build.ps1 [modo]
```

### Modos de compilación

**Modo 1 - Debug**: Compilación básica con símbolos de depuración completos.
```powershell
.\build.ps1 1
```

**Modo 2 - Producción** (Predeterminado): Construcción optimizada que incluye:
- Eliminación de símbolos de depuración (`-ldflags="-s -w"`)
- Recorte de rutas (`-trimpath`)
- Ejecución en ventana oculta (`-H windowsgui`)
- Compresión UPX con LZMA (si está instalado)

```powershell
.\build.ps1 2
# O simplemente
.\build.ps1
```

Para instalar UPX en Windows:
```powershell
winget install upx.upx
```

## Documentación técnica

Para información específica sobre cada módulo:
- **Ransomware**: Consulta `ransomware/README.md` para detalles técnicos sobre el motor de cifrado, la persistencia y la comunicación C2.
- **Vacuna**: Consulta `vaccine/README.md` para especificaciones sobre el proceso de recuperación.
