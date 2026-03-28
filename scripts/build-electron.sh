#!/bin/bash
set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "=== 构建 FMPS Electron 桌面版 ==="

# Build Go binary
echo "1. 编译后端..."
cd "$ROOT/server"
go build -o "$ROOT/electron/resources/fmps-server" .
echo "   后端编译完成"

# Build web
echo "2. 构建前端..."
cd "$ROOT/web"
npm install
npm run build
echo "   前端构建完成"

# Build electron
echo "3. 构建 Electron..."
cd "$ROOT/electron"
npm install
npm run build
echo "   Electron 构建完成"

echo ""
echo "=== 构建完成！==="
echo "安装包位于: electron/dist/"
