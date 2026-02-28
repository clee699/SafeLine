# SafeLine WAF Docker部署指南

## 目录结构

```
SafeLine/
├── management/
│   └── webserver/
│       ├── Dockerfile          # 管理服务Dockerfile
│       ├── config.yml          # 管理服务配置文件
│       └── static/             # 静态资源目录
├── docker-compose.yml          # Docker Compose配置文件
├── .env.example                # 环境变量示例
└── README_DOCKER.md            # 本部署指南
```

## 部署准备

1. 安装Docker和Docker Compose
   - 参考官方文档：[Docker安装](https://docs.docker.com/get-docker/)、[Docker Compose安装](https://docs.docker.com/compose/install/)

2. 克隆仓库
   ```bash
   git clone <your-repository-url>
   cd SafeLine
   ```

3. 配置环境变量（可选）
   - 复制示例环境变量文件
   ```bash
   cp .env.example .env
   ```
   - 根据需要修改`.env`文件中的配置项

## 部署步骤

### 1. 启动所有服务

```bash
docker-compose up -d
```

该命令会启动以下服务：
- **postgres**: PostgreSQL数据库服务（端口5432）
- **webserver**: 管理服务（端口8000）
- **detector**: 检测服务（端口8080和8443）

### 2. 访问管理界面

部署完成后，通过以下地址访问管理界面：
```
http://<your-server-ip>:8000/alert
```

### 3. 初始化数据库（首次部署）

如果是首次部署，需要初始化数据库：

```bash
docker exec -it safeline-webserver /app/safeline-webserver -reset_user admin
```

## 服务说明

### postgres

PostgreSQL数据库服务，用于存储SafeLine WAF的配置和日志数据。

**配置项**：
- POSTGRES_USER: 数据库用户名（默认：safeline）
- POSTGRES_PASSWORD: 数据库密码（默认：safeline）
- POSTGRES_DB: 数据库名称（默认：safeline）

### webserver

SafeLine WAF的管理服务，提供RESTful API和Web管理界面。

**配置项**：
- DB_HOST: 数据库主机（默认：postgres）
- DB_PORT: 数据库端口（默认：5432）
- DB_USER: 数据库用户名（默认：safeline）
- DB_PASSWORD: 数据库密码（默认：safeline）
- DB_NAME: 数据库名称（默认：safeline）

**访问地址**：
- 管理界面：http://<your-server-ip>:8000/alert
- API地址：http://<your-server-ip>:8000/api

### detector

SafeLine WAF的检测服务，负责流量检测和拦截。

**访问地址**：
- HTTP: http://<your-server-ip>:8080
- HTTPS: https://<your-server-ip>:8443

## 配置文件说明

### management/webserver/config.yml

管理服务的配置文件，包含数据库连接、服务器配置、告警配置等。

**告警相关配置项**：
```yaml
alert:
  enabled: true              # 是否启用告警
  check_interval: 10         # 检测间隔（秒）
  default_level: high        # 默认告警级别
  max_alert_count: 100       # 最大告警数量
```

## 常用命令

### 查看服务状态

```bash
docker-compose ps
```

### 查看服务日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看指定服务日志
docker-compose logs webserver

# 实时查看服务日志
docker-compose logs -f webserver
```

### 停止所有服务

```bash
docker-compose down
```

### 重启所有服务

```bash
docker-compose restart
```

### 重建服务

```bash
# 重建所有服务
docker-compose up -d --build

# 重建指定服务
docker-compose up -d --build webserver
```

## 更新服务

1. 更新代码
   ```bash
git pull
```

2. 重建并重启服务
   ```bash
docker-compose up -d --build
```

## 数据持久化

- **数据库数据**：存储在`safeline-postgres-data`卷中
- **检测服务数据**：存储在`safeline-detector-data`卷中
- **配置文件**：通过挂载主机文件到容器中实现持久化

## 端口映射

| 服务       | 容器端口 | 主机端口 | 说明               |
|------------|----------|----------|--------------------|
| postgres   | 5432     | 5432     | 数据库服务端口     |
| webserver  | 8000     | 8000     | 管理服务端口       |
| detector   | 8080     | 8080     | 检测服务HTTP端口   |
| detector   | 8443     | 8443     | 检测服务HTTPS端口  |

## 安全建议

1. 修改默认数据库密码
2. 配置防火墙，限制访问IP
3. 启用HTTPS访问管理界面
4. 定期备份数据库数据
5. 定期更新服务镜像

## 故障排查

### 服务无法启动

1. 查看服务日志
   ```bash
docker-compose logs -f <service-name>
```

2. 检查端口是否被占用
   ```bash
netstat -tuln | grep <port>
```

### 数据库连接失败

1. 检查数据库服务是否正常运行
   ```bash
docker-compose ps postgres
```

2. 检查数据库连接配置
   ```bash
cat management/webserver/config.yml
```

### 管理界面无法访问

1. 检查webserver服务是否正常运行
   ```bash
docker-compose ps webserver
```

2. 检查防火墙设置
   ```bash
iptables -L
```

## 联系支持

如果遇到问题，请联系技术支持或查看官方文档。