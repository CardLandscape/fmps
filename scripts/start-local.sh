#!/bin/bash
set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "=== 启动 FMPS 本地服务 ==="

# Start Go API in background
echo "启动后端 API..."
cd "$ROOT/server"
go build -o "$ROOT/fmps-server" .
FMPS_PORT=8080 "$ROOT/fmps-server" &
SERVER_PID=$!
echo "后端 API 已启动 (PID: $SERVER_PID)"

# Wait for server
echo "等待服务器就绪..."
SERVER_READY=0
for i in $(seq 1 30); do
    if curl -s http://localhost:8080/api/login > /dev/null 2>&1; then
        echo "服务器已就绪"
        SERVER_READY=1
        break
    fi
    sleep 1
done

if [ "$SERVER_READY" -eq 0 ]; then
    echo "错误: 服务器未能在 30 秒内就绪，请检查日志"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

echo ""
echo "=== FMPS 已启动 ==="
echo "访问地址: http://localhost:8080"
echo "默认账号: admin / 123456"
echo ""
echo "按 Ctrl+C 停止服务"

# Wait
trap "kill $SERVER_PID 2>/dev/null; echo '服务已停止'" INT TERM
wait $SERVER_PID
