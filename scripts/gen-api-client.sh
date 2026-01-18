#!/bin/bash
set -e

# get the directory of the currently currently script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$DIR/.."

echo "🔄 Step 1: Generating Swagger documentation in Backend..."
cd "$PROJECT_ROOT/backend"
make swagger

echo "🔄 Step 2: Updating Shared API specs..."
cd "$PROJECT_ROOT"
# Ensure shared/api exists
mkdir -p shared/api
cp backend/docs/swagger.json shared/api/swagger.json
cp backend/docs/swagger.yaml shared/api/swagger.yaml
echo "   ✅ Copied swagger.json/yaml to shared/api/"

echo "🔄 Step 3: Generating TypeScript client into Shared..."
# We use the openapi-generator-cli installed in frontend's node_modules
cd "$PROJECT_ROOT/frontend"
# Generate into ../shared/types
pnpm exec openapi-generator-cli generate \
    -i "$PROJECT_ROOT/shared/api/swagger.json" \
    -g typescript-axios \
    -o "$PROJECT_ROOT/shared/types" \
    --additional-properties=supportsES6=true,npmName=@piemdm/api-client,withSeparateModelsAndApi=true,apiPackage=api,modelPackage=models,useSingleRequestParameter=false \
    --skip-validate-spec

echo "   ✅ Generated types in shared/types/"

echo "✅ API utilization workflow complete!"
echo "   - Spec: shared/api/swagger.json"
echo "   - Types: shared/types/ (@piemdm/api-client)"
echo "   - Frontend uses workspace reference: @piemdm/api-client"
