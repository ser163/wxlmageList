# 微信素材media_id(媒体ID)获取工具

这是一个用 Go 语言编写的控制台程序，用于通过微信公众平台 API 获取图片素材列表。它支持访问令牌（access_token）的本地缓存管理，并从配置文件中读取必要的凭证。

## 功能

1. 从本地 `config.yaml` 文件读取微信应用的 `appid` 和 `secret`。
2. 检查并管理 `access_token`：
    - 如果本地令牌有效，则直接使用。
    - 如果令牌过期或不存在，自动通过 API 获取新令牌并缓存。
3. 使用 POST 请求从微信 API 获取图片素材（每次最多 20 条）。
4. 遍历并打印素材的 `media_id`（媒体ID）、`name`（名称）和 `url`（链接）。

## 依赖

- Go 语言（建议版本：1.18 或更高）
- 第三方库：
    - `gopkg.in/yaml.v2`（用于解析 YAML 配置文件）

## 项目结构

```
wechat-material-fetcher/
├── main.go           # 主程序文件
├── config.yaml       # 配置文件（需手动创建）
├── access_token.json # 访问令牌缓存文件（运行后自动生成）
└── README.md         # 项目说明文件
```

## 使用方法

### 1. 配置环境

确保您的计算机已安装 Go 环境，可以通过以下命令检查：
```bash
go version
```

### 2. 安装依赖

在项目目录下运行以下命令，安装 YAML 解析库：
```bash
go get gopkg.in/yaml.v2
```

### 3. 创建配置文件

在项目根目录下创建 `config.yaml` 文件，并填入您的微信应用凭证。示例内容如下：
```yaml
appid: 你的appid
secret: 你的secret
```
将 `你的appid` 和 `你的secret` 替换为实际的微信公众平台应用 ID 和密钥。

### 4. 运行程序

#### 方法一：直接运行
在项目目录下执行：
```bash
go run main.go
```

#### 方法二：编译后运行
先编译程序：
```bash
go build -o wechat-material-fetcher
```
然后运行生成的可执行文件：
```bash
./wechat-material-fetcher   # Linux/macOS
wechat-material-fetcher.exe # Windows
```

### 5. 输出结果

程序运行后，会输出类似以下内容：
```
2025/02/20 10:00:00 加载配置文件...
2025/02/20 10:00:00 获取访问令牌...
2025/02/20 10:00:00 使用本地缓存的访问令牌
2025/02/20 10:00:00 获取媒体素材...
2025/02/20 10:00:00 打印素材列表:
媒体ID: xxxx, 名称: image1.jpg, 链接: https://example.com/image1.jpg
媒体ID: yyyy, 名称: image2.png, 链接: https://example.com/image2.png
```

## 编译说明

### 1. 单文件编译
如果您只想生成一个独立的可执行文件，直接在项目目录下运行：
```bash
go build -o wechat-material-fetcher
```
生成的可执行文件可以在支持的平台上运行（需确保 `config.yaml` 在同一目录下）。

### 2. 跨平台编译
如果需要为其他操作系统编译（例如在 Linux 上为 Windows 编译），可以使用以下命令：
- 为 Windows 编译：
  ```bash
  GOOS=windows GOARCH=amd64 go build -o wechat-material-fetcher.exe
  ```
- 为 macOS 编译：
  ```bash
  GOOS=darwin GOARCH=amd64 go build -o wechat-material-fetcher
  ```

完成后，将生成的可执行文件和 `config.yaml` 一起分发到目标机器上。

## 注意事项

1. **网络要求**：程序需要访问微信 API，确保网络连接正常。
2. **凭证有效性**：请确保 `appid` 和 `secret` 正确，否则会导致 API 请求失败。
3. **访问令牌缓存**：令牌保存在 `access_token.json` 中，过期时间为 `expires_in - 5` 秒，程序会自动刷新。
4. **日志时间**：日志中的时间戳基于当前系统时间，示例中为 2025 年 2 月 20 日（根据需求日期）。

## 错误排查

- 如果出现 `读取配置文件失败`，检查 `config.yaml` 是否存在且格式正确。
- 如果出现 `获取令牌失败`，检查网络连接或确认 `appid` 和 `secret` 是否有效。
- 如果素材列表为空，可能是账户下无图片素材或 API 返回异常。

## 贡献

欢迎提交问题或改进建议！请通过 GitHub Issues 或 Pull Requests 联系我。
