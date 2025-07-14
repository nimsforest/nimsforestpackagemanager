# NimsForest Package Manager Installation Script for Windows
# Usage: irm get.nimsforest.com/windows | iex

param(
    [string]$InstallDir = "$env:USERPROFILE\.nimsforest\bin",
    [string]$Version = "latest"
)

# Configuration
$BinaryName = "nimsforestpm.exe"
$GitHubRepo = "nimsforest/nimsforestpackagemanager"

# Helper functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    exit 1
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE.ToLower()
    switch ($arch) {
        "amd64" { return "amd64" }
        "x86_64" { return "amd64" }
        "arm64" { return "arm64" }
        default { 
            Write-Error "Unsupported architecture: $arch"
        }
    }
}

# Get latest release version from GitHub
function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$GitHubRepo/releases/latest" -ErrorAction Stop
        return $response.tag_name
    }
    catch {
        Write-Error "Failed to get latest version from GitHub: $($_.Exception.Message)"
    }
}

# Download and install binary
function Install-Binary {
    param(
        [string]$Version,
        [string]$Architecture
    )
    
    $downloadUrl = "https://github.com/$GitHubRepo/releases/download/$Version/${BinaryName.Replace('.exe', '')}_windows_$Architecture.exe"
    $tempFile = Join-Path $env:TEMP $BinaryName
    
    Write-Info "Downloading $BinaryName $Version for Windows/$Architecture..."
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -ErrorAction Stop
    }
    catch {
        Write-Error "Failed to download binary from $downloadUrl`: $($_.Exception.Message)"
    }
    
    # Create install directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        Write-Info "Creating install directory: $InstallDir"
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # Move binary to install directory
    $targetPath = Join-Path $InstallDir $BinaryName
    try {
        Move-Item -Path $tempFile -Destination $targetPath -Force
        Write-Success "Installed $BinaryName to $targetPath"
    }
    catch {
        Write-Error "Failed to install binary to $targetPath`: $($_.Exception.Message)"
    }
    
    return $targetPath
}

# Add to PATH
function Add-ToPath {
    param([string]$Directory)
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$Directory*") {
        Write-Info "Adding $Directory to user PATH..."
        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Success "Added to PATH. You may need to restart your terminal."
        
        # Also add to current session
        $env:PATH = "$env:PATH;$Directory"
    }
    else {
        Write-Info "$Directory is already in PATH"
    }
}

# Verify installation
function Test-Installation {
    param([string]$BinaryPath)
    
    Write-Info "Verifying installation..."
    
    try {
        & $BinaryPath hello
        Write-Success "Installation verified successfully!"
    }
    catch {
        Write-Warning "Could not run system check. Binary installed but may not be in PATH."
        Write-Info "Try running: $BinaryPath hello"
    }
}

# Main installation flow
function Main {
    Write-Info "Installing NimsForest Package Manager for Windows..."
    
    # Check prerequisites
    if (-not (Get-Command "Invoke-WebRequest" -ErrorAction SilentlyContinue)) {
        Write-Error "PowerShell with Invoke-WebRequest is required"
    }
    
    # Detect system
    $architecture = Get-Architecture
    $version = if ($Version -eq "latest") { Get-LatestVersion } else { $Version }
    
    Write-Info "Detected system: Windows/$architecture"
    Write-Info "Installing version: $version"
    
    # Install
    $binaryPath = Install-Binary -Version $version -Architecture $architecture
    
    # Add to PATH
    Add-ToPath -Directory $InstallDir
    
    # Verify
    Test-Installation -BinaryPath $binaryPath
    
    Write-Success "Installation complete!"
    Write-Host ""
    Write-Host "Get started with:" -ForegroundColor White
    Write-Host "  nimsforestpm hello" -ForegroundColor Cyan
    Write-Host "  nimsforestpm create-organization-workspace my-org" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "For help: nimsforestpm --help" -ForegroundColor White
    Write-Host ""
    Write-Host "Note: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
}

# Run main function
Main