# Ubuntu ECS 部署指南

## 准备工作

1. 编译后端二进制:
   ```bash
   cd server && GOOS=linux GOARCH=amd64 go build -o fmps-server .
   ```

2. 构建前端:
   ```bash
   cd web && npm install && npm run build
   ```

## 部署步骤

1. 创建目录:
   ```bash
   sudo mkdir -p /opt/fmps/data
   sudo mkdir -p /opt/fmps/web
   ```

2. 复制文件:
   ```bash
   sudo cp server/fmps-server /opt/fmps/
   sudo cp -r web/dist/* /opt/fmps/web/
   sudo chmod +x /opt/fmps/fmps-server
   ```

3. 安装 systemd 服务:
   ```bash
   sudo cp deploy/fmps.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable fmps
   sudo systemctl start fmps
   ```

4. 检查状态:
   ```bash
   sudo systemctl status fmps
   ```

## Nginx 反向代理（可选）

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location /api/ {
        proxy_pass http://localhost:8080;
    }
    
    location / {
        root /opt/fmps/web;
        try_files $uri $uri/ /index.html;
    }
}
```
