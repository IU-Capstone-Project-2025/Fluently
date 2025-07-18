# Platform Support Documentation

## Overview

The Fluently project is designed to run on multiple platforms but has specific requirements and limitations that users should be aware of.

## Supported Platforms

### ✅ Fully Supported
- **Ubuntu 20.04+** (Linux)
- **macOS 12+** (Intel and Apple Silicon)
- **Windows 10/11** (with Docker Desktop)

### 🐳 Docker Container Architecture

**Important**: All Fluently Docker images are built for **Linux containers only**. This means:

- **Linux**: Native support ✅
- **macOS**: Runs Linux containers via Docker Desktop ✅
- **Windows**: Requires Docker Desktop in **Linux container mode** ✅

## Platform-Specific Requirements

### Linux (Ubuntu/Debian)
- Docker and Docker Compose
- Make utility
- Git
- No additional configuration needed

### macOS
- Docker Desktop for Mac
- Make utility (included with Xcode Command Line Tools)
- Git (included with Xcode Command Line Tools)
- Rosetta 2 (for Intel compatibility on Apple Silicon)

### Windows
- **Docker Desktop for Windows** (with Linux container support)
- **WSL2** (Windows Subsystem for Linux 2)
- **Make** (can be installed via Chocolatey: `choco install make`)
- **Git for Windows**

#### Windows-Specific Configuration

1. **Docker Desktop Must Use Linux Containers**:
   - Right-click Docker Desktop system tray icon
   - Select "Switch to Linux containers" if currently using Windows containers
   - This is **required** - Windows containers are not supported

2. **WSL2 Backend**:
   - Enable WSL2 in Docker Desktop settings
   - This provides better performance and compatibility

3. **Resource Allocation**:
   - Allocate at least 8GB RAM to Docker Desktop
   - Ensure at least 50GB of disk space available

## Docker Image Architecture Support

All Fluently images are built with multi-arch support for:
- `linux/amd64` (Intel/AMD x86_64)
- `linux/arm64` (ARM64/Apple Silicon)

**Not supported**:
- `windows/amd64` (Windows containers)
- `windows/arm64` (Windows ARM containers)

## Why Linux Containers Only?

The Fluently project uses several components that are Linux-specific:

1. **Python ML Libraries**: Many Python packages (especially ML/AI libraries) have Linux-optimized builds
2. **PostgreSQL**: Database optimized for Linux environments
3. **Nginx**: Web server with Linux-specific optimizations
4. **Base Images**: All service images are based on Linux distributions (Ubuntu, Alpine, etc.)

## Platform Testing

Our CI/CD pipeline tests the following:

### Automated Testing
- **Ubuntu**: Full integration testing
- **macOS**: Docker Desktop startup and basic functionality
- **Windows**: Docker Desktop startup, Linux container verification, and basic functionality

### Manual Testing Required
- Production deployments (Linux servers)
- Complex multi-service interactions
- Performance testing under load

## Troubleshooting

### Windows Issues

**Problem**: "image operating system \"linux\" cannot be used on this platform"
**Solution**: Switch Docker Desktop to Linux container mode

**Problem**: WSL2 integration not working
**Solution**: 
1. Enable WSL2 in Windows Features
2. Update WSL2 kernel: `wsl --update`
3. Restart Docker Desktop

**Problem**: Performance issues
**Solution**: 
1. Allocate more resources to Docker Desktop
2. Enable WSL2 integration
3. Move project to WSL2 filesystem for better performance

### macOS Issues

**Problem**: Docker Desktop won't start
**Solution**: 
1. Restart Docker Desktop
2. Check available disk space
3. Reset Docker Desktop to factory defaults if needed

**Problem**: Apple Silicon compatibility
**Solution**: 
1. Ensure Rosetta 2 is installed: `softwareupdate --install-rosetta`
2. Use ARM64 images when available (automatically selected)

### General Issues

**Problem**: Services fail to start
**Solution**: 
1. Check Docker daemon is running: `docker info`
2. Verify system resources (RAM, disk space)
3. Check for port conflicts: `make check-ports`

**Problem**: Build failures
**Solution**: 
1. Clear Docker cache: `docker builder prune`
2. Pull latest images: `docker compose pull`
3. Rebuild images: `docker compose build --no-cache`

## Development Environment Setup

### Recommended Resources
- **RAM**: 16GB+ (8GB minimum)
- **Disk Space**: 100GB+ available
- **CPU**: 4+ cores for smooth development

### Environment Variables
The following environment variables should be set based on your platform:

```bash
# Docker image tags (automatically set by CI/CD)
BACKEND_TAG=latest-develop
NGINX_TAG=latest-develop
TELEGRAM_TAG=latest-develop
ML_API_TAG=latest-develop
```

## Production Deployment

**Production deployments are Linux-only** and typically use:
- Ubuntu 20.04+ LTS servers
- Docker CE (not Docker Desktop)
- Optimized resource allocation
- Load balancing and scaling

## Support and Compatibility Matrix

| Platform | Local Development | CI/CD Testing | Production |
|----------|------------------|---------------|------------|
| Ubuntu Linux | ✅ Full | ✅ Full | ✅ Full |
| macOS | ✅ Full | ✅ Basic | ❌ No |
| Windows | ✅ Full* | ✅ Basic | ❌ No |
| Docker Desktop | ✅ Required | ✅ Tested | ❌ No |
| Docker CE | ✅ Supported | ✅ Tested | ✅ Required |

*Requires Linux container mode

## Getting Help

If you encounter platform-specific issues:

1. Check this documentation first
2. Review the troubleshooting section
3. Check GitHub Issues for known problems
4. Create a new issue with:
   - Platform details (OS, version, architecture)
   - Docker version and configuration
   - Error messages and logs
   - Steps to reproduce

## Future Improvements

Planned improvements for platform support:

- [ ] Better Windows integration with WSL2
- [ ] Automated platform detection and configuration
- [ ] Performance optimizations for different platforms
- [ ] Better error messages for platform-specific issues
- [ ] Docker Desktop configuration validation scripts
