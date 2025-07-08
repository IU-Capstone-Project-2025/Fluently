# SonarCloud Setup Guide

This guide shows how to set up SonarCloud for code quality analysis instead of running a heavy local SonarQube server.

## 1. Create SonarCloud Account

1. Go to [sonarcloud.io](https://sonarcloud.io)
2. Sign up with your GitHub account
3. Import your Fluently repository
4. Create a new project with key: `fluently-app`

## 2. Generate Token

1. Go to your SonarCloud account settings
2. Navigate to Security → Generate Tokens
3. Create a token named "Fluently-Local-Analysis"
4. Copy the token (you'll need this)

## 3. Local Setup

```bash
# Install SonarScanner CLI
make install-sonar-scanner

# Set your SonarCloud token (replace with your actual token)
export SONAR_TOKEN=your_token_here

# Add to your .bashrc or .zshrc for persistence
echo 'export SONAR_TOKEN=your_token_here' >> ~/.bashrc
```

## 4. Running Analysis

```bash
# Quick quality check (generates coverage reports)
make code-quality

# Full analysis with SonarCloud upload
make quality-scan

# Or run individually
make sonar-scan
```

## 5. Configuration

The analysis is configured in `sonar-project.properties`:

- **Project Key**: `fluently-app`
- **Sources**: Backend Go, Python ML API, Frontend, Telegram Bot
- **Coverage Reports**: Go, Python coverage automatically included
- **Exclusions**: Test files, build artifacts, dependencies

## 6. Viewing Results

After running `make quality-scan`, view results at:
https://sonarcloud.io/dashboard?id=fluently-app

## 7. CI/CD Integration

For GitHub Actions, add the SONAR_TOKEN as a repository secret:

```yaml
- name: SonarCloud Scan
  uses: SonarSource/sonarcloud-github-action@master
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

## GitHub Actions Integration

For automatic code quality checks in CI/CD, add the SONAR_TOKEN as a repository secret:

1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `SONAR_TOKEN`
5. Value: Your SonarCloud token from step 2 above

The workflow will now run SonarCloud analysis on every push and block deployment if quality gate fails.

## Quality Gate Configuration

The deployment workflow (`deploy.yml`) includes:
- **Quality Check Job**: Runs before deployment
- **Coverage Generation**: Go and Python coverage reports
- **SonarCloud Analysis**: Code quality and security analysis
- **Deployment Gate**: Deployment only proceeds if quality checks pass

### Workflow Structure:
```
1. Setup (determine environment)
2. Quality Check (SonarCloud analysis) ← NEW
3. Deploy (only if quality check passes)
```

## Benefits vs Local SonarQube

✅ **Much faster startup** (no heavy container)  
✅ **No memory usage** (runs in cloud)  
✅ **Always up-to-date** (managed by SonarSource)  
✅ **Better for TAs** (no complex setup)  
✅ **Free for open source** projects  
✅ **Integrated with GitHub** (PR analysis)

## Troubleshooting

**Token Issues:**
```bash
# Verify token is set
echo $SONAR_TOKEN

# Test connection
sonar-scanner -Dsonar.token=$SONAR_TOKEN -Dsonar.projectKey=fluently-app -X
```

**Coverage Issues:**
```bash
# Ensure coverage files exist
ls -la backend/coverage.out
ls -la analysis/distractor_api/coverage.xml
```
