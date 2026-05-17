# htop 快速使用指南

## ✅ 安装状态

✓ htop.exe 已编译: `D:\trae\todolist\tools\htop\htop.exe`
✓ PATH 环境变量已配置: `D:\trae\todolist\tools\htop`
✓ PowerShell Profile 已配置

## 🚀 立即使用

**方法 1: 重新打开终端（推荐）**

关闭当前 TRAE IDE 终端，然后重新打开，新终端中就可以直接使用：

```powershell
htop
```

**方法 2: 当前终端快速加载**

在当前终端中运行：

```powershell
$env:Path = $env:Path + ";D:\trae\todolist\tools\htop"
```

然后就可以运行：

```powershell
htop
```

## 📋 快捷键

- `↑` / `↓` - 上下滚动进程列表
- `P` - 按 CPU 使用率排序
- `M` - 按内存使用排序
- `L` - 显示帮助信息
- `Q` - 退出程序

## 🔧 辅助文件

项目中包含以下辅助脚本：

- `install.ps1` - PATH 安装脚本
- `add_to_path.bat` - 批处理快速添加
- `setup_profile.ps1` - PowerShell Profile 配置
- `htop_autoload.ps1` - 自动加载脚本

## 📁 项目结构

```
D:\trae\todolist\tools\htop\
├── htop.exe              ← 主程序（可独立运行）
├── main.go              ← 主程序源码
├── README.md            ← 项目说明
├── INSTALL_GUIDE.md     ← 安装指南
├── QUICKSTART.md        ← 本文档
├── install.ps1          ← PATH 安装脚本
├── add_to_path.bat      ← 快速添加批处理
├── setup_profile.ps1    ← Profile 配置脚本
├── htop_autoload.ps1    ← 自动加载脚本
├── system/              ← 系统信息模块
├── process/             ← 进程管理模块
├── terminal/            ← 终端渲染模块
└── input/              ← 输入处理模块
```

## ⚠️ 常见问题

### Q: 为什么还是不能用 htop 命令？

**A:** 环境变量更改需要重新加载。请：
1. 关闭当前终端
2. 重新打开新终端
3. 在新终端中直接输入 `htop`

### Q: 如何确认配置成功？

**A:** 运行以下命令检查：

```powershell
[Environment]::GetEnvironmentVariable("Path", "User") -split ';' | Where-Object { $_ -like '*htop*' }
```

如果看到 `D:\trae\todolist\tools\htop`，说明配置成功。

### Q: 如何手动运行 htop？

**A:** 直接使用完整路径：

```powershell
D:\trae\todolist\tools\htop\htop.exe
```

## 🎯 下一步

1. **关闭并重新打开 TRAE IDE 终端**
2. **在终端中直接输入 `htop`**
3. **开始监控您的系统进程！**

祝您使用愉快！🎉
