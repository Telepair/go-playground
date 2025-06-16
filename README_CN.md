# Go Playground

_[English Version / 英文版本](README.md)_

一个用 Go 语言实现各种有趣小程序的项目集合。每个子项目都是一个独立的程序，展示不同的编程概念、算法或有趣的实现。

## 项目列表

### 🧬 [元胞自动机 (Cellular Automaton)](./cellular-automaton/)

一个交互式的一维元胞自动机程序，使用 Bubble Tea 实现美观的终端用户界面。

[Wikipedia - Cellular Automaton](https://en.wikipedia.org/wiki/Cellular_automaton)

**演示:**

[![asciicast](https://asciinema.org/a/723358.svg)](https://asciinema.org/a/723358)

### 🎮 [康威生命游戏 (Conway's Game of Life)](./conway-game-of-life/)

一个具有多种预定义模式和高度可定制渲染选项的康威生命游戏终端用户界面(TUI)实现。

[Wikipedia - Conway's Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life)

**演示:**

[![asciicast](https://asciinema.org/a/723376.svg)](https://asciinema.org/a/723376)

## 项目结构

```
go-playground/
├── README.md                    # 主项目说明
├── cellular-automaton/          # 元胞自动机
├── conway-game-of-life/         # 康威生命游戏
├── LICENSE                     # 项目许可证
└── .gitignore                 # Git 忽略文件
```

## 使用 Asciinema 录制演示

1. 安装 asciinema:

   ```bash
   # macOS
   brew install asciinema

   # Linux
   pip install asciinema
   ```

2. 录制演示:

   ```bash
   # Start recording
   # Note: After the program finishes running, press 'Q' to quit the program and complete the recording

   # Cellular Automaton
   asciinema rec ./cellular-automaton.cast --title "Cellular Automaton" --command "./bin/cellular-automaton"

   # Conway Game of Life
   asciinema rec ./conway-game-of-life.cast --title "Conway Game of Life" --command "./bin/conway-game-of-life"
   ```

3. 播放演示:

   ```bash
   # Cellular Automaton
   asciinema play ./cellular-automaton.cast

   # Conway Game of Life
   asciinema play ./conway-game-of-life.cast
   ```

4. 上传到 asciinema.org (可选):

   ```bash
   # Cellular Automaton
   asciinema upload ./cellular-automaton.cast

   # Conway Game of Life
   asciinema upload ./conway-game-of-life.cast
   ```

## 技术特点

- **现代 Go 开发**：使用 Go 1.24+ 的最新特性
- **优雅的用户界面**：使用 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 和 [Lipgloss](https://github.com/charmbracelet/lipgloss) 构建美观的终端界面
- **独立模块设计**：每个子项目都有独立的 `go.mod`，方便管理和使用
- **清晰的代码结构**：注重代码可读性和可维护性
- **详细的文档**：每个项目都有完整的使用说明和示例

## 计划中的项目

未来可能会添加的有趣项目：

- 🎮 **生命游戏 (Game of Life)** - 康威生命游戏实现
- 🧮 **曼德博集合 (Mandelbrot Set)** - 曼德博集合可视化
- 🎵 **音频可视化器 (Audio Visualizer)** - 音频频谱可视化
- 🌊 **波函数坍缩 (Wave Function Collapse)** - 波函数坍缩算法
- 🎲 **随机游走 (Random Walk)** - 随机游走可视化
- 📊 **数据结构可视化 (Data Structures Visualization)** - 数据结构可视化
- 🔍 **算法可视化 (Algorithm Visualization)** - 排序和搜索算法可视化

## 贡献

欢迎提交 Issue 和 Pull Request！如果你有有趣的想法或发现了 bug，请随时联系。

### 贡献指南

1. Fork 这个项目
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 环境要求

- Go 1.24.4 或更高版本
- 支持 Unicode 的终端（推荐使用现代终端如 iTerm2、Windows Terminal 等）

## 许可证

本项目采用 [许可证名称] 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

---

**享受用 Go 语言探索各种有趣概念的乐趣！** 🚀
