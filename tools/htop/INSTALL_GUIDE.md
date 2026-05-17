# htop 安装和使用说明

## 自动安装（推荐）

我已经为您创建了自动安装脚本。运行以下命令即可：

```powershell
.\add_to_path.bat
```

或者复制以下命令直接粘贴到终端：

```powershell
$env:Path = $env:Path + ";D:\trae\todolist\tools\htop"
```

## 永久添加到 PATH

如果您希望在新打开的终端中自动使用 htop 命令，请运行：

```powershell
[Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";D:\trae\todolist\tools\htop", "User")
```

然后**关闭当前终端并重新打开**，新终端中就可以直接使用 `htop` 命令了。

## 使用方法

安装完成后，您可以：

1. **在任意位置直接运行**：
   ```powershell
   htop
   ```

2. **使用完整路径运行**：
   ```powershell
   D:\trae\todolist\tools\htop\htop.exe
   ```

3. **相对路径运行**（在 htop 目录中）：
   ```powershell
   .\htop.exe
   ```

## 快捷键说明

- `↑` / `↓` - 上下滚动进程列表
- `P` - 按 CPU 使用率排序
- `M` - 按内存使用排序
- `L` - 显示帮助信息
- `Q` - 退出程序

## 故障排除

### Q: 为什么在新终端中还是不能用 htop 命令？

A: 请确保：
1. 运行了安装脚本
2. **关闭并重新打开终端**（重要！）

### Q: 如何确认 htop 已添加到 PATH？

A: 运行以下命令检查：
```powershell
[Environment]::GetEnvironmentVariable("Path", "User") -split ';' | Where-Object { $_ -like '*htop*' }
```

### Q: 如何手动添加 htop 到当前会话？

A: 运行：
```powershell
$env:Path = $env:Path + ";D:\trae\todolist\tools\htop"
```

## 文件位置

- htop.exe: `D:\trae\todolist\tools\htop\htop.exe`
- 安装脚本: `D:\trae\todolist\tools\htop\add_to_path.bat`
- 完整项目: `D:\trae\todolist\tools\htop\`
