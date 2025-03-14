package api

// API 文档和示例
const APIDocsHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>go-Sublist3r API 文档</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 1000px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        pre {
            background-color: #f5f5f5;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
        }
        code {
            font-family: Consolas, Monaco, 'Andale Mono', monospace;
        }
        .endpoint {
            background-color: #e8f4f8;
            padding: 15px;
            border-left: 5px solid #3498db;
            margin-bottom: 20px;
        }
        .method {
            font-weight: bold;
            color: #3498db;
        }
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
        .auth-note {
            background-color: #ffeeba;
            border-left: 5px solid #ffc107;
            padding: 15px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <h1>go-Sublist3r API 文档</h1>
    <p>go-Sublist3r 提供了 RESTful API，允许通过 HTTP 请求进行子域名枚举操作。</p>

    <div class="auth-note">
        <h3>API 认证</h3>
        <p>如果服务器启用了API密钥认证，所有API请求都需要包含有效的API密钥。有两种方式提供API密钥：</p>
        <ol>
            <li><strong>HTTP请求头</strong>: 在请求头中添加 <code>X-API-Key: your_api_key</code></li>
            <li><strong>URL查询参数</strong>: 在URL中添加 <code>?api_key=your_api_key</code></li>
        </ol>
        <p>API文档页面无需认证即可访问。</p>
    </div>

    <h2>API 端点</h2>

    <div class="endpoint">
        <h3><span class="method">POST</span> /api/v1/scan</h3>
        <p>启动一个新的子域名枚举扫描任务。</p>
        
        <h4>请求参数 (JSON):</h4>
        <table>
            <tr>
                <th>参数</th>
                <th>类型</th>
                <th>必需</th>
                <th>描述</th>
            </tr>
            <tr>
                <td>domain</td>
                <td>string</td>
                <td>是</td>
                <td>要枚举子域名的域名</td>
            </tr>
            <tr>
                <td>bruteforce</td>
                <td>boolean</td>
                <td>否</td>
                <td>是否启用暴力破解模块 (默认: false)</td>
            </tr>
            <tr>
                <td>ports</td>
                <td>string</td>
                <td>否</td>
                <td>要扫描的端口，如 "80,443"</td>
            </tr>
            <tr>
                <td>verbose</td>
                <td>boolean</td>
                <td>否</td>
                <td>是否启用详细输出 (默认: false)</td>
            </tr>
            <tr>
                <td>threads</td>
                <td>integer</td>
                <td>否</td>
                <td>线程数 (默认: 30)</td>
            </tr>
            <tr>
                <td>engines</td>
                <td>string</td>
                <td>否</td>
                <td>使用的搜索引擎，如 "baidu,yahoo"</td>
            </tr>
            <tr>
                <td>no_color</td>
                <td>boolean</td>
                <td>否</td>
                <td>是否禁用颜色输出 (默认: false)</td>
            </tr>
            <tr>
                <td>callback_url</td>
                <td>string</td>
                <td>否</td>
                <td>扫描完成后的回调 URL</td>
            </tr>
        </table>
        
        <h4>请求示例:</h4>
        <pre><code>{
  "domain": "example.com",
  "bruteforce": true,
  "ports": "80,443",
  "verbose": true,
  "threads": 40,
  "engines": "baidu,yahoo",
  "no_color": false,
  "callback_url": "https://your-callback-server.com/notify"
}</code></pre>

        <h4>响应:</h4>
        <pre><code>{
  "id": "5f3a9c21-7b1e-4c49-8a0d-3e51735ce83a",
  "message": "Scan started"
}</code></pre>
    </div>

    <div class="endpoint">
        <h3><span class="method">POST</span> /api/v1/scan/sync</h3>
        <p>启动一个新的子域名枚举扫描任务并等待扫描完成，返回完整结果。这是一个同步接口，请求会阻塞直到扫描完成或超时。</p>
        
        <h4>请求参数 (JSON):</h4>
        <table>
            <tr>
                <th>参数</th>
                <th>类型</th>
                <th>必需</th>
                <th>描述</th>
            </tr>
            <tr>
                <td>domain</td>
                <td>string</td>
                <td>是</td>
                <td>要枚举子域名的域名</td>
            </tr>
            <tr>
                <td>bruteforce</td>
                <td>boolean</td>
                <td>否</td>
                <td>是否启用暴力破解模块 (默认: false)</td>
            </tr>
            <tr>
                <td>ports</td>
                <td>string</td>
                <td>否</td>
                <td>要扫描的端口，如 "80,443"</td>
            </tr>
            <tr>
                <td>verbose</td>
                <td>boolean</td>
                <td>否</td>
                <td>是否启用详细输出 (默认: false)</td>
            </tr>
            <tr>
                <td>threads</td>
                <td>integer</td>
                <td>否</td>
                <td>线程数 (默认: 30)</td>
            </tr>
            <tr>
                <td>engines</td>
                <td>string</td>
                <td>否</td>
                <td>使用的搜索引擎，如 "baidu,yahoo"</td>
            </tr>
            <tr>
                <td>no_color</td>
                <td>boolean</td>
                <td>否</td>
                <td>是否禁用颜色输出 (默认: false)</td>
            </tr>
        </table>
        
        <p><strong>注意</strong>: 此接口有5分钟超时限制，不支持回调URL参数。</p>
        
        <h4>请求示例:</h4>
        <pre><code>{
  "domain": "example.com",
  "bruteforce": false,
  "ports": "80,443",
  "verbose": true,
  "threads": 40
}</code></pre>

        <h4>响应:</h4>
        <pre><code>{
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
}</code></pre>
    </div>

    <div class="endpoint">
        <h3><span class="method">GET</span> /api/v1/scan/{id}</h3>
        <p>获取指定扫描任务的状态和结果。</p>
        
        <h4>URL 参数:</h4>
        <table>
            <tr>
                <th>参数</th>
                <th>描述</th>
            </tr>
            <tr>
                <td>id</td>
                <td>扫描任务 ID</td>
            </tr>
        </table>
        
        <h4>响应:</h4>
        <pre><code>{
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
}</code></pre>

        <p>可能的状态值:</p>
        <ul>
            <li><code>pending</code>: 任务已创建但尚未开始执行</li>
            <li><code>running</code>: 任务正在执行中</li>
            <li><code>completed</code>: 任务已成功完成</li>
            <li><code>failed</code>: 任务执行失败</li>
        </ul>
    </div>

    <div class="endpoint">
        <h3><span class="method">GET</span> /api/v1/scans</h3>
        <p>获取所有扫描任务的列表。</p>
        
        <h4>响应:</h4>
        <pre><code>[
  {
    "id": "5f3a9c21-7b1e-4c49-8a0d-3e51735ce83a",
    "status": "completed",
    "domain": "example.com",
    "start_time": "2023-08-15T10:30:45Z",
    "end_time": "2023-08-15T10:32:12Z"
  },
  {
    "id": "8e2c7d34-5f1a-42b3-9d6c-1a2b3c4d5e6f",
    "status": "running",
    "domain": "test.com",
    "start_time": "2023-08-15T10:45:12Z",
    "end_time": "0001-01-01T00:00:00Z"
  }
]</code></pre>
    </div>

    <h2>使用示例</h2>

    <h3>使用 curl 发起扫描 (无认证)</h3>
    <pre><code>curl -X POST http://localhost:8080/api/v1/scan \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "bruteforce": false,
    "ports": "80,443",
    "verbose": true
  }'</code></pre>

    <h3>使用 curl 发起扫描 (有认证)</h3>
    <pre><code>curl -X POST http://localhost:8080/api/v1/scan \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "domain": "example.com",
    "bruteforce": false,
    "ports": "80,443",
    "verbose": true
  }'</code></pre>

    <h3>使用 URL 参数进行认证</h3>
    <pre><code>curl -X GET "http://localhost:8080/api/v1/scan/5f3a9c21-7b1e-4c49-8a0d-3e51735ce83a?api_key=your_api_key"</code></pre>

    <h3>获取所有扫描</h3>
    <pre><code>curl -X GET http://localhost:8080/api/v1/scans -H "X-API-Key: your_api_key"</code></pre>

    <h3>使用 Python 进行调用</h3>
    <pre><code>import requests
import json

# API密钥
api_key = "your_api_key"

# 添加认证头
headers = {
    "Content-Type": "application/json",
    "X-API-Key": api_key
}

# 启动扫描
scan_data = {
    "domain": "example.com",
    "bruteforce": True,
    "ports": "80,443",
    "verbose": True
}

response = requests.post("http://localhost:8080/api/v1/scan", 
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
        print(subdomain)</code></pre>
</body>
</html>
`

// GetAPIDocsHTML 返回API文档HTML
func GetAPIDocsHTML() string {
	return APIDocsHTML
}
