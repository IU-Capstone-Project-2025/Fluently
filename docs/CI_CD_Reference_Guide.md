
### CI/CD Setup and Architecture

The project implements a complex CI/CD pipeline using **GitHub Actions** with environment-specific deployments. The main configuration is located at:
- **Main CI/CD Configuration**: [deploy.yml](https://github.com/FluentlyOrg/Fluently-fork/blob/main/.github/workflows/deploy.yml)
- **Docker Configuration**: [docker-compose.yml](https://github.com/FluentlyOrg/Fluently-fork/blob/main/docker-compose.yml)

### Pipeline Structure

The CI/CD pipeline consists of two main jobs:

1. **Setup Job**: Environment determination and configuration
   - Automatically determines deployment environment based on branch:
     - `main` branch → Production environment (`fluently-app.ru`)
     - `develop/feature/fix` branches → Staging environment (`fluently-app.online`)
   - Sets environment-specific variables (domains, SSH credentials, ZeroTier IPs)
   - Supports manual environment override and rollback functionality

2. **Deploy Job**: Application deployment and health verification
   - Code synchronization with GitHub repository
   - Environment-specific configuration generation
   - Docker container orchestration
   - Comprehensive health checks
   - Automatic rollback on failure

### Tools and Technologies Used

- **CI/CD Platform**: GitHub Actions
- **Containerization**: Docker & Docker Compose
- **Deployment**: SSH-based deployment using `appleboy/ssh-action`
- **Configuration Management**: Environment variable substitution (`envsubst`)
- **Reverse Proxy**: Nginx with SSL termination
- **SSL/TLS**: Cloudflare Origin Certificates for end-to-end encryption
- **Networking**: ZeroTier for secure server access
- **Monitoring**: Built-in health checks for all services

### Advanced Features

- **Template-based Configuration**: Dynamic generation of environment-specific files
  - Nginx configuration templates
  - Backup script templates
  - Environment-specific variable substitution
- **Automated Backup System**: Pre-deployment backups with automatic cleanup
- **Rollback Capability**: Automatic and manual rollback functionality
- **Health Monitoring**: Multi-service health checks with timeout handling
- **Security**: ZeroTier VPN integration for secure server communication

### Challenges Faced

1. **Complex Service Dependencies**: Managing startup order and health checks for ML API, PostgreSQL, and backend services
2. **Configuration Management**: Avoiding file duplication through template-based approach
3. **Long Build Times**: ML API container builds required careful timeout handling
4. **Network Conflicts**: Docker network cleanup to prevent deployment conflicts
5. **SSL Certificate Management**: Cloudflare Origin Certificate integration for multiple domains

---

### Staging Environment

**Domain**: `fluently-app.online`  
**Server**: Ubuntu 24.04 LTS VDS staging server with dedicated staging user

#### Environment Setup
- **Branch Strategy**: All non-main branches (`develop`, `feature/*`, `fix/*`) deploy to staging
- **Directory Structure**: `/home/deploy-staging/Fluently-fork`
- **Backup Location**: `/home/deploy-staging/backups`
- **Network Access**: ZeroTier VPN for secure remote access

#### Deployment Process
1. **Pre-deployment**: Current state backup creation
2. **Code Update**: Git fetch, checkout, and pull latest changes
3. **Configuration**: Environment-specific `.env` file generation
4. **Template Processing**: Dynamic configuration file generation using `envsubst`
5. **Container Orchestration**: Docker Compose build and deployment
6. **Health Verification**: Multi-service health checks
7. **Cleanup**: Old Docker image removal and backup management

#### Services Architecture
- **Backend API**: Go-based REST API with Swagger documentation
- **ML API**: Python-based machine learning service with HuggingFace models
- **PostgreSQL**: Database with health checks and optimized configuration
- **Nginx**: Reverse proxy with SSL termination
- **Directus CMS**: Content management system
- **Monitoring Stack**: Prometheus, Grafana, Loki for observability

### Production Environment

**Domain**: `fluently-app.ru`  
**Server**: Ubuntu 24.04 LTS Production VDS with dedicated production user

#### Environment Setup
- **Branch Strategy**: Only `main` branch deploys to production
- **Directory Structure**: `/home/deploy/Fluently-fork`
- **Backup Location**: `/home/deploy/backups`
- **SSL/TLS**: Cloudflare Origin Certificates with Full (Strict) SSL mode
- **Network Access**: ZeroTier VPN for secure remote access

#### Production-Specific Features
1. **Enhanced Backup Strategy**: 
   - Automatic pre-deployment backups
   - Retention policy (last 5 backups)
   - Rollback capability with state restoration

2. **SSL/TLS Security**:
   - Cloudflare Origin Certificates for end-to-end encryption
   - HSTS headers and security configurations
   - Certificate management for `fluently-app.ru`

3. **Performance Optimizations**:
   - PostgreSQL tuned for production workloads
   - Docker BuildKit for faster builds
   - Resource limits and reservations for ML services

4. **Monitoring and Observability**:
   - Prometheus metrics collection
   - Grafana dashboards
   - Loki log aggregation
   - cAdvisor for container metrics
   - Node exporter for system metrics

#### VDS Setup and Infrastructure
- **Operating System**: Ubuntu 24.04 LTS
- **Network**: ZeroTier mesh network for secure access
- **SSL Certificates**: Cloudflare Origin Certificates managed automatically
- **Docker**: Latest Docker Engine with Compose V2
- **Storage**: External Docker volumes for data persistence
- **Backup Strategy**: Automated backup scripts with cron scheduling

#### Security Measures
- **Network Isolation**: ZeroTier VPN for all administrative access
- **SSL/TLS**: Full (Strict) SSL mode with Cloudflare
- **Container Security**: Non-root containers where possible
- **Secret Management**: Environment-specific secrets via GitHub Actions
- **Access Control**: Dedicated deployment users with minimal privileges
