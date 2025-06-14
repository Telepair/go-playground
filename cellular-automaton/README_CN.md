# 元胞自动机

_[English Version / 英文版本](README.md)_

一个高度可定制的一维元胞自动机终端用户界面(TUI)实现。

## 功能特性

- **规则驱动生成**: 支持所有 256 种基本元胞自动机规则 (0-255)
- **双运行模式**:
  - 有限模式: 运行指定步数
  - 无限模式: 持续生成并实时可视化
- **自动窗口检测**: 自动检测终端大小或允许手动指定
- **高度可定制渲染**:
  - 可配置元胞大小 (1-3 倍)
  - 自定义颜色 (十六进制格式)
  - 自定义渲染字符
- **动态刷新**: 可配置刷新频率 (默认 1 秒，最小 1 毫秒)
- **灵活窗口大小**: 自动检测或使用 宽度 x 高度 格式手动指定 (如: 100x80)
- **双语支持**: 中英文界面

## 安装

```bash
# 克隆仓库
git clone <repository-url>
cd cellular-automaton

# 构建应用程序
go build -o cellular-automaton
```

## 使用方法

### 基本命令

```bash
# 使用默认设置运行 (规则 30，自动窗口大小)
./cellular-automaton

# 指定规则和自定义参数
./cellular-automaton -rule 110 -steps 100 -size 120x60

# 自动窗口检测的无限模式
./cellular-automaton -rule 30 -steps 0 -size auto

# 自定义渲染样式
./cellular-automaton -rule 90 -cellsize 3 -alive-char "●" -dead-char "○"
```

### 命令行选项

- `-rule <数字>`: 元胞自动机规则 (0-255，默认: 30)
- `-steps <数字>`: 步数 (0 或负数为无限模式，默认: 50)
- `-size <大小>`: 窗口大小 (格式: 宽度 x 高度，如: 100x80，或 'auto' 自动检测，默认: auto)
- `-cellsize <大小>`: 元胞渲染大小 (1-3，默认: 2)
- `-alive-color <颜色>`: 活跃元胞颜色，十六进制格式 (默认: #FFFFFF)
- `-dead-color <颜色>`: 死亡元胞颜色，十六进制格式 (默认: #000000)
- `-alive-char <字符>`: 活跃元胞字符 (默认: █)
- `-dead-char <字符>`: 死亡元胞字符 (默认: 空格)
- `-refresh <秒数>`: 刷新频率(秒)，最小值 0.001 (默认: 0.1)
- `-boundary <类型>`: 边界条件类型 (periodic/fixed/reflect，默认: periodic)
- `-lang <en/cn>`: 界面语言 (默认: en)

### Makefile

```bash
# 显示帮助
make help

# 构建应用程序
make build

# 运行演示
make cellular-automaton-basic
make cellular-automaton-sierpinski
make cellular-automaton-turing
make cellular-automaton-traffic
make cellular-automaton-infinite
make cellular-automaton-colorful
make cellular-automaton-fixed
make cellular-automaton-periodic
make cellular-automaton-reflect
```

### 示例命令

```bash
# 自动检测窗口大小的规则 30
./cellular-automaton -rule 30 -size auto

# 快速刷新的无限模式
./cellular-automaton -rule 30 -steps 0 -refresh 0.1

# ASCII 艺术风格自定义字符
./cellular-automaton -rule 184 -alive-char "■" -dead-char "□" -cellsize 1

# 详细图案的大窗口尺寸
./cellular-automaton -rule 110 -size 200x100 -steps 150

# 固定边界条件 (不循环)
./cellular-automaton -rule 30 -boundary fixed

# 反射边界条件
./cellular-automaton -rule 110 -boundary reflect

# 中文界面
./cellular-automaton -rule 30 -lang cn
```

## 控制按键

### 所有模式

- **q** 或 **Ctrl+C**: 退出应用程序
- **l**: 切换语言 (英文/中文)

### 无限模式专用

- **空格键** 或 **回车键**: 暂停/继续模拟

### 高级控制 (无限模式)

- **r**: 重置模拟到初始状态
- **+** 或 **=**: 提高刷新频率 (加快模拟速度)
- **-** 或 **\_**: 降低刷新频率 (减慢模拟速度)
- **1**, **2**, **3**: 改变元胞渲染大小 (1 倍, 2 倍, 3 倍)

## 有趣的规则推荐

- **规则 30**: 混沌、伪随机图案
- **规则 90**: 谢尔宾斯基三角形图案
- **规则 110**: 图灵完备，复杂的涌现行为
- **规则 184**: 交通流模拟
- **规则 150**: XOR 图案，创建分形结构

## 技术细节

### 边界条件

元胞自动机支持三种边界条件类型:

- **周期性 (默认)**: 最左边元胞的左邻居是最右边的元胞，最右边元胞的右邻居是最左边的元胞 (循环行为)
- **固定**: 最左边元胞的左邻居始终为 0 (死亡)，最右边元胞的右邻居始终为 0 (死亡)
- **反射**: 最左边元胞的左邻居是它自己，最右边元胞的右邻居是它自己 (镜像行为)

### 窗口大小检测

应用程序可以自动检测终端窗口大小:

- 使用 `-size auto` (默认) 进行自动检测
- 应用程序为 UI 元素预留空间 (标题、控制)
- 强制执行最小高度以确保正确显示
- 手动大小指定会覆盖自动检测

### 性能

- **刷新频率**: 支持高频更新 (最小 1 毫秒)
- **内存高效**: 仅存储必要的网格数据，元胞计算零内存分配
- **优化渲染**: 使用预计算样式字符串和高效字符串构建
- **缓冲区复用**: 通过缓冲池减少垃圾回收
- **实时控制**: 无需重启即可动态调整刷新频率和元胞大小

### 图案分析

不同规则产生不同的图案类型:

- **第 1 类**: 演化为均匀状态 (如: 规则 0, 32)
- **第 2 类**: 演化为简单周期结构 (如: 规则 4, 108)
- **第 3 类**: 混沌、非周期行为 (如: 规则 30, 45)
- **第 4 类**: 复杂、局部化结构 (如: 规则 110, 124)

## 贡献

1. Fork 仓库
2. 创建功能分支
3. 进行更改
4. 如适用，添加测试
5. 提交 pull request

## 许可证

此项目是开源的。请查看许可证文件了解详情。
