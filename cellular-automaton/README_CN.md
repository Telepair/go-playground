# 元胞自动机

_[English Version / 英文版本](README.md)_

[Wikipedia - Cellular Automaton](https://en.wikipedia.org/wiki/Cellular_automaton)

一个高度可定制的一维元胞自动机终端用户界面(TUI)实现。

[![asciicast](https://asciinema.org/a/723358.svg)](https://asciinema.org/a/723358)

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
./cellular-automaton
./cellular-automaton -rule 30
./cellular-automaton -rule 90 -alive-color "#00FF00" -dead-color "#FF0000"
./cellular-automaton -rule 110 -alive-char "●" -dead-char "○"
./cellular-automaton -rule 150
./cellular-automaton -rule 184 -alive-char '🚗' -rows 30 -cols 80
```

### 命令行选项

- `-rule <数字>`: 元胞自动机规则 (0-255，默认: 30)
- `-rows <数字>`: 行数 (默认: 30)
- `-cols <数字>`: 列数 (默认: 80)
- `-alive-color <颜色>`: 活跃元胞颜色，十六进制格式 (默认: #FFFFFF)
- `-dead-color <颜色>`: 死亡元胞颜色，十六进制格式 (默认: #000000)
- `-alive-char <字符>`: 活跃元胞字符 (默认: █)
- `-dead-char <字符>`: 死亡元胞字符 (默认: 空格)
- `-lang <en/cn>`: 界面语言 (默认: en)

## 控制按键

### 所有模式

- **t**: 切换规则 (从常用规则中选择或输入自定义规则 0-255)
- **b**: 切换边界类型 (周期性/固定/反射)
- **r**: 重置模拟到初始状态
- **l**: 切换语言 (英文/中文)
- **+** 或 **=**: 提高刷新频率 (加快模拟速度)
- **-** 或 **\_**: 降低刷新频率 (减慢模拟速度)
- **空格键** 或 **回车键**: 暂停/继续模拟
- **q** 或 **Ctrl+C**: 退出应用程序

## 有趣的规则推荐

- **规则 30**: 混沌、伪随机图案
- **规则 90**: 谢尔宾斯基三角形图案
- **规则 110**: 图灵完备，复杂的涌现行为
- **规则 150**: XOR 图案，创建分形结构
- **规则 184**: 交通流模拟

## 技术细节

### 边界条件

元胞自动机支持三种边界条件类型:

- **周期性 (默认)**: 最左边元胞的左邻居是最右边的元胞，最右边元胞的右邻居是最左边的元胞 (循环行为)
- **固定**: 最左边元胞的左邻居始终为 0 (死亡)，最右边元胞的右邻居始终为 0 (死亡)
- **反射**: 最左边元胞的左邻居是它自己，最右边元胞的右邻居是它自己 (镜像行为)
