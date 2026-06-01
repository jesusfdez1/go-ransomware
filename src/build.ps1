# Script de compilacion
# Uso: .\build.ps1 [modo]
# Modos: 1=Debug (desarrollo), 2=Produccion (UPX+optimizaciones)

param([int]$Mode = 2)

$workspace = $PSScriptRoot
if ($Mode -lt 1 -or $Mode -gt 2) { $Mode = 2 }

$modeName = @{1="Debug"; 2="Produccion"}[$Mode]
Write-Host "[*] Modo: $modeName" -ForegroundColor Cyan
Write-Host ""

# Compilar Ransomware
Write-Host "[1/2] Compilando ransomware..." -ForegroundColor Green
Set-Location "$workspace\ransomware"

$success = $false
switch ($Mode) {
    1 {
        # Modo Debug: compilacion basica para desarrollo
        Write-Host "  [*] Compilacion debug..." -ForegroundColor White
        go build -o ejercicio1.exe 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  [+] Compilacion completada" -ForegroundColor Green
            $success = $true
        }
    }
    2 {
        # Modo Produccion: optimizaciones + UPX
        Write-Host "  [*] Compilando con optimizaciones..." -ForegroundColor White
        
        go build -ldflags="-s -w -H windowsgui" -trimpath -o ejercicio1.exe 2>&1 | Out-Null
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  [+] Compilación completada" -ForegroundColor Green
            
            # Aplicar UPX para comprimir
            if (Get-Command upx -ErrorAction SilentlyContinue) {
                Write-Host "  [*] Comprimiendo con UPX..." -ForegroundColor White
                upx --best --lzma ejercicio1.exe 2>&1 | Out-Null
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "  [+] Comprimido con UPX" -ForegroundColor Green
                } else {
                    Write-Host "  [!] UPX fallo (continuando...)" -ForegroundColor Yellow
                }
            } else {
                Write-Host "  [!] UPX no instalado (sin comprimir)" -ForegroundColor Yellow
            }
            $success = $true
        }
    }
}

if ($success -and (Test-Path "ejercicio1.exe")) {
    $sizeBytes = (Get-Item "ejercicio1.exe").Length
    $sizeMB = [math]::Round($sizeBytes / 1MB, 2)
    $sizeKB = [math]::Round($sizeBytes / 1KB, 0)
    Write-Host "  [OK] Ransomware: $sizeMB MB ($sizeKB KB)" -ForegroundColor Green
    $ransomwareOK = $true
} else {
    Write-Host "  [X] Error en ransomware" -ForegroundColor Red
    $ransomwareOK = $false
}

# Compilar Vacuna
if ($ransomwareOK) {
    Write-Host "`n[2/2] Compilando vacuna..." -ForegroundColor Green
    Set-Location "$workspace\vaccine"
    
    $success = $false
    switch ($Mode) {
        1 {
            # Modo Debug: compilacion basica para desarrollo
            Write-Host "  [*] Compilacion debug..." -ForegroundColor White
            go build -o ejercicio2.exe 2>&1 | Out-Null
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  [+] Compilacion completada" -ForegroundColor Green
                $success = $true
            }
        }
        2 {
            # Modo Produccion: optimizaciones + UPX
            Write-Host "  [*] Compilando con optimizaciones..." -ForegroundColor White
            
            go build -ldflags="-s -w -H windowsgui" -trimpath -o ejercicio2.exe 2>&1 | Out-Null
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  [+] Compilación completada" -ForegroundColor Green
                
                # Aplicar UPX para comprimir
                if (Get-Command upx -ErrorAction SilentlyContinue) {
                    Write-Host "  [*] Comprimiendo con UPX..." -ForegroundColor White
                    upx --best --lzma ejercicio2.exe 2>&1 | Out-Null
                    if ($LASTEXITCODE -eq 0) {
                        Write-Host "  [+] Comprimido con UPX" -ForegroundColor Green
                    } else {
                        Write-Host "  [!] UPX fallo (continuando...)" -ForegroundColor Yellow
                    }
                } else {
                    Write-Host "  [!] UPX no instalado (sin comprimir)" -ForegroundColor Yellow
                }
                $success = $true
            }
        }
    }
    
    if ($success -and (Test-Path "ejercicio2.exe")) {
        $sizeBytes = (Get-Item "ejercicio2.exe").Length
        $sizeMB = [math]::Round($sizeBytes / 1MB, 2)
        $sizeKB = [math]::Round($sizeBytes / 1KB, 0)
        Write-Host "  [OK] Vacuna: $sizeMB MB ($sizeKB KB)" -ForegroundColor Green
        $vaccineOK = $true
    } else {
        Write-Host "  [X] Error en vacuna" -ForegroundColor Red
        $vaccineOK = $false
    }
} else {
    Write-Host "`n[!] Vacuna omitida por error en ransomware" -ForegroundColor Yellow
    $vaccineOK = $false
}

Set-Location $workspace


Write-Host ""

if ($ransomwareOK -and $vaccineOK) {
    Write-Host "[OK] Ambos ejecutables compilados" -ForegroundColor Cyan
} elseif ($ransomwareOK) {
    Write-Host "[WARN] Solo ransomware compilado" -ForegroundColor Yellow
} else {
    Write-Host "[ERROR] Fallo en compilacion" -ForegroundColor Red
}


