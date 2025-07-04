#!/bin/bash

# Pipeline Validation Script
# This script validates that the new CI/CD pipeline is set up correctly

echo "🔍 Validating CI/CD Pipeline Setup..."
echo "======================================"

# Check if required files exist
echo "📁 Checking required files..."

required_files=(
    ".github/workflows/deploy.yml"
    "docker-compose.production.yml"
    "docs/Install_Local.md"
    "DEPLOYMENT_GUIDE.md"
    "sonar-project.properties"
)

missing_files=()
for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file"
    else
        echo "❌ $file"
        missing_files+=("$file")
    fi
done

if [ ${#missing_files[@]} -gt 0 ]; then
    echo ""
    echo "❌ Missing required files:"
    printf " - %s\n" "${missing_files[@]}"
    exit 1
fi

echo ""
echo "🐳 Checking Docker configuration..."

# Check if docker-compose.production.yml uses the correct images
if grep -q "docker.io/fluentlyorg/" docker-compose.production.yml; then
    echo "✅ Production compose uses Docker Hub images"
else
    echo "❌ Production compose doesn't use Docker Hub images"
    exit 1
fi

# Check if workflow has the correct structure
echo ""
echo "⚙️ Checking workflow structure..."

if grep -q "build:" .github/workflows/deploy.yml; then
    echo "✅ Build stage found"
else
    echo "❌ Build stage missing"
    exit 1
fi

if grep -q "test:" .github/workflows/deploy.yml; then
    echo "✅ Test stage found"
else
    echo "❌ Test stage missing"
    exit 1
fi

if grep -q "publish:" .github/workflows/deploy.yml; then
    echo "✅ Publish stage found"
else
    echo "❌ Publish stage missing"
    exit 1
fi

if grep -q "deploy:" .github/workflows/deploy.yml; then
    echo "✅ Deploy stage found"
else
    echo "❌ Deploy stage missing"
    exit 1
fi

echo ""
echo "📋 Checking environment configuration..."

if [ -f ".env.example" ]; then
    echo "✅ .env.example exists"
else
    echo "❌ .env.example missing"
    exit 1
fi

echo ""
echo "🎉 Pipeline validation completed successfully!"
echo ""
echo "📝 Next steps:"
echo "1. Set up GitHub secrets (see DEPLOYMENT_GUIDE.md)"
echo "2. Disable SonarCloud Automatic Analysis (see previous instructions)"
echo "3. Push to develop branch to trigger first build"
echo "4. Teaching Assistants can use the updated Install_Local.md guide"
echo ""
echo "🚀 Your new CI/CD pipeline is ready!"
