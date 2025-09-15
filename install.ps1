# GitX Windows Installer Script
# Usage: Invoke-WebRequest -Uri "https://raw.githubusercontent.com/gitxtui/gitx/master/install.ps1" -OutFile "install.ps1"; .\install.ps1

$repo = "gitxtui/gitx"
$installDir = "$env:ProgramFiles\gitx"

function Get-Arch {
    switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { Write-Error "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)"; exit 1 }
    }
}

function Get-LatestRelease {
    $apiUrl = "https://api.github.com/repos/$repo/releases"
    $response = Invoke-RestMethod -Uri $apiUrl
    if ($response -is [System.Collections.IEnumerable] -and $response.Count -gt 0) {
        return $response[0].tag_name
    } else {
        Write-Error "Could not find any release version for $repo."
        exit 1
    }
}

function Main {
    $os = "windows"
    $arch = Get-Arch
    $version = Get-LatestRelease
    $versionNum = $version -replace "^v", ""
    $filename = "gitx_${versionNum}_${os}_${arch}.zip"
    $downloadUrl = "https://github.com/$repo/releases/download/$version/$filename"

    Write-Host "Downloading gitx $version for $os/$arch..."
    $tempDir = New-TemporaryFile | Split-Path
    $zipPath = Join-Path $tempDir $filename
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath

    Write-Host "Extracting..."
    Expand-Archive -Path $zipPath -DestinationPath $tempDir

    if (!(Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir | Out-Null
    }

    $gitxExe = Get-ChildItem -Path $tempDir -Filter "gitx.exe" -Recurse | Select-Object -First 1
    if ($null -eq $gitxExe) {
        Write-Error "gitx.exe not found in the archive."
        exit 1
    }

    Copy-Item -Path $gitxExe.FullName -Destination (Join-Path $installDir "gitx.exe") -Force

    Write-Host "`ngitx has been installed to $installDir"
    Write-Host "Add $installDir to your PATH if not already present."
    Write-Host "Run 'gitx.exe' to get started."
}

Main