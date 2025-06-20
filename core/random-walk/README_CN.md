# 随机游走可视化

[English Version / 英文版本](README.md)

[Wikipedia - Random Walk](https://en.wikipedia.org/wiki/Random_walk)

基于终端的随机游走算法可视化工具，使用 Go 语言和 Bubble Tea 框架实现。

[![asciicast](https://asciinema.org/a/723662.svg)](https://asciinema.org/a/723662)

## 功能特性

- **多种游走模式**：

  - **单粒子模式**：经典的单粒子随机游走
  - **多粒子模式**：多个粒子同时进行随机游走
  - **轨迹模式**：显示单个粒子的运动轨迹
  - **布朗运动**：模拟连续运动的布朗运动
  - **自避行走**：粒子不能重复访问已经走过的位置
  - **莱维飞行**：偶尔进行长距离跳跃的随机游走

- **交互式控制**：
  - 实时可视化，速度可调
  - 暂停/恢复功能
  - 动态调整粒子数量（多粒子模式）
  - 可配置的轨迹长度
  - 双语支持（中文/英文）

## 安装

### 前置要求

- Go 1.22 或更高版本
- 支持 Unicode 的终端

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/telepair/go-playground.git
cd go-playground

# 构建随机游走可视化程序
make build-random-walk

# 或直接构建
go build -o ./bin/random-walk ./random-walk
```

## 使用方法

### 基本使用

```bash
# 使用默认设置运行
./bin/random-walk

# 或使用 make
make random-walk
```

### 命令行选项

```bash
./bin/random-walk [选项]

选项：
  -walker-color string    粒子颜色（十六进制格式）（默认 "#FF00FF"）
  -trail-color string     轨迹颜色（十六进制格式）（默认 "#0088FF"）
  -empty-color string     空白单元格颜色（十六进制格式）（默认 "#000000"）
  -walker-char string     粒子字符（默认 "●"）
  -trail-char string      轨迹字符（默认 "·"）
  -empty-char string      空白单元格字符（默认 " "）
  -lang string           语言：en 或 cn（默认 "en"）
  -profile               启用性能分析和监控
  -profile-port int      性能分析服务器端口（默认 6060）
  -log-file string       调试日志文件路径
```

### 使用示例

```bash
# 使用自定义粒子和轨迹字符
./bin/random-walk -walker-char '🐾' -trail-char '·'

# 自定义颜色
./bin/random-walk -walker-color '#FF00FF' -trail-color '#00FFFF'

# 中文界面运行
./bin/random-walk -lang cn

# 启用性能分析
./bin/random-walk -profile -log-file debug.log
```

## 控制键

| 按键             | 功能                            |
| ---------------- | ------------------------------- |
| `M`              | 切换游走模式                    |
| `W/w`            | 增加/减少粒子数量（多粒子模式） |
| `T/t`            | 增加/减少轨迹长度（轨迹模式）   |
| `+/-` 或 `↑/↓`   | 加速/减速                       |
| `空格` 或 `回车` | 暂停/恢复                       |
| `L`              | 切换语言（中文/英文）           |
| `R`              | 重置模拟                        |
| `Q` 或 `Esc`     | 退出                            |

## 游走模式说明

### 单粒子模式

经典的随机游走，单个粒子在 8 个方向（包括对角线）随机移动。

### 多粒子模式

多个粒子同时进行随机游走，每个粒子有独特的颜色。适合研究碰撞和覆盖模式。

### 轨迹模式

显示单个粒子的运动路径，轨迹会随时间逐渐淡化。

### 布朗运动

模拟布朗运动，粒子以随机方向和距离进行连续运动。

### 自避行走

粒子不能重新访问已经走过的位置。如果没有可用的移动选项，粒子可能会被困住。

### 莱维飞行

粒子偶尔会进行长距离跳跃的随机游走，模拟自然界中发现的莱维飞行模式。

## 技术细节

### 实现

- 使用 Go 语言和 [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI 框架编写
- 使用 [Lip Gloss](https://github.com/charmbracelet/lipgloss) 进行样式设计
- 实现了高效的网格渲染，最小化内存分配
- 支持终端大小调整并保持纵横比

### 性能

- 优化为 50ms 刷新率的流畅动画
- 使用强度衰减的高效轨迹渲染
- 渲染过程中最小化内存分配

## 开发

### 运行测试

```bash
# 运行所有测试
make test

# 运行基准测试
make bench
```

### 项目结构

```
random-walk/
├── main.go          # 程序入口
├── config.go        # 配置和常量
├── walk.go          # 核心随机游走逻辑
├── ui.go            # UI 和交互逻辑
├── styles.go        # 视觉样式和渲染
├── walk_test.go     # 单元测试
├── README.md        # 英文文档
└── README_CN.md     # 中文文档
```

## 许可证

本项目是 go-playground 集合的一部分，遵循相同的许可条款。
