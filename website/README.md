# go-Sublist3r 官方网站

这是 go-Sublist3r 项目的官方网站源代码，设计用于通过 GitHub Pages 进行部署。

## 网站结构

```
website/
├── assets/
│   ├── css/         # 样式文件
│   ├── js/          # JavaScript文件
│   └── img/         # 图像文件
├── docs/            # 文档页面
├── .nojekyll        # 禁用GitHub Pages的Jekyll处理
├── favicon.ico      # 网站图标
└── index.html       # 主页
```

## 本地预览

要在本地预览网站，您可以使用任何简单的HTTP服务器。例如，使用Python：

```bash
# 如果您使用Python 3
cd website
python -m http.server 8000

# 或者使用Python 2
cd website
python -m SimpleHTTPServer 8000
```

然后在浏览器中访问 `http://localhost:8000`。

## 部署到GitHub Pages

### 方法1：使用主分支的/docs目录

1. 将`website`目录重命名为`docs`
2. 将目录推送到GitHub仓库的主分支
3. 在GitHub仓库设置中，启用GitHub Pages，并选择主分支下的`/docs`目录作为源

### 方法2：使用gh-pages分支

1. 创建一个`gh-pages`分支
2. 将`website`目录中的内容复制到`gh-pages`分支的根目录
3. 推送`gh-pages`分支到GitHub
4. 在GitHub仓库设置中，启用GitHub Pages，并选择`gh-pages`分支作为源

### 使用GitHub Actions自动部署

您也可以设置GitHub Actions工作流来自动部署网站。下面是一个示例工作流配置：

创建文件`.github/workflows/deploy-website.yml`：

```yaml
name: Deploy Website to GitHub Pages

on:
  push:
    branches:
      - main
    paths:
      - 'website/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '16'
          
      - name: Deploy to GitHub Pages
        uses: JamesIves/github-pages-deploy-action@v4
        with:
          branch: gh-pages  # 部署到的分支
          folder: website   # 源文件夹
          clean: true       # 自动清理旧文件
```

## 自定义域名

如果您想使用自定义域名，可以创建一个`CNAME`文件：

```bash
echo "your-domain.com" > website/CNAME
```

然后在您的域名提供商处添加适当的DNS记录。

## 贡献

如果您想改进网站，请提交Pull Request。确保您的更改在本地测试通过。 