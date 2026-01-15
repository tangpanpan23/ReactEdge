# ReactEdge 配置文件说明

## 概述

ReactEdge 支持通过 YAML 配置文件进行全面的系统配置，包括服务器设置、AI 参数、日志配置、监控设置等。配置文件位于 `config/app.yaml`，支持环境变量覆盖。

## 配置结构

### 服务器配置 (server)

```yaml
server:
  # Web服务端口 (6000-6999范围)
  port: "6000"

  # 服务器超时配置 (秒)
  read_timeout: 30
  write_timeout: 30

  # 服务器主机地址
  host: "0.0.0.0"

  # 是否启用HTTPS
  tls_enabled: false

  # TLS证书路径 (当启用HTTPS时需要)
  tls_cert_file: ""
  tls_key_file: ""
```

### AI服务配置 (ai)

```yaml
ai:
  # AI模式: internal(对内使用TAL) 或 external(对外使用开放模型)
  mode: "internal"

  # 最大分析时间 (秒)
  max_analysis_time: 60

  # 是否启用缓存
  cache_enabled: true

  # 缓存配置
  cache:
    # 缓存过期时间 (秒)
    ttl: 3600
    # 最大缓存条目数
    max_entries: 1000

  # 熔断器配置
  circuit_breaker:
    # 最大失败次数
    max_failures: 5
    # 熔断恢复时间 (秒)
    timeout: 60

  # 并发控制
  concurrency:
    # 最大并发请求数
    max_concurrent: 10

  # 请求限流
  rate_limit:
    # 每小时最大请求数
    requests_per_hour: 1000
    # 突发请求数
    burst_limit: 100
```

### 日志配置 (logging)

```yaml
logging:
  # 日志级别: debug, info, warn, error
  level: "info"

  # 日志格式: json 或 text
  format: "text"

  # 是否输出到文件
  file_enabled: false

  # 日志文件路径
  file_path: "logs/reactedge.log"

  # 日志轮转配置
  rotation:
    # 最大文件大小 (MB)
    max_size: 100
    # 最大备份文件数
    max_backups: 5
    # 最大保存天数
    max_age: 30
```

### 监控配置 (monitoring)

```yaml
monitoring:
  # 是否启用监控
  enabled: true

  # 监控端口 (与主服务端口不同)
  port: "6060"

  # 健康检查路径
  health_check_path: "/health"

  # 指标收集
  metrics:
    # 是否启用Prometheus指标
    prometheus_enabled: true
    # 指标路径
    path: "/metrics"
```

### 开发环境配置 (development)

```yaml
development:
  # 是否启用调试模式
  debug: false

  # 是否启用CORS
  cors_enabled: true

  # CORS允许的源
  cors_origins: ["http://localhost:3000", "http://localhost:8080", "http://localhost:6000"]

  # 是否启用请求日志
  request_logging: true

  # 是否启用SQL查询日志
  sql_logging: false
```

### 生产环境配置 (production)

```yaml
production:
  # 是否启用Gzip压缩
  gzip_enabled: true

  # 是否启用请求限流
  rate_limiting_enabled: true

  # 静态文件缓存时间 (秒)
  static_cache_ttl: 86400

  # 是否启用安全头
  security_headers_enabled: true
```

## 环境变量覆盖

所有配置项都可以通过环境变量进行覆盖，环境变量优先级高于配置文件。

### 服务器配置环境变量

```bash
# 服务器端口
SERVER_PORT=6000

# 服务器主机
SERVER_HOST=0.0.0.0

# 超时配置
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30
```

### AI配置环境变量

```bash
# AI模式
AI_MODE=internal

# 最大分析时间
AI_MAX_ANALYSIS_TIME=60

# 缓存开关
AI_CACHE_ENABLED=true
```

### 日志配置环境变量

```bash
# 日志级别
LOG_LEVEL=info

# 日志格式
LOG_FORMAT=text
```

### 开发环境配置环境变量

```bash
# 调试模式
DEBUG=false
```

## 配置验证

系统会对关键配置进行验证：

- **端口范围**: 必须在 6000-6999 范围内
- **AI模式**: 只能是 "internal" 或 "external"
- **超时设置**: 不能小于1秒

## 使用示例

### 1. 修改端口配置

```yaml
# config/app.yaml
server:
  port: "6001"  # 修改为6001端口
```

### 2. 切换AI模式

```yaml
# config/app.yaml
ai:
  mode: "external"  # 切换到外部AI模式
```

### 3. 启用HTTPS

```yaml
# config/app.yaml
server:
  tls_enabled: true
  tls_cert_file: "/path/to/cert.pem"
  tls_key_file: "/path/to/key.pem"
```

### 4. 环境变量覆盖

```bash
# 使用环境变量覆盖端口
SERVER_PORT=6002 go run .

# 切换AI模式
AI_MODE=external go run .
```

## 默认配置

如果 `config/app.yaml` 文件不存在或无法加载，系统将使用内置的默认配置：

- 端口: 6000
- AI模式: internal
- 超时: 30秒
- 缓存: 启用
- 日志级别: info

## 热重载

当前版本暂不支持配置热重载，修改配置后需要重启服务。未来版本将支持配置文件的热重载功能。

## 注意事项

1. **端口范围**: Web服务端口必须在6000-6999范围内，这是项目的端口规范
2. **权限要求**: 确保应用有权限绑定指定的端口和读取证书文件
3. **环境变量**: 环境变量会覆盖配置文件中的相同设置
4. **配置验证**: 系统会在启动时验证配置的有效性，无效配置会导致启动失败
