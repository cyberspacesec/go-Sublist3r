# go-Sublist3r

一个用Go语言实现的[Sublist3r](https://github.com/aboul3la/Sublist3r)子域名枚举工具包装器。该项目提供了一种使用Docker运行Sublist3r的便捷方法，无需在系统上直接安装Python依赖项。

## 特性

- 自动使用Docker Hub上的公共镜像或构建Sublist3r的Docker镜像
- 与Sublist3r选项一致的命令行界面
- 首次运行时自动下载或创建Docker镜像
- 处理容器和主机之间的输出文件
- 当Docker不可用时提供模拟功能
- 提供RESTful API，支持远程调用和集成
- 支持API密钥认证，保护API访问安全

## 项目结构

```
.
├── Dockerfile         # 构建Sublist3r Docker镜像的定义文件
├── Makefile           # 构建和安装脚本
├── README.md          # 项目文档
├── go.mod             # Go模块定义
├── main.go            # 主程序入口
├── scripts            # 辅助脚本目录
│   └── sublist3r-simulator.sh # Docker不可用时的模拟脚本
└── pkg                # 项目包
    ├── api            # API服务器实现
    │   ├── docs.go    # API文档
    │   └── server.go  # API服务器实现
    ├── cmd            # 命令行界面定义
    │   ├── api.go     # API服务器命令
    │   ├── build.go   # 构建Docker镜像子命令
    │   ├── pull.go    # 拉取Docker镜像子命令
    │   ├── register.go # 命令注册
    │   ├── root.go    # 根命令
    │   └── scan.go    # 扫描子命令
    ├── docker         # Docker相关功能
    │   └── docker.go  # Docker操作实现
    └── ui             # 用户界面相关功能
        └── ui.go      # UI实现
```

## 安装

### 前提条件

- Go 1.16或更高版本
- Docker (可选，如果不可用会使用模拟模式)

### 从源代码构建

```bash
git clone https://github.com/cyberspacesec/go-Sublist3r.git
cd go-Sublist3r
make build
```

要全局安装:

```bash
make install
```

## 使用方法

### 获取Docker镜像

程序会自动获取镜像，但您也可以使用以下命令显式操作:

#### 从Docker Hub拉取镜像 (推荐)
```bash
go-sublist3r pull-docker-image
```

#### 构建本地Docker镜像
```bash
go-sublist3r build-docker-image
```

默认情况下，程序会优先尝试从Docker Hub拉取镜像，仅当拉取失败时才会尝试构建本地镜像。

### 运行子域名枚举

基本用法:

```bash
go-sublist3r scan -d example.com
```

使用其他选项:

```bash
go-sublist3r scan -d example.com -b -p 80,443 -v -o results.txt
```

### 启动API服务器

启动API服务器，允许通过HTTP请求进行子域名枚举:

```bash
go-sublist3r api --port 8080 --workers 5 --capacity 100
```

启动带有API密钥认证的API服务器:

```bash
go-sublist3r api --port 8080 --api-key "your-secret-key"
```

参数说明:
- `--port`: API服务器监听的端口号 (默认: 8080)
- `--workers`: 处理扫描请求的工作线程数 (默认: 5)
- `--capacity`: 最大并发扫描请求数 (默认: 100)
- `--api-key`: API密钥，用于请求认证（留空表示不启用认证）

服务器启动后，可以通过浏览器访问API文档: `http://localhost:8080/docs`

### 命令行选项

#### 扫描命令选项 (scan)

以下选项与原始Sublist3r相对应:

- `-d, --domain`: 要枚举子域名的域名（必需）
- `-b, --bruteforce`: 启用subbrute暴力破解模块
- `-p, --ports`: 扫描找到的子域名上的指定tcp端口
- `-v, --verbose`: 启用详细模式并实时显示结果
- `-t, --threads`: 用于subbrute暴力破解的线程数（默认：30）
- `-e, --engines`: 指定以逗号分隔的搜索引擎列表
- `-o, --output`: 将结果保存到文本文件
- `-n, --no-color`: 不带颜色输出

## API使用

go-Sublist3r提供了RESTful API，支持远程调用和与其他系统集成。

### API认证

如果API服务器启用了API密钥认证，所有API请求都需要包含有效的API密钥。有两种方式提供API密钥：

1. **HTTP请求头**: 在请求头中添加 `X-API-Key: your_api_key`
2. **URL查询参数**: 在URL中添加 `?api_key=your_api_key`

API文档页面 (`/docs`) 无需认证即可访问。

### API端点

#### 启动扫描

```
POST /api/v1/scan
```

请求体示例:
```json
{
  "domain": "example.com",
  "bruteforce": true,
  "ports": "80,443",
  "verbose": true,
  "threads": 40,
  "engines": "baidu,yahoo",
  "no_color": false,
  "callback_url": "https://your-callback-server.com/notify"
}
```

响应示例:
```json
{
  "id": "5f3a9c21-7b1e-4c49-8a0d-3e51735ce83a",
  "message": "Scan started"
}
```

#### 同步扫描(等待结果)

```
POST /api/v1/scan/sync
```

这是一个同步接口，请求会阻塞直到扫描完成或达到5分钟超时限制，结果会直接返回而不需要额外的轮询。

请求体示例:
```json
{
  "domain": "example.com",
  "bruteforce": false,
  "ports": "80,443",
  "verbose": true,
  "threads": 40
}
```

响应示例:
```json
{
  "id": "5f3a9c21-7b1e-4c49-8a0d-3e51735ce83a",
  "status": "completed",
  "domain": "example.com",
  "start_time": "2023-08-15T10:30:45Z",
  "end_time": "2023-08-15T10:32:12Z",
  "results": [
    "www.example.com",
    "mail.example.com",
    "blog.example.com",
    "api.example.com",
    "dev.example.com"
  ]
}
```

#### 获取扫描状态和结果

```
GET /api/v1/scan/{id}
```

响应示例:
```json
{
  "id": "5f3a9c21-7b1e-4c49-8a0d-3e51735ce83a",
  "status": "completed",
  "domain": "example.com",
  "start_time": "2023-08-15T10:30:45Z",
  "end_time": "2023-08-15T10:32:12Z",
  "results": [
    "www.example.com",
    "mail.example.com",
    "blog.example.com",
    "api.example.com",
    "dev.example.com"
  ]
}
```

#### 获取所有扫描

```
GET /api/v1/scans
```

### API使用示例

#### 使用curl

```bash
# 启动扫描 (无认证)
curl -X POST http://localhost:8080/api/v1/scan \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "bruteforce": false,
    "ports": "80,443",
    "verbose": true
  }'

# 使用同步接口 (无认证)
curl -X POST http://localhost:8080/api/v1/scan/sync \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "bruteforce": false,
    "ports": "80,443"
  }'

# 启动扫描 (有认证)
curl -X POST http://localhost:8080/api/v1/scan \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "domain": "example.com",
    "bruteforce": false,
    "ports": "80,443",
    "verbose": true
  }'

# 使用同步接口 (有认证)
curl -X POST http://localhost:8080/api/v1/scan/sync \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "domain": "example.com",
    "bruteforce": false
  }'

# 使用URL参数进行认证
curl -X GET "http://localhost:8080/api/v1/scan/5f3a9c21-7b1e-4c49-8a0d-3e51735ce83a?api_key=your_api_key"

# 获取所有扫描
curl -X GET http://localhost:8080/api/v1/scans -H "X-API-Key: your_api_key"
```

#### 使用Python

```python
import requests
import json

# API密钥
api_key = "your_api_key"

# 添加认证头
headers = {
    "Content-Type": "application/json",
    "X-API-Key": api_key
}

# 方法1: 使用异步接口
def async_scan():
    # 启动扫描
    scan_data = {
        "domain": "example.com",
        "bruteforce": True,
        "ports": "80,443",
        "verbose": True
    }

    response = requests.post(
        "http://localhost:8080/api/v1/scan", 
        headers=headers, 
        json=scan_data
    )
    scan_id = response.json()["id"]
    print(f"Scan started with ID: {scan_id}")

    # 检查扫描状态
    status_response = requests.get(
        f"http://localhost:8080/api/v1/scan/{scan_id}",
        headers=headers
    )
    status_data = status_response.json()
    print(f"Scan status: {status_data['status']}")

    # 当状态为 completed 时获取结果
    if status_data["status"] == "completed":
        for subdomain in status_data["results"]:
            print(subdomain)

# 方法2: 使用同步接口(更简单)
def sync_scan():
    # 使用同步接口启动扫描并等待结果
    scan_data = {
        "domain": "example.com",
        "bruteforce": False,
        "ports": "80,443",
    }

    response = requests.post(
        "http://localhost:8080/api/v1/scan/sync", 
        headers=headers, 
        json=scan_data,
        timeout=300  # 5分钟超时
    )
    
    if response.status_code == 200:
        result = response.json()
        print(f"Scan completed for: {result['domain']}")
        print(f"Found {len(result['results'])} subdomains:")
        for subdomain in result["results"]:
            print(subdomain)
    else:
        print(f"Error: {response.status_code} - {response.text}")

# 选择使用同步或异步方法
sync_scan()  # 更简单的方式
```

## Docker镜像信息

程序使用以下Docker镜像:

- 默认从Docker Hub获取: [trickest/sublist3r](https://hub.docker.com/r/trickest/sublist3r)
- 本地构建: 使用项目目录中的Dockerfile

## 官方网站

go-Sublist3r提供了一个官方网站，您可以在GitHub Pages上访问：https://cyberspacesec.github.io/go-Sublist3r

### 自动部署网站

项目配置了GitHub Actions自动化工作流，会在提交到master/main分支或合并PR时自动部署website目录下的内容到GitHub Pages。

工作流程配置文件位于`.github/workflows/deploy-website.yml`，其工作原理如下：

1. 在推送到master/main分支且修改了website目录下文件时触发
2. 在PR被合并到master/main分支且修改了website目录下文件时触发
3. 自动将website目录部署到gh-pages分支
4. GitHub Pages配置为使用gh-pages分支作为源

### 手动部署网站

如果需要手动部署网站，可以参考`website/README.md`中的详细说明。基本步骤包括：

1. 将website目录中的内容复制到gh-pages分支
2. 推送gh-pages分支到GitHub
3. 在仓库设置中启用GitHub Pages，并选择gh-pages分支作为源

## 贡献指南

欢迎通过以下方式为项目贡献：

1. 提交Bug报告或功能请求
2. 提交Pull Request改进代码
3. 改进文档和示例
4. 分享使用经验和案例研究

请确保代码遵循Go的编码规范，并包含适当的测试和文档。

## 致谢

本项目是Ahmed Aboul-Ela的[Sublist3r](https://github.com/aboul3la/Sublist3r)的封装器。 