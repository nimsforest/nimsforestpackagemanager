# NimsForest Package Manager Installation Scripts

This directory contains installation scripts for **get.nimsforest.com** that enable easy cross-platform installation of nimsforestpm.

## Installation Methods

### Quick Install (Universal)
```bash
curl -fsSL get.nimsforest.com | sh
```

### Platform-Specific Install

**Linux:**
```bash
curl -fsSL get.nimsforest.com/linux | sh
```

**macOS:**
```bash
curl -fsSL get.nimsforest.com/macos | sh
```

**Windows (PowerShell):**
```powershell
irm get.nimsforest.com/windows | iex
```

**Developers (All Platforms):**
```bash
go install github.com/nimsforest/nimsforestpackagemanager/cmd@latest
```

## Files

- `install_universal.sh` - Detects OS and redirects to appropriate installer (served at `/`)
- `install_linux.sh` - Linux installation script (served at `/linux`)
- `install_macos.sh` - macOS installation script (served at `/macos`) 
- `install_windows.ps1` - Windows PowerShell script (served at `/windows`)

## URL Mapping for get.nimsforest.com

| URL | Script | Content-Type |
|-----|--------|--------------|
| `/` | `install_universal.sh` | `text/plain` |
| `/linux` | `install_linux.sh` | `text/plain` |
| `/macos` | `install_macos.sh` | `text/plain` |
| `/windows` | `install_windows.ps1` | `text/plain` |

## Features

### All Scripts Include:
- ✅ Architecture detection (amd64, arm64)
- ✅ Latest version fetching from GitHub releases
- ✅ Secure HTTPS downloads
- ✅ Proper error handling
- ✅ PATH management
- ✅ Installation verification
- ✅ Colored output for better UX

### Platform-Specific Features:

**Linux:**
- Installs to `/usr/local/bin` by default
- Handles sudo requirements automatically
- Works with all major distributions

**macOS:**
- Intel and Apple Silicon support
- Gatekeeper compatibility guidance
- Homebrew-style installation to `/usr/local/bin`

**Windows:**
- Installs to `%USERPROFILE%\.nimsforest\bin`
- Automatically adds to user PATH
- PowerShell Core and Windows PowerShell compatible

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NIMSFOREST_INSTALL_DIR` | Installation directory | `/usr/local/bin` (Unix), `%USERPROFILE%\.nimsforest\bin` (Windows) |

## Expected GitHub Release Structure

Scripts expect releases with these naming conventions:
```
nimsforestpm_linux_amd64
nimsforestpm_linux_arm64
nimsforestpm_darwin_amd64
nimsforestpm_darwin_arm64
nimsforestpm_windows_amd64.exe
nimsforestpm_windows_arm64.exe
```

## Testing Locally

Test scripts locally before deployment:

```bash
# Test Linux script
bash install/install_linux.sh

# Test macOS script  
bash install/install_macos.sh

# Test Windows script (in PowerShell)
.\install\install_windows.ps1

# Test universal detector
bash install/install_universal.sh
```

## Server Implementation

The get.nimsforest.com server should:

1. Serve scripts based on URL path
2. Set appropriate `Content-Type: text/plain` headers
3. Enable HTTPS/TLS
4. Add caching headers for better performance
5. Log downloads for analytics

Example nginx config:
```nginx
location = / {
    return 200;
    add_header Content-Type text/plain;
    alias /path/to/install_universal.sh;
}

location = /linux {
    return 200;
    add_header Content-Type text/plain;
    alias /path/to/install_linux.sh;
}

location = /macos {
    return 200;
    add_header Content-Type text/plain;
    alias /path/to/install_macos.sh;
}

location = /windows {
    return 200;
    add_header Content-Type text/plain;
    alias /path/to/install_windows.ps1;
}
```