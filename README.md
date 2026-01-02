# Dadoubao (大豆包)

一个语音输入桥接工具，让你可以在手机上用语音输入，文字自动粘贴到电脑光标位置。

## 功能特点

- 手机语音输入，电脑自动接收并粘贴
- 支持自动回车发送
- 扫码即连，无需安装手机 App
- 支持 Windows 和 macOS

## 使用方法

### macOS

1. 下载或编译程序

```bash
# 克隆项目
git clone https://github.com/mxzs123/Dadoubao.git
cd Dadoubao

# 编译
go build -o dadoubao

# 运行
./dadoubao
```

2. 首次运行时，系统会提示授权"辅助功能"权限，请在 系统设置 -> 隐私与安全性 -> 辅助功能 中允许

3. 运行后会在当前目录生成二维码文件 `voicebridge-qr.png`，用手机扫码打开网页

4. 将电脑光标放到需要输入的位置，在手机上使用语音输入，文字会自动粘贴到电脑

### Windows

1. 下载或编译程序

```bash
# 克隆项目
git clone https://github.com/mxzs123/Dadoubao.git
cd Dadoubao

# 编译
go build -o dadoubao.exe

# 运行
dadoubao.exe
```

2. 确保手机和电脑在同一局域网内

3. 运行后会在当前目录生成二维码文件 `voicebridge-qr.png`，用手机扫码打开网页

4. 将电脑光标放到需要输入的位置，在手机上使用语音输入，文字会自动粘贴到电脑

## 编译要求

- Go 1.21 或更高版本

## 注意事项

- 手机和电脑需要在同一局域网内
- macOS 用户需要授权辅助功能权限
- 默认端口为 8765，请确保该端口未被占用

## License

MIT
