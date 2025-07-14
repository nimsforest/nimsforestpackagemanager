# get.nimsforest.com

This directory contains the GitHub Pages site for **get.nimsforest.com** - the official installation portal for NimsForest Package Manager.

## URL Structure

| URL | File | Description |
|-----|------|-------------|
| `get.nimsforest.com` | `index.html` | Landing page with installation instructions |
| `get.nimsforest.com/linux` | `linux` | Linux installation script |
| `get.nimsforest.com/macos` | `macos` | macOS installation script |
| `get.nimsforest.com/windows` | `windows` | Windows PowerShell script |
| `get.nimsforest.com/install.sh` | `install.sh` | Universal installer (fallback) |

## Setup Instructions

### 1. Enable GitHub Pages
1. Go to repository Settings → Pages
2. Source: **Deploy from branch**
3. Branch: **main** 
4. Folder: **/docs**
5. Save

### 2. Configure Custom Domain
1. In repository Settings → Pages → Custom domain
2. Enter: `get.nimsforest.com`
3. Enable "Enforce HTTPS"

### 3. DNS Configuration
Point your domain to GitHub Pages:

```
CNAME get.nimsforest.com nimsforest.github.io
```

Or if using apex domain:
```
A get.nimsforest.com 185.199.108.153
A get.nimsforest.com 185.199.109.153  
A get.nimsforest.com 185.199.110.153
A get.nimsforest.com 185.199.111.153
```

## Files

- `index.html` - Beautiful landing page with platform detection
- `linux` - Linux installation script (no extension)
- `macos` - macOS installation script (no extension)  
- `windows` - Windows PowerShell script (no extension)
- `install.sh` - Universal fallback installer
- `_config.yml` - Jekyll configuration
- `CNAME` - Custom domain configuration
- `.htaccess` - HTTP headers and caching

## Testing

Test locally with Jekyll:
```bash
cd docs
bundle exec jekyll serve
```

Test install URLs:
```bash
curl -fsSL localhost:4000/linux
curl -fsSL localhost:4000/macos
curl -fsSL localhost:4000/windows
```

## Features

✅ **Cross-platform** - Works on Linux, macOS, Windows  
✅ **Smart detection** - Auto-detects user's platform  
✅ **Copy buttons** - Easy copy-paste installation  
✅ **Mobile friendly** - Responsive design  
✅ **Fast loading** - Optimized for performance  
✅ **Secure** - HTTPS with security headers  

## Deployment

Changes are automatically deployed when pushed to the `main` branch. GitHub Pages will rebuild the site within a few minutes.

Monitor deployment at: `https://github.com/nimsforest/nimsforestpackagemanager/deployments`